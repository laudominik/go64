package emu

import (
	"errors"
	"fmt"
	"go64/config"
	"go64/emu/peripherals"
	"os"
)

type Machine struct {
	cpu Cpu

	/* memory map */
	cardridgeROM Memory         // Cartridge Domain 1 Address 2
	rspData      Memory         // SP DMEM
	rspInstr     Memory         // SP INSTR
	mi           peripherals.Mi // MIPS Interface (MI)
}

func inRange(arg, left, right uint64) bool {
	return arg >= left && arg <= right
}

func (m *Machine) Tick() {
	pc := m.cpu.pc
	if config.CONFIG.Pc {
		fmt.Printf("[%x] ", pc)
	}

	instrb := m.readDWord(pc)
	if m.cpu.exception {
		// exception can happen during instruction fetch
		m.cpu.handleException()
		return
	}
	instr := decode(instrb)
	m.execute(instr)
}

func (m *Machine) execute(instr Instruction) {
	m.cpu.pc += 4
	instr.callback(m, instr)
	if m.cpu.exception {
		m.cpu.handleException()
	}
	m.cpu.exception = false
}

func (m *Machine) readDWord(virtualAddress uint64) uint32 {
	if virtualAddress&0b11 != 0 {
		panic(fmt.Sprintf("Unaligned read at %x: %x", m.cpu.pc, virtualAddress))
	}
	virtualAddress &= 0xFFFFFFFF // todo: remove after implementing 64-bit mode

	physicalAddress := m.cpu.translate(virtualAddress, true)
	if m.cpu.exception {
		return 0
	}
	if inRange(physicalAddress, 0x10000000, 0x1FBFFFFF) {
		return m.cardridgeROM.Read(physicalAddress - 0x10000000)
	} else if inRange(physicalAddress, 0x04000000, 0x04000FFF) {
		return m.rspData.Read(physicalAddress - 0x04000000)
	} else if inRange(physicalAddress, 0x04001000, 0x04001FFF) {
		return m.rspInstr.Read(physicalAddress - 0x04001000)
	} else if inRange(physicalAddress, 0x04300000, 0x043FFFFF) {
		return m.mi.Read(physicalAddress)
	} else if inRange(physicalAddress, 0x04700000, 0x047FFFFF) {
		return 0 // Control RDRAM settings (timings?) Irrelevant for emulation.
	}
	panic("Reading unmapped memory")
}

func (m *Machine) writeDWord(virtualAddress uint64, value uint32) {
	if virtualAddress&0b11 != 0 {
		panic(fmt.Sprintf("Unaligned write at %x: %x", m.cpu.pc, virtualAddress))
	}
	virtualAddress &= 0xFFFFFFFF

	physicalAddress := m.cpu.translate(virtualAddress, false)
	if m.cpu.exception {
		return
	}
	if inRange(physicalAddress, 0x10000000, 0x1FBFFFFF) {
		panic("Trying to write to cardridge ROM")
	} else if inRange(physicalAddress, 0x04000000, 0x04000FFF) {
		m.rspData.Write(physicalAddress-0x04000000, value)
	} else if inRange(physicalAddress, 0x04001000, 0x04001FFF) {
		m.rspInstr.Write(physicalAddress-0x04001000, value)
	} else if inRange(physicalAddress, 0x04700000, 0x047FFFFF) {
		// Control RDRAM settings (timings?) Irrelevant for emulation.
	} else {
		panic("Writing to unmapped memory")
	}

}

func (m *Machine) Reset() {
	m.cpu.reset()
	m.rspData = make([]byte, 0x1000)
	m.rspInstr = make([]byte, 0x1000)
	/*
		The first 0x1000 bytes from the cartridge are copied to SP DMEM.
		This is implemented as a copy of 0x1000 bytes from 0xB0000000 to 0xA4000000.
	*/
	for i := uint64(0); i < 0x1000/4; i++ {
		m.writeDWord(0xA4000000+4*i, m.readDWord(0xB0000000+4*i))
	}
}

func (m *Machine) LoadRom(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("ROM file does not exist: %s", filePath)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read ROM file: %w", err)
	}
	if len(data) == 0 {
		return errors.New("ROM file is empty")
	}

	m.cardridgeROM = data
	fmt.Printf("Successfully loaded ROM: %s (%d B)\n", filePath, len(data))
	return nil
}
