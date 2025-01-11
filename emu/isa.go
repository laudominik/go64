package emu

import "fmt"

type Instruction struct {
	/* decoded */
	mnemonic string
	ty       int
	emulate  func(m *Machine, instr Instruction)
	/* general */
	instrb uint32
	opcode uint32
	/* I-type */
	rs  uint32
	rt  uint32
	imm uint32
	/* R-type */
	rd    uint32
	sa    uint32
	funct uint32
	/* J-type */
	tgt uint32
}

const INSTR_TYPE_R = 0
const INSTR_TYPE_I = 1
const INSTR_TYPE_J = 2

func bits(num, end, start uint32) uint32 {
	mask := uint32((1 << (end - start + 1)) - 1)
	return (num >> start) & mask
}

func decode(instrb uint32) Instruction {
	instr := Instruction{
		mnemonic: "",
		ty:       -1,
		emulate:  nil,
		instrb:   instrb,
		opcode:   bits(instrb, 31, 26),
		rs:       bits(instrb, 25, 21),
		rt:       bits(instrb, 20, 16),
		imm:      bits(instrb, 15, 0),
		rd:       bits(instrb, 15, 11),
		sa:       bits(instrb, 10, 6),
		funct:    bits(instrb, 5, 0),
		tgt:      bits(instrb, 25, 0),
	}
	return instr
}

func (instr Instruction) disassemble() string {
	switch instr.ty {
	case INSTR_TYPE_I:
		return fmt.Sprintf("%s %d %d %d", instr.mnemonic, instr.rs, instr.rd, instr.imm)
	case INSTR_TYPE_J:
		return fmt.Sprintf("%s %d", instr.mnemonic, instr.tgt)
	case INSTR_TYPE_R:
		return fmt.Sprintf("%s %d %d %d %d", instr.mnemonic)
	}
	panic("Trying to disassemble invalid/undecoded function")
}
