// Package overseer implements daemonizable
// self-upgrading binaries in Go (golang).
package overseer

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

const (
	envSlaveID = "OVERSEER_SLAVE_ID"
	envIsSlave = "OVERSEER_IS_SLAVE"
	envNumFDs  = "OVERSEER_NUM_FDS"
	envBinPath = "OVERSEER_BIN_PATH"
)

var (
	isSlave = os.Getenv(envIsSlave) == "1"
	slaveID = os.Getenv(envSlaveID)
)

// Config defines overseer's run-time configuration
type Config struct {
	//Required will prevent overseer from fallback to running
	//running the program in the main process on failure.
	Required bool
	//Program's main function
	Program func(state State)
	//Program's zero-downtime socket listening address (set this or Addresses)
	Address string
	//Program's zero-downtime socket listening addresses (set this or Address)
	Addresses []string
	//RestartSignal will manually trigger a graceful restart. Defaults to SIGUSR2.
	RestartSignal os.Signal
	//TerminateTimeout controls how long overseer should
	//wait for the program to terminate itself. After this
	//timeout, overseer will issue a SIGKILL.
	TerminateTimeout time.Duration
	//Debug enables all [overseer] logs.
	Debug bool
	//NoWarn disables warning [overseer] logs.
	NoWarn bool
	//NoRestart disables all restarts, this option essentially converts
	//the RestartSignal into a "ShutdownSignal".
	NoRestart bool
	//SingleAccept allows accept same listener by only one slave.
	//By default, when graceful stop signal received by slave,
	//slave do not close listeners, you should close all listeners by youself,
	//that means old slave and new slave will accept the same listeners simultaneously.
	//Otherwise, slave close all listeners immediately.
	SingleAccept bool
}

func validate(c *Config) error {
	//validate
	if c.Program == nil {
		return errors.New("overseer.Config.Program required")
	}
	if c.Address != "" {
		if len(c.Addresses) > 0 {
			return errors.New("overseer.Config.Address and Addresses cant both be set")
		}
		c.Addresses = []string{c.Address}
	} else if len(c.Addresses) > 0 {
		c.Address = c.Addresses[0]
	}
	if c.RestartSignal == nil {
		c.RestartSignal = SIGUSR2
	}
	if c.TerminateTimeout <= 0 {
		c.TerminateTimeout = 30 * time.Second
	}
	return nil
}

//RunErr allows manual handling of any
//overseer errors.
func RunErr(c Config) error {
	return runErr(&c)
}

//Run executes overseer, if an error is
//encountered, overseer fallsback to running
//the program directly (unless Required is set).
func Run(c Config) {
	err := runErr(&c)
	if err != nil {
		if c.Required {
			log.Fatalf("[overseer] %s", err)
		} else if c.Debug || !c.NoWarn {
			log.Printf("[overseer] disabled. run failed: %s", err)
		}
		c.Program(DisabledState)
		return
	}
	os.Exit(0)
}

//abstraction over master/slave
var currentProcess interface {
	triggerRestart()
	run() error
}

func runErr(c *Config) error {
	//os not supported
	if !supported {
		return fmt.Errorf("os (%s) not supported", runtime.GOOS)
	}
	if err := validate(c); err != nil {
		return err
	}
	//run either in master or slave mode
	if IsSlave() {
		currentProcess = &slave{Config: c}
	} else {
		currentProcess = &master{Config: c}
	}
	return currentProcess.run()
}

//IsMaster returns current running mode is master or not
func IsMaster() bool {
	return !isSlave
}

//IsSlave returns current running mode is slave or not
func IsSlave() bool {
	return isSlave
}

//SlaveID returns current slave's number id
func SlaveID() string {
	return slaveID
}

//Restart programmatically triggers a graceful restart. If NoRestart
//is enabled, then this will essentially be a graceful shutdown.
func Restart() {
	if currentProcess != nil {
		currentProcess.triggerRestart()
	}
}

//IsSupported returns whether overseer is supported on the current OS.
func IsSupported() bool {
	return supported
}
