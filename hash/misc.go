package hash

import "unsafe"

func add(p unsafe.Pointer, amt int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + uintptr(amt))
}

func readUnaligned32(p unsafe.Pointer) uint32 {
	return *(*uint32)(p)
}

func readUnaligned64(p unsafe.Pointer) uint64 {
	return *(*uint64)(p)
}
