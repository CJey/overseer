// +build windows

package overseer

import (
	"syscall"
)

var (
	supported = true
	uid       = syscall.Getuid()
	gid       = syscall.Getgid()
	SIGUSR1   = syscall.SIGTERM
	SIGUSR2   = syscall.SIGTERM
	SIGTERM   = syscall.SIGTERM
)
