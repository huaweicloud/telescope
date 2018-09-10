package mockfn

import (
	"fmt"
	"reflect"
	"syscall"
	"unsafe"
)

func rawMemoryAccess(p uintptr, length int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p,
		Len:  length,
		Cap:  length,
	}))
}

func pageStart(ptr uintptr) uintptr {
	return ptr & ^(uintptr(syscall.Getpagesize() - 1))
}

// from is a pointer to the actual function
// to is a pointer to a go func value
func replaceFunction(from, to uintptr) (original []byte) {
	jumpData := jmpToFunctionValue(to)
	f := rawMemoryAccess(from, len(jumpData))
	original = make([]byte, len(f))
	copy(original, f)
	if isAlreadyReplaced(original) {
		panic("func already mocked")
	}

	copyToLocation(from, jumpData)
	return
}

func isAlreadyReplaced(bytes []byte) bool {
	return bytes[0] == 0x48 && bytes[1] == 0xBA && bytes[10] == 0xFF && bytes[11] == 0x22
}

func printRawData(bytes []byte) {
	fmt.Println("\nraw data: ")
	var i = 0
	for _, v := range bytes {
		i++
		fmt.Printf("0x%02X ", v)
		if i > 0 && (i%8) == 0 {
			fmt.Printf("\n")
		}
	}

	fmt.Println()
}

func memcopy(from uintptr, len int) []byte {
	f := rawMemoryAccess(from, len)
	original := make([]byte, len)
	copy(original, f)

	return original
}

func replaceJBE(target, alias uintptr) (targetOffset uintptr, aliasOffset uintptr, aliasOriginal []byte) {
	aHead := copyMoreStack(memcopy(alias, 60))
	tHead := copyMoreStack(memcopy(target, 60))

	_, aAddrLen, moreStackOffset, ok := findJBEorJE(alias, aHead, 0)
	if !ok {
		printRawData(aHead)
		panic("jbe not found at alias head\n")
	}

	tPos, tAddrLen, _, ok := findJBEorJE(target, tHead, 0)
	if !ok {
		printRawData(tHead)
		panic("jbe not found at target head\n")
	}

	for ok {
		if tAddrLen < aAddrLen {
			panic(fmt.Sprintf("tAddrLen(%d) < aAddrlen(%d)\n", tAddrLen, aAddrLen))
		}

		addr := make([]byte, tAddrLen)
		copy(addr, int2Bytes(int32(moreStackOffset-tPos-tAddrLen)))
		for i := 0; i < tAddrLen; i++ {
			tHead[tPos+i] = addr[i]
		}

		tPos, tAddrLen, _, ok = findJBEorJE(target, tHead, tPos+tAddrLen)
	}

	original := memcopy(alias, len(tHead))

	copyToLocation(alias, tHead)

	return uintptr(len(tHead)), uintptr(len(tHead)), original
}

func findJBEorJE(ptr uintptr, buffer []byte, offset int) (pos int, addrLen int, moreStackOffset int, ok bool) {
	m := []map[string][]byte{ //                             instruction | offset
		{"code": {0x76}, "insLen": {2}},       // jbe addr   0x76          0x61
		{"code": {0x0f, 0x86}, "insLen": {6}}, // jbe addr   0x0f 0x86     0xd1 0x00 0x00 0x00
		{"code": {0x0f, 0x84}, "insLen": {6}}, // je  addr   0x0f 0x84     0x12 0x01 0x00 0x00
	}

	bufLen := len(buffer)
	for i := offset; i < bufLen; i++ {
		for _, ins := range m {
			code := ins["code"]
			insLen := int(ins["insLen"][0])
			if i+insLen > bufLen {
				continue
			}

			ins := buffer[i : i+len(code)]
			addr := buffer[i+len(code) : i+insLen]
			if code[0] == ins[0] && (len(code) == 1 || code[1] == ins[1]) {
				v := int(bytes2Int(addr))
				if v < 0 {
					continue
				}

				v += insLen + i
				if isCallQ(ptr + uintptr(v)) {
					return i + len(code), len(addr), v, true
				}
			}
		}
	}

	return 0, 0, 0, false
}

// Find the position of first sub instruction in function header
//  return all the code between the function header and the instruction.
//  find the sub instruction by it's feature, there may be exceptions
//  need add more features if catch exceptions
func copyMoreStack(head []byte) []byte {
	var inArray = func(v byte, arr []byte) bool {
		for _, vv := range arr {
			if v == vv {
				return true
			}
		}
		return false
	}

	mid := []byte{0x81, 0x83, 0x8d}
	for _, v := range subFeatures {
		if head[v] == byte(0x48) && head[v+2] == byte(0xec) && inArray(head[v+1], mid) {
			return head[0:v]
		}
	}

	printRawData(head)
	panic("offset not found\n")

	return []byte{}
}

func bytes2Int(data []byte) int32 {
	d := uint32(0)

	for i, v := range data {
		x := uint(i) * 8
		d |= (uint32(v) << x) & (0xff << x)
	}

	return int32(d)
}

func int2Bytes(v int32) []byte {
	vv := uint32(v)
	return []byte{
		byte(vv & 0xff),
		byte((vv & 0xff00) >> 8),
		byte((vv & 0xff0000) >> 16),
		byte((vv & 0xff000000) >> 24),
	}
}

func isCallQ(addr uintptr) bool {
	f := rawMemoryAccess(addr, 1)
	return f[0] == 0xe8
}
