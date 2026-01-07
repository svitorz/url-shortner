package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DbName     string
	DbHost     string
	DbPort     int64
	DbUser     string
	DbPassword string

	AppPort int64

	RedisHost     string
	RedisPort     int64
	RedisDb       int64
	RedisProtocol int64
}

func LoadConfig() (*Config, error) {
	var cfg *Config

	err := godotenv.Load()

	if err != nil {
		return cfg, err
	}

	dbport, err := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 64)

	if err != nil {
		return cfg, err
	}

	appport, err := strconv.ParseInt(os.Getenv("APP_PORT"), 10, 64)

	if err != nil {
		return cfg, err
	}

	redisport, err := strconv.ParseInt(os.Getenv("REDIS_PORT"), 10, 64)

	if err != nil {
		return cfg, err
	}

	redisdb, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 64)

	if err != nil {
		return cfg, err
	}

	redprotocol, err := strconv.ParseInt(os.Getenv("REDIS_PROTOCOL"), 10, 64)

	if err != nil {
		return cfg, err
	}
	return &Config{
		DbName:        os.Getenv("DB_NAME"),
		DbHost:        os.Getenv("DB_HOST"),
		DbPort:        dbport,
		DbUser:        os.Getenv("DB_USER"),
		DbPassword:    os.Getenv("DB_PASSWORD"),
		AppPort:       appport,
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     redisport,
		RedisDb:       redisdb,
		RedisProtocol: redprotocol,
	}, nil
}
