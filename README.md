# overseer

[![GoDoc](https://godoc.org/github.com/cjey/overseer?status.svg)](https://godoc.org/github.com/cjey/overseer)

`overseer` is a package for creating monitorable, gracefully restarting, self-upgrading binaries in Go (golang). The main goal of this project is to facilitate the creation of self-upgrading binaries which play nice with standard process managers, secondly it should expose a small and simple API with reasonable defaults.

![overseer diagram](https://docs.google.com/drawings/d/1o12njYyRILy3UDs2E6JzyJEl0psU4ePYiMQ20jiuVOY/pub?w=566&h=284)

Commonly, graceful restarts are performed by the active process (*dark blue*) closing its listeners and passing these matching listening socket files (*green*) over to a newly started process. This restart causes any **foreground** process monitoring to incorrectly detect a program crash. `overseer` attempts to solve this by using a small process to perform this socket file exchange and proxying signals and exit code from the active process.

### Features

* Simple
* Works with process managers (systemd, upstart, supervisor, etc)
* Graceful, zero-down time restarts
* Easy self-upgrading binaries

### Install

```sh
go get github.com/cjey/overseer
```

### Quick example

This program works with process managers, supports graceful, zero-down time restarts and self-upgrades its own binary.

``` go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cjey/overseer"
)

//create another main() to run the overseer process
//and then convert your old main() into a 'prog(state)'
func main() {
	overseer.Run(overseer.Config{
		Program: prog,
		Address: ":3000",
	})
}

//prog(state) runs in a child process
func prog(state overseer.State) {
	log.Printf("app (%s) listening...", state.ID)
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "app (%s) says hello\n", state.ID)
	}))
	http.Serve(state.Listener, nil)
}
```

**How it works:**

* `overseer` uses the main process to check for and install upgrades and a child process to run `Program`.
* The main process retrieves the files of the listeners described by `Address/es`.
* The child process is provided with these files which is converted into a `Listener/s` for the `Program` to consume.
* All child process pipes are connected back to the main process.
* All signals received on the main process are forwarded through to the child process.
* Except for scheduled restarts, the active child process exiting will cause the main process to exit with the same code. So, **`overseer` is not a process manager**.

See [Config](https://godoc.org/github.com/cjey/overseer#Config)uration options [here](https://godoc.org/github.com/cjey/overseer#Config) and the runtime [State](https://godoc.org/github.com/cjey/overseer#State) available to your program [here](https://godoc.org/github.com/cjey/overseer#State).

### More examples

See the [example/](example/) directory and run `example.sh`, you should see the following output:

```sh
$ cd example/
$ sh example.sh
BUILT APP (1)
RUNNING APP
app#1 listening...
app#1 says hello
app#1 says hello
BUILT APP (2)
app#2 listening...
app#2 says hello
app#2 says hello
app#1 says hello
app#1 exiting...
BUILT APP (3)
app#3 listening...
app#3 says hello
app#3 says hello
app#2 says hello
app#2 exiting...
app#3 says hello
```

**Note:** `app#1` stays running until the last request is closed.

#### Only use graceful restarts

```go
func main() {
	overseer.Run(overseer.Config{
		Program: prog,
		Address: ":3000",
	})
}
```

Send `main` a `SIGUSR2` (`Config.RestartSignal`) to manually trigger a restart

### Known issues

* The master process's `overseer.Config` cannot be changed via an upgrade, the master process must be restarted.
	* Therefore, `Addresses` can only be changed by restarting the main process.
* Currently shells out to `mv` for moving files because `mv` handles cross-partition moves unlike `os.Rename`.
* Only supported on darwin and linux.
* Package `init()` functions will run twice on start, once in the main process and once in the child process.

### More documentation

* [Core `overseer` package](https://godoc.org/github.com/cjey/overseer)

### Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)
