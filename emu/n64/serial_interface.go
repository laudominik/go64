package n64

import "fmt"

type Si struct{}

func (f *Si) Read(reg uint64) uint32 {
	switch reg {
	default:
		panic(fmt.Sprintf("Reading from unimplemented SI register 0x%x", reg))
	}
}

func (f *Si) Write(reg uint64, value uint32) {
	switch reg {
	default:
		panic(fmt.Sprintf("Writing to unimplemented SI register 0x%x", reg))
	}
}
