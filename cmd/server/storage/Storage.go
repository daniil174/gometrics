package storage

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

func (m *MemStorage) AddCounter(name string, value int64) error {
	m.Counter[name] += value
	return nil
}

func (m *MemStorage) RewriteGauge(name string, value float64) error {
	m.Gauge[name] = value
	return nil
}
