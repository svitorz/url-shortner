package config

import (
	"testing"
	"time"
)

// Antes: time.Duration(minutos * time.Now().Minute()) — nanossegundos e
// dependente do minuto do relógio no boot (17ns às hh:17, 0 às hh:00).
func TestRateLimitWindow(t *testing.T) {
	for _, tc := range []struct {
		minutes int
		want    time.Duration
	}{{1, time.Minute}, {5, 5 * time.Minute}, {60, time.Hour}} {
		if got := (RateLimitConfig{Minutes: tc.minutes}).Window(); got != tc.want {
			t.Errorf("Minutes=%d: Window() = %v, quer %v", tc.minutes, got, tc.want)
		}
	}
}

func TestLoadRejectsNonPositiveRateLimit(t *testing.T) {
	setValidEnv(t)
	t.Setenv("REQUISITION_LIMIT", "0")
	t.Setenv("MINUTES_OF_DURATION", "-1")
	if _, err := Load(); err == nil {
		t.Fatal("Load aceitou REQUISITION_LIMIT=0 e MINUTES_OF_DURATION=-1")
	}
}
