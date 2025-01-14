package emu

var ISA_R_TABLE = map[uint32]InstructionCallback{
	0b000000: r_sll,
	0b000010: r_srl,
	0b000011: stub,
	0b000100: stub,
	0b000110: stub,
	0b001000: r_jr,
	0b001001: stub,
	0b010000: r_mfhi,
	0b010001: stub,
	0b010010: stub,
	0b010011: stub,
	0b011000: stub,
	0b011001: stub,
	0b011010: stub,
	0b011011: stub,
	0b100000: r_add,
	0b100001: stub,
	0b100010: stub,
	0b100011: stub,
	0b100100: stub,
	0b100101: r_or,
	0b100110: r_xor,
	0b100111: stub,
	0b101010: stub,
	0b101011: stub,
}

var ISA_R_MNEMONIC = map[uint32]string{
	0b000000: "SLL",
	0b000010: "SRL",
	0b000011: "SRA",
	0b000100: "SLLV",
	0b000110: "SRLV",
	0b001000: "JR",
	0b001001: "JALR",
	0b010000: "MFHI",
	0b010001: "MTHI",
	0b010010: "MFLO",
	0b010011: "MTLO",
	0b011000: "MULT",
	0b011001: "MULTU",
	0b011010: "DIV",
	0b011011: "DIVU",
	0b100000: "ADD",
	0b100001: "ADDU",
	0b100010: "SUB",
	0b100011: "SUBU",
	0b100100: "AND",
	0b100101: "OR",
	0b100110: "XOR",
	0b100111: "NOR",
	0b101010: "SLT",
	0b101011: "SLTU",
}

func r_add(m *Machine, instr Instruction) {
	m.cpu.r[instr.rd] = m.cpu.r[instr.rs] + m.cpu.r[instr.rt]
}

func r_xor(m *Machine, instr Instruction) {
	m.cpu.r[instr.rd] = m.cpu.r[instr.rs] ^ m.cpu.r[instr.rt]
}

func r_mfhi(m *Machine, instr Instruction) {
	m.cpu.r[instr.rd] = m.cpu.hi
}

func r_sll(m *Machine, instr Instruction) {
	m.cpu.r[instr.rd] = m.cpu.r[instr.rt] << uint64(instr.sa)
}

func r_or(m *Machine, instr Instruction) {
	m.cpu.r[instr.rd] = m.cpu.r[instr.rs] | m.cpu.r[instr.rt]
}

func r_jr(m *Machine, instr Instruction) {
	m.cpu.planJump(m.cpu.r[instr.rs])
}

func r_srl(m *Machine, instr Instruction) {
	m.cpu.r[instr.rd] = m.cpu.r[instr.rt] >> instr.sa
}
