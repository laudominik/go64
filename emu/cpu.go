package emu

import "go64/emu/peripherals"

type Cpu struct {
	pc            uint64
	r             Registers
	hi            uint64
	lo            uint64
	cop0          Registers
	tlb           [32][2]uint64
	exception     bool
	exceptionCode int
	delaySlot     struct {
		/*
			this is supposed to imitate delay slots
			delay slot should occur every jump
		*/
		in        bool
		nextPCVal uint64
	}
	mi peripherals.Mi
}

type Registers [32]uint64

func (cpu *Cpu) reset() {
	/*
		we don't emulate PIF ROM, just its effects
	*/
	cpu.r[11] = 0xFFFFFFFFA4000040 // pointer to ipl3 bootcode
	cpu.r[29] = 0xFFFFFFFFA4001FF0 // SP
	cpu.r[31] = 0xFFFFFFFFA4001550 // Return Address
	cpu.r[2] = 0xFFFFFFFFF58B0FBF
	cpu.r[20] = 0x0000000000000001
	cpu.r[22] = 0x000000000000003F
	cpu.cop0[1] = 0x0000001F
	cpu.cop0[12] = 0x34000000
	cpu.cop0[15] = 0x00000B00
	cpu.cop0[16] = 0x0006E463
	cpu.pc = 0xA4000040
	cpu.mi.InterruptMask = 0b111111
}

func (cpu *Cpu) planJump(addr uint64) {
	cpu.delaySlot.nextPCVal = addr
	cpu.delaySlot.in = true
}

func (cpu *Cpu) doJump() {
	if !cpu.delaySlot.in {
		panic("Trying to jump outside the delay slot")
	}
	cpu.pc = cpu.delaySlot.nextPCVal
	cpu.delaySlot.in = false
	cpu.delaySlot.nextPCVal = 0
}
