package peripherals

import (
	"fmt"
	"go64/emu/util"
)

type Mi struct {
	modeReg       uint32
	Interrupt     uint32
	InterruptMask uint32
}

const SP_BIT = 0
const SI_BIT = 1
const AI_BIT = 2
const VI_BIT = 3
const PI_BIT = 4
const DP_BIT = 5

var SP_MASK uint32 = util.Mask(SP_BIT)
var SI_MASK uint32 = util.Mask(SI_BIT)
var AI_MASK uint32 = util.Mask(AI_BIT)
var VI_MASK uint32 = util.Mask(VI_BIT)
var PI_MASK uint32 = util.Mask(PI_BIT)
var DP_MASK uint32 = util.Mask(DP_BIT)

func (mi *Mi) Read(reg uint64) uint32 {
	switch reg {
	case 0x0:
		return mi.modeReg
	case 0x4:
		return 0x02020102
	case 0x8:
		return mi.Interrupt
	case 0xC:
		return mi.InterruptMask
	default:
		panic(fmt.Sprintf("Reading from unimplemented MI register 0x%x", reg))
	}
}

func (mi *Mi) Write(reg uint64, value uint32) {
	switch reg {
	case 0x0:
		mi.modeReg = value // TODO: set vs clear bits
	case 0xC:
		mi.InterruptMask =
			(util.SetBit(value, 1, SP_BIT) & util.ClearBit(value, 0, SP_BIT)) |
				(util.SetBit(value, 3, SI_BIT) & util.ClearBit(value, 2, SI_BIT)) |
				(util.SetBit(value, 5, AI_BIT) & util.ClearBit(value, 4, AI_BIT)) |
				(util.SetBit(value, 7, VI_BIT) & util.ClearBit(value, 6, VI_BIT)) |
				(util.SetBit(value, 9, PI_BIT) & util.ClearBit(value, 8, PI_BIT)) |
				(util.SetBit(value, 11, DP_BIT) & util.ClearBit(value, 10, DP_BIT))

	default:
		panic(fmt.Sprintf("Writing to unimplemented/read-only MI register 0x%x", reg))
	}
}
