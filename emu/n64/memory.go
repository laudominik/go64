package n64

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
