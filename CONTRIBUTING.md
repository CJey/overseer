
### Issues

If you've found a bug, please create an [issue](https://github.com/jpillora/overseer/issues) or if possible, create a fix and send in a pull request.

## Contributing

If you'd like to contribute, please see the notes below and create an issue mentioning want you want to work on and if you're creating an addition to the core `overseer` repo, then also include the proposed API.

## Issues and bug fixes

### Tests

`overseer` needs a test which suite should drive an:

* HTTP client for verifying application version
* HTTP server for providing application upgrades
* an `overseer` process via `exec.Cmd`

And as it operates, confirm each phase.

### Updatable config

Child process should pass new config back to the main process and:
* Update logging settings
* Update socket bindings
