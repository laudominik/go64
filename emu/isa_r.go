package emu

var ISA_R_TABLE = map[uint32]InstructionCallback{
	0x20: r_add,
	0x22: stub,
	0x00: stub,
}

var ISA_R_MNEMONIC = map[uint32]string{
	0x20: "ADD",
}

func r_add(m *Machine, instr Instruction) {
	m.cpu.r[instr.rd] = m.cpu.r[instr.rs] + m.cpu.r[instr.rt]
}
