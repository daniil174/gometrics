package memstats

import (
	"math/rand"
	"runtime"
)

type CounterMetric struct {
	Name  string
	Value int64
}

type GaugeMetric struct {
	Name  string
	Value float64
}

var m runtime.MemStats
var PullCount int64 = 0

func CollectGaugeMetrics() []GaugeMetric {
	runtime.ReadMemStats(&m)
	return []GaugeMetric{
		{Name: "Alloc", Value: float64(m.Alloc)},
		{Name: "BuckHashSys", Value: float64(m.BuckHashSys)},
		{Name: "Frees", Value: float64(m.Frees)},
		{Name: "GCCPUFraction", Value: m.GCCPUFraction},
		{Name: "GCSys", Value: float64(m.GCSys)},
		{Name: "HeapAlloc", Value: float64(m.HeapAlloc)},
		{Name: "HeapIdle", Value: float64(m.HeapIdle)},
		{Name: "HeapInuse", Value: float64(m.HeapInuse)},
		{Name: "HeapObjects", Value: float64(m.HeapObjects)},
		{Name: "HeapReleased", Value: float64(m.HeapReleased)},
		{Name: "HeapSys", Value: float64(m.HeapSys)},
		{Name: "LastGC", Value: float64(m.LastGC)},
		{Name: "Lookups", Value: float64(m.Lookups)},
		{Name: "MCacheInuse", Value: float64(m.MCacheInuse)},
		{Name: "MCacheSys", Value: float64(m.MCacheSys)},
		{Name: "MSpanInuse", Value: float64(m.MSpanInuse)},
		{Name: "MSpanSys", Value: float64(m.MSpanSys)},
		{Name: "Mallocs", Value: float64(m.Mallocs)},
		{Name: "NextGC", Value: float64(m.NextGC)},
		{Name: "NumForcedGC", Value: float64(m.NumForcedGC)},
		{Name: "NumGC", Value: float64(m.NumGC)},
		{Name: "OtherSys", Value: float64(m.OtherSys)},
		{Name: "PauseTotalNs", Value: float64(m.PauseTotalNs)},
		{Name: "StackInuse", Value: float64(m.StackInuse)},
		{Name: "StackSys", Value: float64(m.StackSys)},
		{Name: "Sys", Value: float64(m.Sys)},
		{Name: "TotalAlloc", Value: float64(m.TotalAlloc)},
		{Name: "RandomValue", Value: 100 * rand.Float64()},
	}
}

//func ReturnListOfCouners() []CounterMetric {
//	return []CounterMetric{
//		{Name: "PollCount", Value: PullCount},
//	}
//}
//
//func ReturnListOfGauges() []GaugeMetric {
//	return []GaugeMetric{
//		{Name: "Alloc", Value: float64(m.Alloc)},
//		{Name: "BuckHashSys", Value: float64(m.BuckHashSys)},
//		{Name: "Frees", Value: float64(m.Frees)},
//		{Name: "GCCPUFraction", Value: m.GCCPUFraction},
//		{Name: "GCSys", Value: float64(m.GCSys)},
//		{Name: "HeapAlloc", Value: float64(m.HeapAlloc)},
//		{Name: "HeapIdle", Value: float64(m.HeapIdle)},
//		{Name: "HeapInuse", Value: float64(m.HeapInuse)},
//		{Name: "HeapObjects", Value: float64(m.HeapObjects)},
//		{Name: "HeapReleased", Value: float64(m.HeapReleased)},
//		{Name: "HeapSys", Value: float64(m.HeapSys)},
//		{Name: "LastGC", Value: float64(m.LastGC)},
//		{Name: "Lookups", Value: float64(m.Lookups)},
//		{Name: "MCacheInuse", Value: float64(m.MCacheInuse)},
//		{Name: "MCacheSys", Value: float64(m.MCacheSys)},
//		{Name: "MSpanInuse", Value: float64(m.MSpanInuse)},
//		{Name: "MSpanSys", Value: float64(m.MSpanSys)},
//		{Name: "Mallocs", Value: float64(m.Mallocs)},
//		{Name: "NextGC", Value: float64(m.NextGC)},
//		{Name: "NumForcedGC", Value: float64(m.NumForcedGC)},
//		{Name: "NumGC", Value: float64(m.NumGC)},
//		{Name: "OtherSys", Value: float64(m.OtherSys)},
//		{Name: "PauseTotalNs", Value: float64(m.PauseTotalNs)},
//		{Name: "StackInuse", Value: float64(m.StackInuse)},
//		{Name: "StackSys", Value: float64(m.StackSys)},
//		{Name: "Sys", Value: float64(m.Sys)},
//		{Name: "TotalAlloc", Value: float64(m.TotalAlloc)},
//		// {Name: "PollCount", Type: "counter", Value: int64(pollCount)},
//		{Name: "RandomValue", Value: 100 * rand.Float64()},
//	}
//}
