package emu

func translate(virtualAddress uint64) uint64 {
	/* KSEG0 */
	if 0x80000000 <= virtualAddress && virtualAddress < 0xA0000000 {
		return virtualAddress - 0x80000000
	}
	/* KSEG1 */
	if 0xA0000000 <= virtualAddress && virtualAddress < 0xC0000000 {
		return virtualAddress - 0xA0000000
	}
	/* KUSEG, KSSEG, KSEG3 */
	return tlbTranslate()
}

func tlbTranslate() uint64 {
	// TODO
	return 0
}
