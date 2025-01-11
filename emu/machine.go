package emu

import (
	"errors"
	"fmt"
	"os"
)

type Machine struct {
	cpu          Cpu
	cardridgeROM []byte // Cartridge Domain 1 Address 2
}

func (m *Machine) write(virtualAddress uint64) {
	physicalAddress := translate(virtualAddress)
	_ = physicalAddress
	// TODO: call to actual hardware
}

func (m *Machine) read(virtualAddress uint64) uint64 {
	physicalAddress := translate(virtualAddress)
	_ = physicalAddress

	// TODO: call to actual hardware
	return 0
}

func (m *Machine) Reset() {
	m.cpu.reset()

	/*
		The first 0x1000 bytes from the cartridge are then copied to SP DMEM.
		This is implemented as a copy of 0x1000 bytes from 0xB0000000 to 0xA4000000.
	*/

	// todo: copy
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
