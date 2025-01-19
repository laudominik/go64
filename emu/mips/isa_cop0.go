package mips

var ISA_COP0_G1_TABLE = map[uint32]InstructionCallback{
	0b00100: cop0_g1_mtc0,
}

var ISA_COP0_G1_MNEMONIC = map[uint32]string{
	0b00100: "MTC0",
}

func cop0_g1_mtc0(cpu *Cpu, instr Instruction) {
	cpu.cop0[instr.rd+(instr.sa&0b111)] = cpu.r[instr.rt]
}
