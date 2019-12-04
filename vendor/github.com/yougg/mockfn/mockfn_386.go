package mockfn

// Assembles a jump to a function value
func jmpToFunctionValue(to uintptr) []byte {
	return []byte{
		0xBA,
		byte(to),
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24), // mov edx,to
		0xFF, 0x22,     // jmp DWORD PTR [edx]
	}
}

func isAlreadyReplaced(bytes []byte) bool {
	return bytes[0] == 0xBA && bytes[5] == 0xFF && bytes[6] == 0x22
}
