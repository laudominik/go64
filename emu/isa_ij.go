package emu

var ISA_IJ_TABLE = map[uint32]InstructionCallback{
	0b000101: i_bne,
	0b001000: i_addi,
	0b001100: i_andi,
	0b100011: i_lw,
	0b101011: i_sw,
}

var ISA_IJ_MNEMONIC = map[uint32]string{
	0b000101: "BNE",
	0b001000: "ADDI",
	0b001100: "ANDI",
	0b100011: "LW",
	0b101011: "SW",
}

var ISA_IJ_TYPE = map[uint32]int{
	0b000101: INSTR_TYPE_I,
	0b001000: INSTR_TYPE_I,
	0b001100: INSTR_TYPE_I,
	0b100011: INSTR_TYPE_I,
	0b101011: INSTR_TYPE_I,
}

func i_bne(m *Machine, instr Instruction) {
	if m.cpu.r[instr.rs] == m.cpu.r[instr.rt] {
		return
	}
	m.cpu.pc += uint64(sext32(instr.imm, 16)) * 4
}

func i_addi(m *Machine, instr Instruction) {
	m.cpu.r[instr.rt] = m.cpu.r[instr.rs] + uint64(sext32(instr.imm, 16))
}

func i_andi(m *Machine, instr Instruction) {
	m.cpu.r[instr.rt] = m.cpu.r[instr.rs] & uint64(instr.imm)
}

func i_lw(m *Machine, instr Instruction) {
	addr := m.cpu.r[instr.rs] + uint64(sext32(instr.imm, 16))
	v := m.readDWord(addr)
	if m.cpu.exception { // exception when reading
		return
	}
	m.cpu.r[instr.rt] = uint64(v)
}

func i_sw(m *Machine, instr Instruction) {
	v := m.cpu.r[instr.rt]
	addr := m.cpu.r[instr.rs] + uint64(sext32(instr.imm, 16))
	m.writeDWord(addr, uint32(v))
	// exception when writing can happen so this instruction won't have any effect
}
