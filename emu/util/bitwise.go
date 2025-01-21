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

func Bits(num, end, start uint32) uint32 {
	mask := uint32((1 << (end - start + 1)) - 1)
	return (num >> start) & mask
}

func Reg32(num uint64) uint32 {
	return uint32(num & 0xFFFFFFFF)
}

func Bit(num, no uint32) uint32 {
	return Bits(num, no, no)
}

// if bit inputBit is set, return number with outputBit set to 1
func SetBit(value uint32, inputBit uint32, outputBit uint32) uint32 {
	if value&(1<<inputBit) != 0 {
		return 1 << outputBit
	}
	return 0
}

// if bit inputBit is set, return number with outputBit set to 0
func ClearBit(value uint32, inputBit uint32, outputBit uint32) uint32 {
	if value&(1<<inputBit) != 0 {
		return ^(1 << outputBit)
	}
	return 0xffffffff
}

func Mask(ithBit uint32) uint32 {
	return (1 << ithBit)
}
