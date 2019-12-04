package manager

import (
	"syscall"
)

var (
	SigUserStop = syscall.Signal(0x14)
)
