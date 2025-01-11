package emu

type Cpu struct {
	pc   uint64
	r    Registers
	cop0 Registers
}

type Registers [32]uint64

func (cpu *Cpu) reset() {
	/*
		we don't emulate PIF ROM, just its effects
	*/
	cpu.r[11] = 0xFFFFFFFFA4000040
	cpu.r[20] = 0x0000000000000001
	cpu.r[22] = 0x000000000000003F
	cpu.r[29] = 0xFFFFFFFFA4001FF0
	cpu.cop0[1] = 0x0000001F
	cpu.cop0[12] = 0x34000000
	cpu.cop0[15] = 0x00000B00
	cpu.cop0[16] = 0x0006E463
	cpu.pc = 0xA4000040
}

func sext32(num uint32, ogBits int) uint32 {
	signBit := uint32(1 << (ogBits - 1))
	if num&signBit != 0 {
		return uint32(num | ^((1 << ogBits) - 1))
	}
	return uint32(num & ((1 << ogBits) - 1))
}
