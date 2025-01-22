package mips

var ISA_COP0_G1_TABLE = map[uint32]InstructionCallback{
	0b00000: cop0_g1_mfc0,
	0b00010: cop0_g1_mfhc0,
	0b00100: cop0_g1_mtc0,
}

var ISA_COP0_G1_MNEMONIC = map[uint32]string{
	0b00000: "MFC0",
	0b00010: "MFHC0",
	0b00100: "MTC0",
}

func cop0_g1_mtc0(cpu *Cpu, instr Instruction) {
	cpu.cop0[instr.rd+(instr.sa&0b111)] = cpu.r[instr.rt]
}

func cop0_g1_mfc0(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rt] = cpu.cop0[instr.rd+(instr.sa&0b111)]
}

func cop0_g1_mfhc0(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rt] = cpu.cop0[instr.rd+(instr.sa&0b111)] >> 32
}
