package emu

var ISA_SPECIAL3_TABLE = map[uint32]InstructionCallback{
	0b000000: special3_ext,
}

var ISA_SPECIAL3_MNEMONIC = map[uint32]string{
	0b000000: "EXT",
}

func special3_ext(m *Machine, instr Instruction) {
	panic("Calling unimplemented EXT")
}
