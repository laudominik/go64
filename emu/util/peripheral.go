package util

type Peripheral interface {
	Read(reg uint64) uint32
	Write(reg uint64, value uint32)
}

type MemoryRange struct {
	Start uint64
	End   uint64
	Name  string
	P     Peripheral
}
