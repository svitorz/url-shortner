package config

import "time"

type RateLimitConfig struct {
	Limit   int // REQUISITION_LIMIT
	Minutes int // MINUTES_OF_DURATION
}

// Window é a duração da janela fixa do rate limiter.
func (c RateLimitConfig) Window() time.Duration {
	return time.Duration(c.Minutes) * time.Minute
}
