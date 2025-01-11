package main

import (
	"fmt"
	"go64/config"
	"go64/emu"
	"os"
)

func main() {
	var machine emu.Machine

	config.ParseConfig()

	if len(os.Args) != 2 {
		fmt.Println("Usage: go64 <path-to-rom>")
		os.Exit(1)
	}

	romPath := os.Args[1]
	err := machine.LoadRom(romPath)
	if err != nil {
		fmt.Printf("Error loading ROM: %v\n", err)
		os.Exit(1)
	}
	machine.Reset()

	for {
		machine.Tick()
	}
}
