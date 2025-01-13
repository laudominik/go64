package emu

import (
	"fmt"
	"go64/config"
)

type Instruction struct {
	/* decoded */
	mnemonic string
	ty       int
	callback InstructionCallback
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

type InstructionCallback func(m *Machine, instr Instruction)

const INSTR_TYPE_R = 0b000000
const INSTR_TYPE_I = -1
const INSTR_TYPE_J = -2
const INSTR_TYPE_REGIMM = 0b000001
const INSTR_TYPE_COP0 = 0b010000

func bits(num, end, start uint32) uint32 {
	mask := uint32((1 << (end - start + 1)) - 1)
	return (num >> start) & mask
}

func decode(instrb uint32) Instruction {
	instr := Instruction{
		mnemonic: "",
		ty:       -1,
		callback: nil,
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

	var callback InstructionCallback
	var mnemonic string
	var valid bool
	var validMnemonic bool
	ty := int(instr.opcode)
	validTy := true

	switch instr.opcode {
	case INSTR_TYPE_R:
		callback, valid = ISA_R_TABLE[instr.funct]
		mnemonic, validMnemonic = ISA_R_MNEMONIC[instr.funct]
	case INSTR_TYPE_REGIMM:
		callback, valid = ISA_REGIMM_TABLE[instr.rt]
		mnemonic, validMnemonic = ISA_REGIMM_MNEMONIC[instr.rt]
	case INSTR_TYPE_COP0:
		if instr.rs&0b10000 == 0 {
			callback, valid = ISA_COP0_G1_TABLE[instr.rs]
			mnemonic, validMnemonic = ISA_COP0_G1_MNEMONIC[instr.rs]
		} else {

		}
	default:
		callback, valid = ISA_IJ_TABLE[instr.opcode]
		mnemonic, validMnemonic = ISA_IJ_MNEMONIC[instr.opcode]
		ty, validTy = ISA_IJ_TYPE[instr.opcode]
	}

	if !valid || !validMnemonic || !validTy {
		panic(fmt.Sprintf("Invalid instruction, Opcode: 0b%b Funct: 0b%b", instr.opcode, instr.funct))
	}

	instr.callback = callback
	instr.mnemonic = mnemonic
	instr.ty = ty

	if config.CONFIG.Disassemble {
		fmt.Println(instr.disassemble())
	}
	return instr
}

func stub(m *Machine, instr Instruction) {
	panic(fmt.Sprintf("Calling a stub %s", instr.mnemonic))
}
