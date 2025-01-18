package emu

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

func (cpu *Cpu) raiseInterrupt() {
	cpu.exception = true
	cpu.exceptionCode = EXCEPTION_INTERRUPT
}

func (cpu *Cpu) raiseException(code int) {
	cpu.exception = true
	cpu.exceptionCode = code
}

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
		panic("Early boot stage exception, not implemented")
	}
}
