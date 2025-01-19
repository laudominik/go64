package emu

type Unused []byte

func (f *Unused) Read(reg uint64) uint32 {
	return 0
}

func (f *Unused) Write(reg uint64, value uint32) {

}
