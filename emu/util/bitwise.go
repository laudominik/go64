package util

func Sext32(num uint32, ogBits int) uint32 {
	signBit := uint32(1 << (ogBits - 1))
	if num&signBit != 0 {
		return uint32(num | ^((1 << ogBits) - 1))
	}
	return uint32(num & ((1 << ogBits) - 1))
}

func Sext64(num uint64, ogBits int) uint64 {
	signBit := uint64(1 << (ogBits - 1))
	if num&signBit != 0 {
		return uint64(num | ^((1 << ogBits) - 1))
	}
	return uint64(num & ((1 << ogBits) - 1))
}
