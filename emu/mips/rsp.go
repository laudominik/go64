package mips

import "fmt"

type Rsp struct {
	pc      uint32
	semHeld bool
}

type SpRegs struct {
	rsp      *Rsp
	memAddr  uint32
	dramAddr uint32
	status   uint32
}

func CreateSpRegs(rsp *Rsp) *SpRegs {
	var sp SpRegs
	sp.rsp = rsp
	return &sp
}

func (sr *SpRegs) Read(reg uint64) uint32 {
	switch reg {
	case 0x0: /* SP Mem address */
		return sr.memAddr
	case 0x4: /* DRAM Address */
		return sr.dramAddr
	default:
		panic(fmt.Sprintf("Reading from unimplemented SP register 0x%x", reg))
	}

}

func (sr *SpRegs) Write(reg uint64, value uint32) {
	rsp := sr.rsp
	switch reg {
	case 0x0: /* SP Mem address */
		sr.memAddr = value
	case 0x4: /* DRAM Address */
		sr.dramAddr = value
	case 0x10: /* Status */
		sr.status = value
	case 0x1C: /* Semaphore */
		rsp.semHeld = false
	case 0x40000: /* PC */
		sr.rsp.pc = value >> 2
	default:
		panic(fmt.Sprintf("Writing to unimplemented SP register 0x%x", reg))
	}
}
