package manager

import (
	"os"
	"syscall"
)

// common variables (chans and vars)
var (
	// Channels
	chOsSignal chan os.Signal
	SIG_STOP   = syscall.SIGQUIT
)

// Initialize the os signal channel
func init() {
	chOsSignal = make(chan os.Signal, 1)
}

// Get the os signal channel
func getchOsSignal() chan os.Signal {
	return chOsSignal
}