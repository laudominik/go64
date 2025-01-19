package n64

import (
	"errors"
	"fmt"
	"go64/config"
	"go64/emu/mips"
	"go64/emu/util"
	"os"
)

type Machine struct {
	cpu mips.Cpu
	rsp mips.Rsp

	memoryMap    []util.MemoryRange
	cartridgeRom Memory
}

func (m *Machine) Tick() {
	m.cpu.Tick()
}

func inRange(arg, left, right uint64) bool {
	return arg >= left && arg <= right
}

func (m *Machine) Read(virtualAddress uint64) uint32 {
	if virtualAddress&0b11 != 0 {
		panic(fmt.Sprintf("Unaligned read at %x: %x", m.cpu.Pc(), virtualAddress))
	}
	virtualAddress &= 0xFFFFFFFF // todo: remove after implementing 64-bit mode

	physicalAddress := m.cpu.Translate(virtualAddress, true)
	if m.cpu.Exception {
		return 0
	}
	return m.readDwordPhys(physicalAddress)
}

func (m *Machine) readDwordPhys(physicalAddress uint64) uint32 {
	for _, memoryRange := range m.memoryMap {
		if !inRange(physicalAddress, memoryRange.Start, memoryRange.End) {
			continue
		}
		if config.CONFIG.LogMemory.Read {
			fmt.Printf("Memory read (0x%x) from %s\n", physicalAddress, memoryRange.Name)
		}
		return memoryRange.P.Read(physicalAddress - memoryRange.Start)
	}

	panic(fmt.Sprintf("Reading unmapped memory 0x%x", physicalAddress))
}

func (m *Machine) Write(virtualAddress uint64, value uint32) {
	if virtualAddress&0b11 != 0 {
		panic(fmt.Sprintf("Unaligned write at %x: %x", m.cpu.Pc(), virtualAddress))
	}
	virtualAddress &= 0xFFFFFFFF

	physicalAddress := m.cpu.Translate(virtualAddress, false)
	if m.cpu.Exception {
		return
	}
	m.writeDWordPhys(physicalAddress, value)
}

func (m *Machine) writeDWordPhys(physicalAddress uint64, value uint32) {
	for _, memoryRange := range m.memoryMap {
		if !inRange(physicalAddress, memoryRange.Start, memoryRange.End) {
			continue
		}
		if config.CONFIG.LogMemory.Write {
			fmt.Printf("Memory write (0x%x -> 0x%x) to %s\n", physicalAddress, value, memoryRange.Name)
		}
		memoryRange.P.Write(physicalAddress-memoryRange.Start, value)
		return
	}

	panic(fmt.Sprintf("Writing to unmapped memory 0x%x -> 0x%x", physicalAddress, value))
}

func (m *Machine) InitPeripherals() {
	m.cpu.AddressSpace = m // connect cpu to bus

	m.memoryMap = []util.MemoryRange{
		{Start: 0x10000000, End: 0x1FBFFFFF, Name: "Cartridge ROM", P: &m.cartridgeRom},
		{Start: 0x00000000, End: 0x003FFFFF, Name: "RDRAM", P: make(Memory, 0x400000)},
		{Start: 0x03F00000, End: 0x03FFFFFF, Name: "RDRAM MMIO", P: &Unused{}},
		{Start: 0x04000000, End: 0x04000FFF, Name: "RSP Data Memory", P: make(Memory, 0x1000)},
		{Start: 0x04001000, End: 0x04001FFF, Name: "RSP Instruction Memory", P: make(Memory, 0x1000)},
		{Start: 0x04040000, End: 0x040FFFFF, Name: "SP Registers", P: mips.CreateSpRegs(&m.rsp)},
		{Start: 0x04300000, End: 0x043FFFFF, Name: "MIPS Interface", P: &m.cpu.Mi},
		{Start: 0x04600000, End: 0x046FFFFF, Name: "Peripheral Interface", P: CreatePi(m)},
		{Start: 0x04700000, End: 0x047FFFFF, Name: "RDRAM settings", P: &Unused{}},
	}
}

func (m *Machine) Reset() {
	m.cpu.Reset()
	/*
		The first 0x1000 bytes from the cartridge are copied to SP DMEM.
		This is implemented as a copy of 0x1000 bytes from 0xB0000000 to 0xA4000000.
	*/
	for i := uint64(0); i < 0x1000/4; i++ {
		m.Write(0xA4000000+4*i, m.Read(0xB0000000+4*i))
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

	m.cartridgeRom = Memory(data)
	fmt.Printf("Successfully loaded ROM: %s (%d B)\n", filePath, len(data))
	return nil
}
