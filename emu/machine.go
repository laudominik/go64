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
	rsp Rsp

	memoryMap []MemoryRange
}

type MemoryRange struct {
	start uint64
	end   uint64
	name  string
	p     peripherals.Peripheral
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

	if m.cpu.delaySlot.in {
		m.execute(instr)
		m.cpu.doJump()
	} else {
		m.execute(instr)
	}
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

	for _, memoryRange := range m.memoryMap {
		if !inRange(physicalAddress, memoryRange.start, memoryRange.end) {
			continue
		}
		if config.CONFIG.LogMemory.Read {
			fmt.Printf("Memory read (0x%x) from %s\n", physicalAddress, memoryRange.name)
		}
		return memoryRange.p.Read(physicalAddress - memoryRange.start)
	}

	panic(fmt.Sprintf("Reading unmapped memory 0x%x", physicalAddress))
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

	for _, memoryRange := range m.memoryMap {
		if !inRange(physicalAddress, memoryRange.start, memoryRange.end) {
			continue
		}
		if config.CONFIG.LogMemory.Write {
			fmt.Printf("Memory write (0x%x -> 0x%x) to %s\n", physicalAddress, value, memoryRange.name)
		}
		memoryRange.p.Write(physicalAddress-memoryRange.start, value)
		return
	}

	panic(fmt.Sprintf("Writing to unmapped memory 0x%x -> 0x%x", physicalAddress, value))
}

func (m *Machine) InitPeripherals() {
	m.memoryMap = []MemoryRange{
		MemoryRange{0x10000000, 0x1FBFFFFF, "Cardridge ROM", Memory{}}, // keep first
		MemoryRange{0x00000000, 0x003FFFFF, "RDRAM", make(Memory, 0x400000)},
		MemoryRange{0x03F00000, 0x03FFFFFF, "RDRAM MMIO", &peripherals.Unused{}},
		MemoryRange{0x04000000, 0x04000FFF, "RSP Data Memory", make(Memory, 0x1000)},
		MemoryRange{0x04001000, 0x04001FFF, "RSP Instruction Memory", make(Memory, 0x1000)},
		MemoryRange{0x04040000, 0x040FFFFF, "SP Registers", CreateSpRegs(m)},
		MemoryRange{0x04300000, 0x043FFFFF, "MIPS Interface", &peripherals.Mi{}},
		MemoryRange{0x04600000, 0x046FFFFF, "Peripheral Interface", &peripherals.Pi{}},
		MemoryRange{0x04700000, 0x047FFFFF, "RDRAM settings", &peripherals.Unused{}},
	}
}

func (m *Machine) Reset() {
	m.cpu.reset()
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

	m.memoryMap[0].p = Memory(data) // cardridge ROM
	fmt.Printf("Successfully loaded ROM: %s (%d B)\n", filePath, len(data))
	return nil
}
