package mips

import "go64/emu/util"

var ISA_REGIMM_TABLE = map[uint32]InstructionCallback{
	0b000000: regimm_bltz,
	0b010001: regimm_bgezal,
}

var ISA_REGIMM_MNEMONIC = map[uint32]string{
	0b000000: "BLTZ",
	0b010001: "BGEZAL",
}

func regimm_bltz(cpu *Cpu, instr Instruction) {
	if int64(cpu.r[instr.rs]) >= 0 {
		return
	}
	cpu.planJump(cpu.pc + uint64(util.Sext32(instr.imm, 16))*4)
}

func regimm_bgezal(cpu *Cpu, instr Instruction) {
	if int64(cpu.r[instr.rs]) < 0 {
		return
	}
	cpu.r[31] = cpu.pc + 4
	cpu.planJump(cpu.pc + uint64(util.Sext32(instr.imm, 16))*4)
}
