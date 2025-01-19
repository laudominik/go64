package mips

import "go64/emu/util"

var ISA_REGIMM_TABLE = map[uint32]InstructionCallback{
	0b000000: regimm_bltz,
}

var ISA_REGIMM_MNEMONIC = map[uint32]string{
	0b000000: "BLTZ",
}

func regimm_bltz(cpu *Cpu, instr Instruction) {
	if int64(cpu.r[instr.rs]) >= 0 {
		return
	}
	cpu.pc += uint64(util.Sext32(instr.imm, 16)) * 4
}
