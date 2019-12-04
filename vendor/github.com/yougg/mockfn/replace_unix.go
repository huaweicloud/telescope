// +build !windows,!android,!plan9,!nacl

package mockfn

import (
	"golang.org/x/sys/unix"
)

var subFeatures = []int{15, 19, 24, 27, 31, 49}

// this function is super unsafe
// aww yeah
// It copies a slice to a raw memory location, disabling all memory protection before doing so.
func copyToLocation(location uintptr, data []byte) {
	f := rawMemoryAccess(location, len(data))
	mprotectCrossPage(location, len(data), unix.PROT_READ|unix.PROT_WRITE|unix.PROT_EXEC)
	copy(f, data[:])
	mprotectCrossPage(location, len(data), unix.PROT_READ|unix.PROT_EXEC)
}

func mprotectCrossPage(addr uintptr, len int, prot int) {
	pageSize := unix.Getpagesize()
	for p := pageStart(addr); p <= addr+uintptr(len); p += uintptr(pageSize) {
		page := rawMemoryAccess(p, pageSize)
		err := unix.Mprotect(page, prot)
		if err != nil {
			panic(err)
		}
	}
}
