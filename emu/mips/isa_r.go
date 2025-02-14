package mips

import (
	"go64/emu/util"
	b "math/bits"
)

var ISA_R_TABLE = map[uint32]InstructionCallback{
	0b000000: r_sll,
	0b000010: r_srl,
	0b000011: stub,
	0b000100: r_sllv,
	0b000110: r_srlv,
	0b001000: r_jr,
	0b001001: stub,
	0b010000: r_mfhi,
	0b010001: stub,
	0b010010: r_mflo,
	0b010011: stub,
	0b011000: r_mult,
	0b011001: r_mult,
	0b011010: stub,
	0b011011: stub,
	0b100000: r_add,
	0b100001: r_add,
	0b100010: r_sub,
	0b100011: r_sub,
	0b100100: r_and,
	0b100101: r_or,
	0b100110: r_xor,
	0b100111: stub,
	0b101010: r_slt,
	0b101011: r_sltu,
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

func r_mult(cpu *Cpu, instr Instruction) {
	cpu.hi, cpu.lo = b.Mul64(cpu.r[instr.rs], cpu.r[instr.rt])
}

func r_add(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rd] = cpu.r[instr.rs] + cpu.r[instr.rt]
}

func r_sub(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rd] = cpu.r[instr.rs] - cpu.r[instr.rt]
}

func r_xor(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rd] = cpu.r[instr.rs] ^ cpu.r[instr.rt]
}

func r_mfhi(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rd] = cpu.hi
}

func r_mflo(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rd] = cpu.lo
}

func r_sll(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rd] = cpu.r[instr.rt] << uint64(instr.sa)
}

func r_or(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rd] = cpu.r[instr.rs] | cpu.r[instr.rt]
}

func r_and(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rd] = cpu.r[instr.rs] & cpu.r[instr.rt]
}

func r_jr(cpu *Cpu, instr Instruction) {
	cpu.planJump(cpu.r[instr.rs])
}

func r_srl(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rd] = cpu.r[instr.rt] >> instr.sa
}

func r_slt(cpu *Cpu, instr Instruction) {
	if int64(cpu.r[instr.rs]) < int64(cpu.r[instr.rt]) {
		cpu.r[instr.rd] = 1
		return
	}
	cpu.r[instr.rd] = 0
}

func r_sltu(cpu *Cpu, instr Instruction) {
	if util.Reg32(cpu.r[instr.rs]) < util.Reg32(cpu.r[instr.rt]) {
		cpu.r[instr.rd] = 1
		return
	}
	cpu.r[instr.rd] = 0
}

func r_srlv(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rd] = cpu.r[instr.rt] >> (cpu.r[instr.rs] & 0b11111)
}

func r_sllv(cpu *Cpu, instr Instruction) {
	cpu.r[instr.rd] = cpu.r[instr.rt] << (cpu.r[instr.rs] & 0b11111)
}
