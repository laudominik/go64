package emu

import (
	"errors"
	"fmt"
	"os"
)

type Machine struct {
	cpu          Cpu
	cardridgeROM []byte // Cartridge Domain 1 Address 2
	rspData      []byte // SP DMEM
}

func inRange(arg, left, right uint64) bool {
	return arg >= left && arg <= right
}

func (m *Machine) Tick() {
	pc := m.cpu.pc
	instrb := m.readDWord(pc)
	instr := decode(instrb)
	m.execute(instr)
}

func (m *Machine) execute(instr Instruction) {
	m.cpu.pc += 4
	instr.callback(m, instr)
}

func (m *Machine) readDWord(virtualAddress uint64) uint32 {
	if virtualAddress&0b11 != 0 {
		panic(fmt.Sprintf("Unaligned read at %x: %x", m.cpu.pc, virtualAddress))
	}
	virtualAddress &= 0xFFFFFFFF // todo: remove after implementing 64-bit mode
	/* big endian */
	hh := uint32(m.read(virtualAddress))
	hl := uint32(m.read(virtualAddress + 1))
	lh := uint32(m.read(virtualAddress + 2))
	ll := uint32(m.read(virtualAddress + 3))
	return ll + (lh << 8) + (hl << 16) + (hh << 24)
}

func (m *Machine) writeDWord(virtualAddress uint64, value uint32) {
	if virtualAddress&0b11 != 0 {
		panic(fmt.Sprintf("Unaligned write at %x: %x", m.cpu.pc, virtualAddress))
	}
	hh := byte((value >> 24) & 0xFF)
	hl := byte((value >> 16) & 0xFF)
	lh := byte((value >> 8) & 0xFF)
	ll := byte(value & 0xFF)
	m.write(virtualAddress, hh)
	m.write(virtualAddress+1, hl)
	m.write(virtualAddress+2, lh)
	m.write(virtualAddress+3, ll)
}

func (m *Machine) write(virtualAddress uint64, value byte) {
	physicalAddress := translate(virtualAddress)
	if inRange(physicalAddress, 0x10000000, 0x1FBFFFFF) {
		panic("Trying to write to cardridge ROM")
	} else if inRange(physicalAddress, 0x04000000, 0x04000FFF) {
		m.rspData[physicalAddress-0x04000000] = value
	}
}

func (m *Machine) read(virtualAddress uint64) byte {
	physicalAddress := translate(virtualAddress)
	if inRange(physicalAddress, 0x10000000, 0x1FBFFFFF) {
		return m.cardridgeROM[physicalAddress-0x10000000]
	} else if inRange(physicalAddress, 0x04000000, 0x04000FFF) {
		return m.rspData[physicalAddress-0x04000000]
	}
	return 0
}

func (m *Machine) Reset() {
	m.cpu.reset()
	m.rspData = make([]byte, 0x1000)
	/*
		The first 0x1000 bytes from the cartridge are copied to SP DMEM.
		This is implemented as a copy of 0x1000 bytes from 0xB0000000 to 0xA4000000.
	*/
	for i := uint64(0); i < 0x1000; i++ {
		m.write(0xA4000000+i, m.read(0xB0000000+i))
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
