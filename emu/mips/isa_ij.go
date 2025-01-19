package mips

import "go64/emu/util"

var ISA_IJ_TABLE = map[uint32]InstructionCallback{
	0b000010: stub,
	0b000011: j_jal,
	0b000100: i_beq,
	0b000101: i_bne,
	0b000110: stub,
	0b000111: stub,
	0b001000: i_addi,
	0b001001: i_addiu,
	0b001010: i_slti,
	0b001011: stub,
	0b001100: i_andi,
	0b001101: i_ori,
	0b001110: i_xori,
	0b001111: i_lui,
	0b010100: i_beql,
	0b010101: i_bnel,
	0b010110: i_blezl,
	0b011000: stub,
	0b011001: stub,
	0b011010: stub,
	0b100000: stub,
	0b100001: stub,
	0b100011: i_lw,
	0b100100: stub,
	0b100101: stub,
	0b100110: i_lwr,
	0b101000: stub,
	0b101001: stub,
	0b101011: i_sw,
	0b101111: i_cache,
}

var ISA_IJ_MNEMONIC = map[uint32]string{
	0b000010: "J",
	0b000011: "JAL",
	0b000100: "BEQ",
	0b000101: "BNE",
	0b000110: "BLEZ",
	0b000111: "BGTZ",
	0b001000: "ADDI",
	0b001001: "ADDIU",
	0b001010: "SLTI",
	0b001011: "SLTIU",
	0b001100: "ANDI",
	0b001101: "ORI",
	0b001110: "XORI",
	0b001111: "LUI",
	0b010100: "BEQL",
	0b010101: "BNEL",
	0b010110: "BLEZL",
	0b011000: "LLO",
	0b011001: "LHI",
	0b011010: "TRAP",
	0b100000: "LB",
	0b100001: "LH",
	0b100011: "LW",
	0b100100: "LBU",
	0b100101: "LHU",
	0b100110: "LWR",
	0b101000: "SB",
	0b101001: "SH",
	0b101011: "SW",
	0b101111: "CACHE",
}

var ISA_IJ_TYPE = map[uint32]int{
	0b000010: INSTR_TYPE_J,
	0b000011: INSTR_TYPE_J,
	0b000100: INSTR_TYPE_I,
	0b000101: INSTR_TYPE_I,
	0b000110: INSTR_TYPE_I,
	0b000111: INSTR_TYPE_I,
	0b001000: INSTR_TYPE_I,
	0b001001: INSTR_TYPE_I,
	0b001010: INSTR_TYPE_I,
	0b001011: INSTR_TYPE_I,
	0b001100: INSTR_TYPE_I,
	0b001101: INSTR_TYPE_I,
	0b001110: INSTR_TYPE_I,
	0b001111: INSTR_TYPE_I,
	0b010100: INSTR_TYPE_I,
	0b010101: INSTR_TYPE_I,
	0b010110: INSTR_TYPE_I,
	0b011000: INSTR_TYPE_I,
	0b011001: INSTR_TYPE_I,
	0b011010: INSTR_TYPE_J,
	0b100000: INSTR_TYPE_I,
	0b100001: INSTR_TYPE_I,
	0b100011: INSTR_TYPE_I,
	0b100100: INSTR_TYPE_I,
	0b100101: INSTR_TYPE_I,
	0b100110: INSTR_TYPE_I,
	0b101000: INSTR_TYPE_I,
	0b101001: INSTR_TYPE_I,
	0b101011: INSTR_TYPE_I,
	0b101111: INSTR_TYPE_I,
}

func i_bne(cpu *Cpu, instr Instruction) {
	if cpu.r[instr.rs] == cpu.r[instr.rt] {
		return
	}
	cpu.planJump(cpu.pc + uint64(util.Sext32(instr.imm, 16))*4)
}

func i_beq(cpu *Cpu, instr Instruction) {
	if cpu.r[instr.rs] != cpu.r[instr.rt] {
		return
	}
	cpu.planJump(cpu.pc + uint64(util.Sext32(instr.imm, 16))*4)
}

func i_addi(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rt] = cpu.r[instr.rs] + util.Sext64(uint64(instr.imm), 16)
}

func i_andi(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rt] = cpu.r[instr.rs] & uint64(instr.imm)
}

func i_lw(cpu *Cpu, instr Instruction) {
	se := util.Sext64(uint64(instr.imm), 16)
	addr := cpu.r[instr.rs] + se
	v := cpu.AddressSpace.Read(addr)
	if cpu.Exception { // exception when reading
		return
	}
	cpu.r[instr.rt] = uint64(v)
}

func i_sw(cpu *Cpu, instr Instruction) {
	v := cpu.r[instr.rt]
	addr := cpu.r[instr.rs] + util.Sext64(uint64(instr.imm), 16)
	cpu.AddressSpace.Write(addr, uint32(v))
	// exception when writing can happen so this instruction won't have any effect
}

func i_lwr(cpu *Cpu, instr Instruction) {
	addr := cpu.r[instr.rs] + uint64(util.Sext32(instr.imm, 16))

	offset := addr % 4 // byte number in word
	aligned := addr - (addr % 4)

	val := cpu.AddressSpace.Read(aligned)
	if cpu.Exception { // exception when reading
		return
	}

	shift := offset * 8
	mask := (1 << (32 - shift)) - 1
	extracted := (val & uint32(mask)) << shift
	cpu.r[instr.rt] = (cpu.r[instr.rt] & ^uint64(mask<<shift)) | uint64(extracted)
}

func i_lui(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rt] = uint64(instr.imm << 16)
}

func i_addiu(cpu *Cpu, instr Instruction) {
	/*
		ADDIU performs the same arithmetic operation but, does not trap on overflow
		should maybe consider that... but let's keep it simple for now
	*/
	i_addi(cpu, instr)
}

func i_ori(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rt] = cpu.r[instr.rs] | uint64(instr.imm)
}

func j_jal(cpu *Cpu, instr Instruction) {
	cpu.r[31] = cpu.pc
	cpu.planJump((cpu.pc & 0xFFFFFFFFF0000000) + uint64(instr.tgt<<2))
}

func i_slti(cpu *Cpu, instr Instruction) {
	if int64(cpu.r[instr.rs]) < int64(instr.imm) {
		cpu.r[instr.rt] = 1
		return
	}
	cpu.r[instr.rt] = 0
}

func i_beql(cpu *Cpu, instr Instruction) {
	/* 	how does it differ from beq?
	I suppose it is used in branch predictor
	so from emulation POV no difference */
	if cpu.r[instr.rs] != cpu.r[instr.rt] {
		return
	}
	cpu.planJump(cpu.pc + uint64(util.Sext32(instr.imm, 16))*4)
}

func i_xori(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rt] = cpu.r[instr.rs] ^ uint64(instr.imm)
}

func i_bnel(cpu *Cpu, instr Instruction) {
	i_bne(cpu, instr)
}

func i_blezl(cpu *Cpu, instr Instruction) {
	if int64(cpu.r[instr.rs]) > 0 {
		return
	}
	cpu.planJump(cpu.pc + uint64(util.Sext32(instr.imm, 16))*4)
}

func i_cache(cpu *Cpu, instr Instruction) {
	/* Not needed
	not 100% sure so sth to keep in mind */
}
