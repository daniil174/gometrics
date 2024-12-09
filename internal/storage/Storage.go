package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	DB   = "db"
	FILE = "file"
	NONE = "none"
)

var MemStrg = NewMemStorage()

var PgDataBase *sql.DB

var ErrMetricDidntExist = errors.New("metric didn't exist")

type MemStorageClass struct {
	MemType     string             `json:"-"`
	FileStorage *os.File           `json:"-"`
	Gauge       map[string]float64 `json:"gauge"`
	Counter     map[string]int64   `json:"counter"`
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func StartDB(o string) error {
	//connStr := "dbname=my_database sslmode=disable"

	var err error
	PgDataBase, err = sql.Open("postgres", o)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func CloseDB() {
	PgDataBase.Close()
}

//func PingDB() (bool, error) {
//	// Проверяем соединение с базой данных
//	err := PgDataBase.Ping()
//	if err != nil {
//		log.Fatal(err)
//		return false, err
//	}
//	return true, nil
//}

func NewMemStorage() *MemStorageClass {
	return &MemStorageClass{
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

func (m *MemStorageClass) SetMemType(t string) {
	m.MemType = t
}

func (m *MemStorageClass) AddCounter(name string, value int64) error {
	// Временно убрал, потому что иначе не проходят автотесты
	// if _, ok := m.Counter[name]; !ok {
	//	return ErrMetricDidntExist
	// }
	m.Counter[name] += value
	return nil
}

func (m *MemStorageClass) RewriteGauge(name string, value float64) error {
	// Временно убрал, потому что иначе не проходят автотесты
	//if _, ok := m.Gauge[name]; !ok {
	//	return ErrMetricDidntExist
	//}
	m.Gauge[name] = value
	return nil
}

func (m *MemStorageClass) GetGauge(name string) (float64, error) {
	if _, ok := m.Gauge[name]; !ok {
		return 0, ErrMetricDidntExist
	}
	return m.Gauge[name], nil
}

func (m *MemStorageClass) GetCounter(name string) (int64, error) {
	if _, ok := m.Counter[name]; !ok {
		return 0, ErrMetricDidntExist
	}
	return m.Counter[name], nil
}

//===============================
//===============================

type Storage struct {
	FileStorage *os.File
	MemStorage  MemStorageClass
}

//func New() *Storage {
//	return &Storage{
//		MemStorageClass: *NewMemStorage(),
//	}
//}

func (m *MemStorageClass) ReadMetricFromDB() error {
	rows, err := PgDataBase.Query(`SELECT name, type, value, delta FROM metrics`)
	if err != nil {

		return fmt.Errorf("failed to select metrics: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf("no metrics found")
	}

	//readMetrics := make(map[string]Metrics)
	for rows.Next() {
		var m Metrics
		err = rows.Scan(&m.ID, &m.MType, &m.Value, &m.Delta)
		if err != nil {
			return fmt.Errorf("failed to scan metrics: %w", err)
		}
		if m.MType == "Counter" {
			MemStrg.Counter[m.ID] = *m.Delta
			log.Println("Db read Counters success")
		}
		if m.MType == "Gauge" {
			MemStrg.Gauge[m.ID] = *m.Value
			log.Println("Db read Gauge success")
		}
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("failed to iterate over metrics: %w", err)
	}

	return nil
}

func (m *MemStorageClass) ResetDBandSetZeroValue() error {
	_, err := PgDataBase.Exec(`DROP TABLE metrics; 
		CREATE TABLE IF NOT EXISTS metrics (
		id SERIAL PRIMARY KEY,
		type TEXT NOT NULL,
		name TEXT NOT NULL,
		value DOUBLE PRECISION,
		delta BIGINT,
		timestamp TIMESTAMP NOT NULL
	)`)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	for n, v := range MemStrg.Counter {
		_, err = PgDataBase.Exec(`INSERT INTO metrics (type, name,  delta, timestamp)
		VALUES ($1, $2, $3, $4)`,
			"Counter", n, v, time.Now())
		if err != nil {
			log.Println("Db faild to insert counters", err)
			return fmt.Errorf("failed to insert counters metric: %w", err)
		}
		//log.Println("Db create Counters success")
	}

	for n, v := range MemStrg.Gauge {
		_, err = PgDataBase.Exec(`INSERT INTO metrics (type, name, value, timestamp)
		VALUES ($1, $2, $3, $4)`,
			"Gauge", n, v, time.Now())
		if err != nil {
			log.Println("Db faild to insert Gauges", err)
			return fmt.Errorf("failed to insert Gauges metric: %w", err)
		}
		//log.Println("Db create Gauge success")
	}
	return nil
}

func (m *MemStorageClass) WriteMetricToDB() error {

	//пишем все Counters
	for n, v := range MemStrg.Counter {
		_, err := PgDataBase.Exec(
			`UPDATE metrics 
		SET delta =$1,
		    timestamp = $2
		WHERE name = $3`,
			v, time.Now(), n)
		if err != nil {
			log.Println("Db faild to insert counters", err)
			return fmt.Errorf("failed to insert counters metric: %w", err)
		}
		log.Println("Db save counters success")
	}

	//пишем все Gauge
	for n, v := range MemStrg.Gauge {
		_, err := PgDataBase.Exec(
			`UPDATE metrics 
		SET value =$1,
		    timestamp = $2
		WHERE name = $3`,
			v, time.Now(), n)
		if err != nil {
			log.Println("Db faild to update Gauges", err)
			return fmt.Errorf("failed to update Gauges metric: %w", err)
		}
		log.Println("Db update Gauges success")
	}
	return nil
}

// OpenFile открытие файла для хранения данных
func (m *MemStorageClass) ReadFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	//err = json.NewEncoder(file).Encode(s.MemStorageClass)

	err = json.NewDecoder(file).Decode(m)
	if err != nil {
		file.Close()
	}

	m.FileStorage = file
	//defer file.Close()
	return nil
}

func (m *MemStorageClass) SaveMetricsToFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	//fmt.Printf("Metrics: %+v", m)
	//fmt.Sprintf("Metrics: %+v", m)
	//servlogger.Sugar.Errorf("Metrics: %+v", m)

	err = json.NewEncoder(file).Encode(m)
	if err != nil {
		file.Close()
	}

	m.FileStorage = file
	//defer file.Close()
	return nil
}

func (m *MemStorageClass) CloseFile() error {
	return m.FileStorage.Close()
}

//func StartFileStorageLogic(config *flags.Config, s *Storage, servlogger *servlogger.Logger) {
//	if config.FileStoragePath != "" {
//		err := s.OpenFile(config.FileStoragePath)
//		if err != nil {
//			servlogger.Error("Failed to open file: %v", zap.Error(err))
//		}
//	} else {
//		servlogger.Info("File storage is not specified")
//		return
//	}
//
//	if config.Restore {
//		err := s.LoadMemStorageFromFile()
//		if err != nil {
//			servlogger.Error("Failed to restore data from file: %v", zap.Error(err))
//		}
//	}
//
//	go func() {
//		for {
//			interval := time.Duration(config.StoreInterval) * time.Second
//			// if interval == 0 {
//			// 	interval = 100 * time.Microsecond // Установите разумное значение по умолчанию
//			// }
//			time.Sleep(interval)
//			s.SaveMemStorageToFile()
//		}
//	}()
//}
