package emu

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
	// for entry := range cpu.tlb {

	// }

	cpu.exception = true
	if isRead {
		cpu.exceptionCode = EXCEPTION_TLB_MISS_STORE
	} else {
		cpu.exceptionCode = EXCEPTION_TLB_MISS_LOAD
	}
	return 0
}
