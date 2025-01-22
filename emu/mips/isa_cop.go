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

var ISA_COP1_TABLE = map[uint32]InstructionCallback{
	0b00000: cop1_mfc1,
	0b00010: cop1_mfhc1,
	0b00100: cop1_mtc1,
}

var ISA_COP1_MNEMONIC = map[uint32]string{
	0b00000: "MFC1",
	0b00010: "MFHC1",
	0b00100: "MTC1",
}

var ISA_COP2_TABLE = map[uint32]InstructionCallback{
	0b00000: cop2_mfc2,
	0b00010: cop2_mfhc2,
	0b00100: cop2_mtc2,
}

var ISA_COP2_MNEMONIC = map[uint32]string{
	0b00000: "MFC2",
	0b00010: "MFHC2",
	0b00100: "MTC2",
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

func cop1_mtc1(cpu *Cpu, instr Instruction) {
	cpu.cop1[instr.rd+(instr.sa&0b111)] = cpu.r[instr.rt]
}

func cop1_mfc1(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rt] = cpu.cop1[instr.rd+(instr.sa&0b111)]
}

func cop1_mfhc1(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rt] = cpu.cop1[instr.rd+(instr.sa&0b111)] >> 32
}

func cop2_mtc2(cpu *Cpu, instr Instruction) {
	cpu.cop1[instr.rd+(instr.sa&0b111)] = cpu.r[instr.rt]
}

func cop2_mfc2(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rt] = cpu.cop1[instr.rd+(instr.sa&0b111)]
}

func cop2_mfhc2(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rt] = cpu.cop2[instr.rd+(instr.sa&0b111)] >> 32
}
