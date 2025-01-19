package emu

import (
	"fmt"
	"go64/emu/util"
)

type Pi struct {
	m        *Machine
	dramAddr uint32
	cartAddr uint32
	len      uint32
	dmaBusy  uint32
	ioBusy   uint32
}

func CreatePi(m *Machine) *Pi {
	var pi Pi
	pi.m = m
	return &pi
}

func (pi *Pi) Read(reg uint64) uint32 {
	m := pi.m
	switch reg {
	case 0x10:
		piInterrupt := util.Bit(m.cpu.mi.Interrupt, PI_BIT)
		return (pi.dmaBusy | pi.ioBusy | piInterrupt<<3)
	default:
		panic(fmt.Sprintf("Reading from unimplemented PI register 0x%x", reg))
	}
}

func (pi *Pi) Write(reg uint64, value uint32) {
	switch reg {
	case 0x0: /* DRAM address */
		pi.dramAddr = value
	case 0x4: /* Cartridge address */
		pi.cartAddr = value
	case 0xc: /* RD Length */
		len := (value & 0x00FFFFFF) + 1
		cartAddr := pi.cartAddr & 0xFFFFFFFE
		dramAddr := pi.dramAddr & 0x007FFFFE

		if (dramAddr&0x7 != 0) && (len >= 0x8) {
			len -= dramAddr & 0x7
		}

		if dramAddr > 0x00800000 {
			panic("Should never happen: DMA address to high")
		}

		pi.len = len
		pi.doDmaTransfer(cartAddr, dramAddr, len)

	default:
		panic(fmt.Sprintf("Writing to unimplemented PI register 0x%x", reg))
	}
}

func (pi *Pi) doDmaTransfer(from uint32, to uint32, len uint32) {
	if len%4 != 0 || (from&0b11 != 0) || (to&0b11 != 0) {
		panic("DMA trying to transfer unaligned address")
	}

	m := pi.m
	for i := uint32(0); i < len; i += 4 {
		m.writeDWordPhys(uint64(to+i), m.readDwordPhys(uint64(from+i)))
	}

	m.cpu.raiseInterrupt(PI_MASK)
}
