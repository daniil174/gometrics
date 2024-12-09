package servconfig

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

const DefaultStoreInterval = 10

type Config struct {
	ServerAddress   string `env:"ADDRESS"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDsn     string `env:"DATABASE_DSN"`
}

// SetConfig устанавливает и получает конфигурацию из переменных или флагов
func SetConfig() (*Config, error) {
	// Set the environment variable names
	var tmpCfg, flagCfg Config
	tmpCfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, err
	}

	//fmt.Printf("Config from ENV: \n ADDRESS=%s \n FileStoragePath=%s \n StoreInterval=%d \n RESTORE=%t \n",
	//	tmpCfg.ServerAddress, tmpCfg.FileStoragePath, tmpCfg.StoreInterval, tmpCfg.Restore)

	flag.StringVar(&flagCfg.ServerAddress, "a", "localhost:8080",
		"server address and port, example 127.0.0.1:8080")
	flag.StringVar(&flagCfg.FileStoragePath, "f", "store.json",
		"full file storage path, example store.json")
	flag.IntVar(&flagCfg.StoreInterval, "i", DefaultStoreInterval,
		"Time interval for saving data, example ")
	flag.BoolVar(&flagCfg.Restore, "r", true,
		"choose to restore data or not, example false ")
	flag.StringVar(&flagCfg.DatabaseDsn, "d",
		//"postgresql://localhost:5432/my_database?sslmode=disable",
		"",
		"database config string")
	flag.Parse()
	// fmt.Printf("Config after flags and default: \n ADDRESS=%s \n FileStoragePath=%s \n StoreInterval=%d \n RESTORE=%t \n",
	//	flagCfg.ServerAddress, flagCfg.FileStoragePath, flagCfg.StoreInterval, flagCfg.Restore)

	if tmpCfg.ServerAddress == "" {
		tmpCfg.ServerAddress = flagCfg.ServerAddress
	}
	if tmpCfg.FileStoragePath == "" {
		tmpCfg.FileStoragePath = flagCfg.FileStoragePath
	}
	if tmpCfg.StoreInterval == 0 {
		tmpCfg.StoreInterval = flagCfg.StoreInterval
	}
	if !tmpCfg.Restore {
		tmpCfg.Restore = flagCfg.Restore
	}
	if tmpCfg.DatabaseDsn == "" {
		tmpCfg.DatabaseDsn = flagCfg.DatabaseDsn
	}

	fmt.Printf("Result cfg: \n ADDRESS=%s \n FileStoragePath=%s \n StoreInterval=%d \n RESTORE=%t \n DatabaseDsn=%s \n",
		tmpCfg.ServerAddress, tmpCfg.FileStoragePath, tmpCfg.StoreInterval, tmpCfg.Restore, tmpCfg.DatabaseDsn)

	return &tmpCfg, nil
}
