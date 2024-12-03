// config.go
package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr      string
	RedisPassword  string
	LimitPerSecond int
	BlockDuration  time.Duration
	StorageType    string
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func LoadConfig() (Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting working directory")
	}

	// debug
	envPath := filepath.Join(wd, "config", ".env")
	log.Printf("Verificando se o arquivo .env existe: %v", fileExists(envPath))
	os.Setenv("STORAGE_TYPE", "mock") // debug - Forçar manualmente para "mock"

	err = godotenv.Load(filepath.Join(wd, "config", ".env"))

	//debug
	log.Printf("Loading .env from: %s", filepath.Join(wd, "config", ".env"))
	log.Printf("STORAGE_TYPE: %s", os.Getenv("STORAGE_TYPE"))
	log.Printf("REDIS_ADDR: %s", os.Getenv("REDIS_ADDR"))
	log.Printf("LIMIT_PER_SECOND: %s", os.Getenv("LIMIT_PER_SECOND"))

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Se STORAGE_TYPE estiver vazio, definir um valor padrão para teste
	if os.Getenv("STORAGE_TYPE") == "" {
		os.Setenv("STORAGE_TYPE", "mock")
	}

	log.Printf("STORAGE_TYPE após o set: %s", os.Getenv("STORAGE_TYPE"))

	blockDuration, err := time.ParseDuration(os.Getenv("BLOCK_DURATION"))
	if err != nil {
		log.Fatalf("Error parsing BLOCK_DURATION: %s", err)
	}

	config := Config{
		RedisAddr:      os.Getenv("REDIS_ADDR"),
		RedisPassword:  os.Getenv("REDIS_PASSWORD"),
		LimitPerSecond: mustAtoi(os.Getenv("LIMIT_PER_SECOND")),
		BlockDuration:  blockDuration,
		StorageType:    os.Getenv("STORAGE_TYPE"),
	}

	return config, nil
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Error converting %s to integer: %s", s, err)
	}
	return i
}
