package peripherals

import "fmt"

type Pi struct {
	dramAddr uint32
	cartAddr uint32
	dmaBusy  uint32
	ioBusy   uint32
}

func (pi *Pi) Read(reg uint64) uint32 {
	switch reg {
	case 0x10:
		return (pi.dmaBusy | pi.ioBusy) /* todo: add PI interrupt */
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
	case 0xc:
		/* TODO: DMA request */
		panic(0)
	default:
		panic(fmt.Sprintf("Writing to unimplemented PI register 0x%x", reg))
	}
}
