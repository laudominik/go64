package peripherals

type Mi struct {
	miModeReg uint32
}

func (mi *Mi) Read(reg uint64) uint32 {
	if reg == 0 {
		return mi.miModeReg
	}
	return 0
}

func (mi *Mi) Write(reg uint64, value uint32) {
	if reg == 0 {
		mi.miModeReg = value
	}
}
