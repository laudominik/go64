package emu

type Memory []byte

func (mem Memory) Read(reg uint64) uint32 {
	hh := uint32(mem[reg])
	hl := uint32(mem[reg+1])
	lh := uint32(mem[reg+2])
	ll := uint32(mem[reg+3])
	return ll + (lh << 8) + (hl << 16) + (hh << 24)
}

func (mem Memory) Write(reg uint64, value uint32) {
	mem[reg] = byte((value >> 24) & 0xFF)
	mem[reg+1] = byte((value >> 16) & 0xFF)
	mem[reg+2] = byte((value >> 8) & 0xFF)
	mem[reg+3] = byte(value & 0xFF)
}

func (cpu *Cpu) translate(virtualAddress uint64, isRead bool) uint64 {
	/* KSEG0 */
	if 0x80000000 <= virtualAddress && virtualAddress < 0xA0000000 {
		return virtualAddress - 0x80000000
	}
	/* KSEG1 */
	if 0xA0000000 <= virtualAddress && virtualAddress < 0xC0000000 {
		return virtualAddress - 0xA0000000
	}
	/* KUSEG, KSSEG, KSEG3 */
	return cpu.tlbTranslate(virtualAddress, isRead)
}

func (cpu *Cpu) tlbTranslate(virtualAddress uint64, isRead bool) uint64 {
	if isRead {
		cpu.raiseException(EXCEPTION_TLB_MISS_STORE)
	} else {
		cpu.raiseException(EXCEPTION_TLB_MISS_LOAD)
	}

	return 0
}
