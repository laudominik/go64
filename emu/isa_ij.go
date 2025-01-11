package emu

var ISA_IJ_TABLE = map[uint32]InstructionCallback{
	0b100011: i_lw,
}

var ISA_IJ_MNEMONIC = map[uint32]string{
	0b100011: "LW",
}

var ISA_IJ_TYPE = map[uint32]int{
	0b100011: INSTR_TYPE_I,
}

func i_lw(m *Machine, instr Instruction) {
	v := m.readDWord(m.cpu.r[instr.rs] + uint64(instr.imm))
	m.cpu.r[instr.rt] = uint64(v)
}
