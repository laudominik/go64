package mips

import (
	"fmt"
	"go64/config"
	"go64/emu/util"
)

type Cpu struct {
	pc            uint64
	r             Registers
	hi            uint64
	lo            uint64
	cop0          Registers
	cop1          Registers
	cop2          Registers
	tlb           [32][2]uint64
	Exception     bool
	exceptionCode int
	delaySlot     struct {
		/*
			this is supposed to imitate delay slots
			delay slot should occur every jump
		*/
		in        bool
		nextPCVal uint64
	}
	Mi           Mi
	AddressSpace util.Peripheral
}

type Registers [32]uint64

func (cpu *Cpu) Reset() {
	/*
		we don't emulate PIF ROM, just its effects
	*/
	cpu.r[11] = 0xFFFFFFFFA4000040 // pointer to ipl3 bootcode
	cpu.r[29] = 0xFFFFFFFFA4001FF0 // SP
	cpu.r[31] = 0xFFFFFFFFA4001550 // Return Address
	cpu.r[2] = 0xFFFFFFFFF58B0FBF
	cpu.r[20] = 0x0000000000000001
	cpu.r[22] = 0x91 // it should be checksum seed
	cpu.cop0[1] = 0x0000001F
	cpu.cop0[12] = 0x34000000
	cpu.cop0[15] = 0x00000B00
	cpu.cop0[16] = 0x0006E463
	cpu.pc = 0xA4000040
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

func (cpu *Cpu) Tick() {
	pc := cpu.pc
	if config.CONFIG.Pc {
		fmt.Printf("[%x] ", pc)
	}

	instrb := cpu.AddressSpace.Read(pc)
	// fmt.Println(cpu.r[22])
	if cpu.Exception {
		// exception can happen during instruction fetch
		cpu.handleException()
		return
	}
	instr := decode(instrb)

	if cpu.delaySlot.in {
		cpu.execute(instr)
		cpu.doJump()
	} else {
		cpu.execute(instr)
	}
}

func (cpu *Cpu) execute(instr Instruction) {
	cpu.pc += 4
	instr.callback(cpu, instr)
	if cpu.Exception {
		cpu.handleException()
	}
	cpu.Exception = false
}

func (cpu *Cpu) Translate(virtualAddress uint64, isRead bool) uint64 {
	/* KSEG0 */
	if 0x80000000 <= virtualAddress && virtualAddress < 0xA0000000 {
		return virtualAddress - 0x80000000
	}
	/* KSEG1 */
	if 0xA0000000 <= virtualAddress && virtualAddress < 0xC0000000 {
		return virtualAddress - 0xA0000000
	}
	/* KUSEG, KSSEG, KSEG3 */
	return cpu.tlbTranslate(virtualAddress, isRead)
}

func (cpu *Cpu) tlbTranslate(virtualAddress uint64, isRead bool) uint64 {
	if isRead {
		cpu.raiseException(EXCEPTION_TLB_MISS_STORE)
	} else {
		cpu.raiseException(EXCEPTION_TLB_MISS_LOAD)
	}

	return 0
}

func (cpu *Cpu) Pc() uint64 {
	return cpu.pc
}
