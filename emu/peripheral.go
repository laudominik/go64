package emu

type Peripheral interface {
	Read(reg uint64) uint32
	Write(reg uint64, value uint32)
}
