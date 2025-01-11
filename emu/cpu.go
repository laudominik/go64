package emu

type Cpu struct {
	pc            uint64
	r             Registers
	cop0          Registers
	tlb           [32][2]uint64
	exception     bool
	exceptionCode int
}

type Registers [32]uint64

type Exceptions struct {
	tlb   bool
	tlb64 bool
}

func (cpu *Cpu) reset() {
	/*
		we don't emulate PIF ROM, just its effects
	*/
	cpu.r[11] = 0xFFFFFFFFA4000040
	cpu.r[20] = 0x0000000000000001
	cpu.r[22] = 0x000000000000003F
	cpu.r[29] = 0xFFFFFFFFA4001FF0
	cpu.cop0[1] = 0x0000001F
	cpu.cop0[12] = 0x34000000
	cpu.cop0[15] = 0x00000B00
	cpu.cop0[16] = 0x0006E463
	cpu.pc = 0xA4000040
}

func sext32(num uint32, ogBits int) uint32 {
	signBit := uint32(1 << (ogBits - 1))
	if num&signBit != 0 {
		return uint32(num | ^((1 << ogBits) - 1))
	}
	return uint32(num & ((1 << ogBits) - 1))
}

const COP0_STATUS = 12
const COP0_CAUSE = 13
const COP0_EPC = 14

const EXCEPTION_INTERRUPT = 0
const EXCEPTION_TLB_MISS_LOAD = 2
const EXCEPTION_TLB_MISS_STORE = 3
const EXCEPTION_COP_UNUSABLE = 11

const STATUS_EXL = 1
const STATUS_ERL = 2
const STATUS_UX = 5
const STATUS_SX = 6
const STATUS_KX = 7
const STATUS_BEV = 22

func (cpu *Cpu) handleException() {
	// cpu.cop0[/* Cause */] = 0
	oldEXL := cpu.cop0[COP0_STATUS] & STATUS_EXL

	if cpu.cop0[COP0_STATUS]&STATUS_EXL == 0 {
		cpu.cop0[COP0_STATUS] |= STATUS_EXL
		cpu.cop0[COP0_EPC] = cpu.pc - 4
	}
	cpu.cop0[COP0_CAUSE] = uint64(cpu.exceptionCode) << 2

	if cpu.exceptionCode == EXCEPTION_COP_UNUSABLE {
		// set the coprocessor error field in $Cause to the coprocessor that caused the error
	}

	if cpu.cop0[COP0_STATUS]&STATUS_BEV == 0 {
		// late boot stage exception
		if cpu.exceptionCode == EXCEPTION_TLB_MISS_LOAD || cpu.exceptionCode == EXCEPTION_TLB_MISS_STORE {
			if oldEXL == 0 {
				cpu.pc = 0x80000000
			} else {
				cpu.pc = 0x80000180
			}
		} else {
			cpu.pc = 0x80000180
		}
	} else {
		// early boot stage exception
		if cpu.exceptionCode == EXCEPTION_TLB_MISS_LOAD || cpu.exceptionCode == EXCEPTION_TLB_MISS_STORE {
			if oldEXL == 0 {
				cpu.pc = 0xBFC00200
			} else {
				cpu.pc = 0xBFC00380
			}
		} else {
			cpu.pc = 0xBFC00380
		}
	}
}
