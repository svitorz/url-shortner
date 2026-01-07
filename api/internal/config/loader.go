package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var AppPort int64

func init() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error to handle .env", err)
	}
}

func LoadAppConfig() (int, error) {
	appport, err := strconv.ParseInt(os.Getenv("APP_PORT"), 10, 64)

	if err != nil {
		return 0, err
	}

	return int(appport), nil
}
