package config

import (
	"os"
	"strconv"
)

type Config struct {
	TimeAdditionMs        int
	TimeSubtractionMs     int
	TimeMultiplicationsMs int
	TimeDivisionsMs       int
	ComputingPower        int
	OrchestratorUrl       string
}

func LoadConfig() *Config {
	return &Config{
		TimeAdditionMs:        getEnvAsInt("TIME_ADDITION_MS", 1000),
		TimeSubtractionMs:     getEnvAsInt("TIME_SUBTRACTION_MS", 1000),
		TimeMultiplicationsMs: getEnvAsInt("TIME_MULTIPLICATIONS_MS", 1000),
		TimeDivisionsMs:       getEnvAsInt("TIME_DIVISIONS_MS", 1000),
		ComputingPower:        getEnvAsInt("COMPUTING_POWER", 5),
		OrchestratorUrl:       getEnv("ORCHESTRATOR_URL", "http://localhost:8080"),
	}
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return defaultValue
	}
	return int(value)
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
