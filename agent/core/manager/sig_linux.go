package manager

import (
	"syscall"
)

var (
	SigUserStop = syscall.SIGUSR2
)
