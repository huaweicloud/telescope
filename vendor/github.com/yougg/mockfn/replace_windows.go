package mockfn

import (
	"syscall"
	"unsafe"
)

const PageExecuteReadwrite = 0x40

var subFeatures = []int{22, 26, 56}
var procVirtualProtect = syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualProtect")

func virtualProtect(lpAddress uintptr, dwSize int, flNewProtect uint32, lpflOldProtect unsafe.Pointer) error {
	ret, _, _ := procVirtualProtect.Call(lpAddress, uintptr(dwSize), uintptr(flNewProtect), uintptr(lpflOldProtect))
	if ret == 0 {
		return syscall.GetLastError()
	}
	return nil
}

// this function is super unsafe
// aww yeah
// It copies a slice to a raw memory location, disabling all memory protection before doing so.
func copyToLocation(location uintptr, data []byte) {
	f := rawMemoryAccess(location, len(data))

	var oldPerms uint32
	err := virtualProtect(location, len(data), PageExecuteReadwrite, unsafe.Pointer(&oldPerms))
	if err != nil {
		panic(err)
	}
	copy(f, data[:])

	// VirtualProtect requires you to pass in a pointer which it can write the
	// current memory protection permissions to, even if you don't want them.
	var tmp uint32
	err = virtualProtect(location, len(data), oldPerms, unsafe.Pointer(&tmp))
	if err != nil {
		panic(err)
	}
}
