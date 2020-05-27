// +build !linux,!darwin,!windows,!freebsd

package overseer

import (
	"os"
)

var (
	supported = false
	uid       = 0
	gid       = 0
	SIGUSR1   = os.Interrupt
	SIGUSR2   = os.Interrupt
	SIGTERM   = os.Kill
)
