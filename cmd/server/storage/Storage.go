package storage

import "errors"

var ErrMetricDidntExist = errors.New("metric didn't exist")

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge: map[string]float64{
			"Alloc":         0,
			"BuckHashSys":   0,
			"Frees":         0,
			"GCCPUFraction": 0,
			"GCSys":         0,
			"HeapAlloc":     0,
			"HeapIdle":      0,
			"HeapInuse":     0,
			"HeapObjects":   0,
			"HeapReleased":  0,
			"HeapSys":       0,
			"LastGC":        0,
			"Lookups":       0,
			"MCacheInuse":   0,
			"MCacheSys":     0,
			"MSpanInuse":    0,
			"MSpanSys":      0,
			"Mallocs":       0,
			"NextGC":        0,
			"NumForcedGC":   0,
			"NumGC":         0,
			"OtherSys":      0,
			"PauseTotalNs":  0,
			"StackInuse":    0,
			"StackSys":      0,
			"Sys":           0,
			"TotalAlloc":    0,
			"RandomValue":   0,
		},
		Counter: map[string]int64{
			"PollCount": 0,
		},
	}
}

func (m *MemStorage) AddCounter(name string, value int64) error {
	if _, ok := m.Counter[name]; !ok {
		return ErrMetricDidntExist
	}
	m.Counter[name] += value
	return nil
}

func (m *MemStorage) RewriteGauge(name string, value float64) error {
	if _, ok := m.Gauge[name]; !ok {
		return ErrMetricDidntExist
	}
	m.Gauge[name] = value
	return nil
}

func (m *MemStorage) GetGauge(name string) (float64, error) {
	if _, ok := m.Gauge[name]; !ok {
		return 0, ErrMetricDidntExist
	}
	return m.Gauge[name], nil
}

func (m *MemStorage) GetCounter(name string) (int64, error) {
	if _, ok := m.Counter[name]; !ok {
		return 0, ErrMetricDidntExist
	}
	return m.Counter[name], nil
}
