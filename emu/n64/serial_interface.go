package n64

import (
	"fmt"
	"go64/emu/mips"
)

type Si struct {
	m *Machine
}

func CreateSi(m *Machine) *Si {
	var si Si
	si.m = m
	return &si
}

func (si *Si) Read(reg uint64) uint32 {
	switch reg {
	default:
		panic(fmt.Sprintf("Reading from unimplemented SI register 0x%x", reg))
	}
}

func (si *Si) Write(reg uint64, value uint32) {
	cpu := si.m.cpu
	switch reg {
	case 0x18:
		cpu.LowerInterrupt(mips.SI_MASK)
		break
	default:
		panic(fmt.Sprintf("Writing to unimplemented SI register 0x%x", reg))
	}
}
