package config

import "time"

// LinkCacheConfig: TTLs do cache do redirect, em segundos.
type LinkCacheConfig struct {
	TTLSeconds         int
	NegativeTTLSeconds int
}

func (c LinkCacheConfig) TTL() time.Duration { return time.Duration(c.TTLSeconds) * time.Second }
func (c LinkCacheConfig) NegativeTTL() time.Duration {
	return time.Duration(c.NegativeTTLSeconds) * time.Second
}
