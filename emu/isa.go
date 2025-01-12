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
const INSTR_TYPE_SPECIAL3 = 0b011111

func bits(num, end, start uint32) uint32 {
	mask := uint32((1 << (end - start + 1)) - 1)
	return (num >> start) & mask
}

var FUNCT_TABLE = map[int]map[uint32]InstructionCallback{
	0b000000: ISA_R_TABLE,
	0b011111: ISA_SPECIAL3_TABLE,
}
var MNEMONIC_TABLE = map[int]map[uint32]string{
	0b000000: ISA_R_MNEMONIC,
	0b011111: ISA_SPECIAL3_MNEMONIC,
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
	var ty int
	var validTy bool

	functTable, isExtension := FUNCT_TABLE[int(instr.opcode)]
	mnemonicTable, _ := MNEMONIC_TABLE[int(instr.opcode)]

	if isExtension {
		callback, valid = functTable[instr.funct]
		mnemonic, validMnemonic = mnemonicTable[instr.funct]
		ty = int(instr.opcode)
		validTy = true
	} else if instr.opcode == INSTR_TYPE_REGIMM {
		callback, valid = ISA_REGIMM_TABLE[instr.rt]
		mnemonic, validMnemonic = ISA_REGIMM_MNEMONIC[instr.rt]
		ty = int(instr.opcode)
		validTy = true
	} else {
		callback, valid = ISA_IJ_TABLE[instr.opcode]
		mnemonic, validMnemonic = ISA_IJ_MNEMONIC[instr.opcode]
		ty, validTy = ISA_IJ_TYPE[instr.opcode]
	}

	if !valid || !validMnemonic || !validTy {
		panic(fmt.Sprintf("Invalid instruction, Type: %d Opcode: 0b%b Funct: 0b%b", ty, instr.opcode, instr.funct))
	}

	instr.callback = callback
	instr.mnemonic = mnemonic
	instr.ty = ty

	if config.CONFIG.Disassemble {
		fmt.Println(instr.disassemble())
	}
	return instr
}

func (instr Instruction) disassemble() string {
	switch instr.ty {
	case INSTR_TYPE_I:
		return fmt.Sprintf("%s r%d r%d 0x%x", instr.mnemonic, instr.rs, instr.rd, instr.imm)
	case INSTR_TYPE_J:
		return fmt.Sprintf("%s 0x%x", instr.mnemonic, instr.tgt)
	case INSTR_TYPE_R, INSTR_TYPE_SPECIAL3:
		return fmt.Sprintf("%s r%d r%d r%d (shift=%d)", instr.mnemonic, instr.rs, instr.rt, instr.rd, instr.sa)
	case INSTR_TYPE_REGIMM:
		return fmt.Sprintf("%s r%d 0x%x", instr.mnemonic, instr.rs, instr.imm)

	}
	panic("Trying to disassemble invalid/undecoded function")
}

func stub(m *Machine, instr Instruction) {
	panic(fmt.Sprintf("Calling a stub %s", instr.mnemonic))
}
