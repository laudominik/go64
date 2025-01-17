package peripherals

type SpRegs struct {
}

func (sr *SpRegs) Read(reg uint64) uint32 {
	return 0
}

func (sr *SpRegs) Write(reg uint64, value uint32) {
	// TODO
	panic("Calling a stub write to sp regs")
}
