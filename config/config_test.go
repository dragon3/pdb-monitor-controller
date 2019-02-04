package config

import "testing"

func TestNewFromEnv(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		c, err := NewFromEnv()
		if err != nil {
			t.Fatal(err)
		}
		if c.LogLevel != "INFO" {
			t.Fatalf("expected LogLevel was INFO but got %s", c.LogLevel)
		}
		if c.Env != "development" {
			t.Fatalf("expected Env was development but got %s", c.Env)
		}
	})
}
