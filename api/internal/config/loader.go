package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var AppPort int64
var HashCost int

type ApiConfig struct {
	ApiKey        string
	TokenLifeSpan int64
}

func init() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error to handle .env", err)
	}
}

func LoadApiConfig() (*ApiConfig, error) {
	tokLifeSpan, err := strconv.ParseInt(os.Getenv("TOKEN_LIFE_SPAN"), 10, 64)

	if err != nil {
		fmt.Println("Erro ao buscar tempo vida do token JWT")
		return &ApiConfig{}, err
	}

	return &ApiConfig{
		ApiKey:        os.Getenv("API_KEY"),
		TokenLifeSpan: tokLifeSpan,
	}, nil
}

func LoadHashCost() (int, error) {
	cost, err := strconv.Atoi(os.Getenv("HASH_COST"))

	if err != nil {
		return 0, err
	}

	return cost, nil
}

func LoadAppConfig() (int, error) {
	appport, err := strconv.ParseInt(os.Getenv("APP_PORT"), 10, 64)

	if err != nil {
		return 0, err
	}

	return int(appport), nil
}
