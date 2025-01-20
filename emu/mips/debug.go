package mips

import "fmt"

var REG_ALT_NAME = [32]string{
	"zero",
	"at",
	"v0",
	"v1",
	"a0", //r4
	"a1",
	"a2",
	"a3",
	"t0", //r8
	"t1",
	"t2",
	"t3",
	"t4",
	"t5",
	"t6",
	"t7",
	"s0", //r16
	"s1",
	"s2",
	"s3",
	"s4",
	"s5",
	"s6",
	"s7",
	"t8",
	"t9",
	"k0",
	"k1",
	"gp",
	"sp",
	"s8",
	"ra",
}

func (instr Instruction) disassemble() string {
	switch instr.ty {
	case INSTR_TYPE_I:
		return fmt.Sprintf("%s %s %s 0x%x", instr.mnemonic, REG_ALT_NAME[instr.rs], REG_ALT_NAME[instr.rt], instr.imm)
	case INSTR_TYPE_J:
		return fmt.Sprintf("%s 0x%x", instr.mnemonic, instr.tgt)
	case INSTR_TYPE_R:
		return fmt.Sprintf("%s %s %s %s (shift=%d)", instr.mnemonic, REG_ALT_NAME[instr.rs], REG_ALT_NAME[instr.rt], REG_ALT_NAME[instr.rd], instr.sa)
	case INSTR_TYPE_REGIMM:
		return fmt.Sprintf("%s %s 0x%x", instr.mnemonic, REG_ALT_NAME[instr.rs], instr.imm)
	case INSTR_TYPE_COP0:
		return fmt.Sprintf("%s %s %s %d", instr.mnemonic, REG_ALT_NAME[instr.rt], REG_ALT_NAME[instr.rd], instr.sa)
	}
	panic("Trying to disassemble invalid/undecoded function")
}

func (cpu *Cpu) dump_regs() {
	for i, r := range cpu.r {
		fmt.Println(fmt.Sprintf("%s = 0x%x", REG_ALT_NAME[i], r))
	}
}
