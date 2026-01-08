package config

import (
	"os"
	"strconv"
)

func LoadRateLimitConf() ([]int, error) {
	limit, err := strconv.Atoi(os.Getenv("REQUISITION_LIMIT"))

	if err != nil {
		return []int{}, err
	}

	duration, err := strconv.Atoi(os.Getenv("MINUTES_OF_DURATION"))

	if err != nil {
		return []int{}, err
	}

	return []int{limit, duration}, nil
}
