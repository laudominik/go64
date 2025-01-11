package emu

var ISA_IJ_TABLE = map[uint32]InstructionCallback{
	0b100011: i_lw,
	0b101011: i_sw,
}

var ISA_IJ_MNEMONIC = map[uint32]string{
	0b100011: "LW",
	0b101011: "SW",
}

var ISA_IJ_TYPE = map[uint32]int{
	0b100011: INSTR_TYPE_I,
	0b101011: INSTR_TYPE_R,
}

func i_lw(m *Machine, instr Instruction) {
	v := m.readDWord(m.cpu.r[instr.rs] + uint64(instr.imm))
	m.cpu.r[instr.rt] = uint64(v)
}

func i_sw(m *Machine, instr Instruction) {
	v := m.cpu.r[instr.rt]
	m.writeDWord(m.cpu.r[instr.rs]+uint64(instr.imm), uint32(v))
}
