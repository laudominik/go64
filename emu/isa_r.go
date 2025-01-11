package emu

var ISA_R_TABLE = map[uint32]InstructionCallback{
	0b100000: r_add,
}

var ISA_R_MNEMONIC = map[uint32]string{
	0b100000: "ADD",
}

func r_add(m *Machine, instr Instruction) {
	m.cpu.r[instr.rd] = m.cpu.r[instr.rs] + m.cpu.r[instr.rt]
}
