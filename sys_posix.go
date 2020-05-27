// +build linux darwin freebsd

package overseer

//this file attempts to contain all posix
//specific stuff, that needs to be implemented
//in some other way on other OSs... TODO!

import (
	"syscall"
)

var (
	supported = true
	uid       = syscall.Getuid()
	gid       = syscall.Getgid()
	SIGUSR1   = syscall.SIGUSR1
	SIGUSR2   = syscall.SIGUSR2
	SIGTERM   = syscall.SIGTERM
)
