package memory

type Memory [100]bool

func New() *Memory {
	return &Memory{}
}

func (m *Memory) Malloc(near int) int {
	for i := 0; i < len(m); i++ {
		if !m[i] {
			m[i] = true
			return i
		}
	}

	panic("Memory is full!")
}

func (m *Memory) Free(p int) {
	m[p] = false
}
