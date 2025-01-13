package peripherals

type Mi struct {
}

func (mi *Mi) Read(reg uint64) uint32 {
	return 0
}

func (mi *Mi) Write(reg uint64, value uint32) {

}
