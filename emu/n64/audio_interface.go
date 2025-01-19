package n64

type Ai struct {
	m *Machine
}

func (si *Ai) Read(reg uint64) uint32 {
	// TODO: leaving audio for now
	return 0
}

func (si *Ai) Write(reg uint64, value uint32) {
	// TODO: leaving audio for now
}
