package httpx

import (
	"testing"
	"time"
)

// go test -v -run TestNewClient
func TestNewClient(t *testing.T) {
	client := NewClient(Config{
		Timeout: 5 * time.Second,
	})

	if client == nil {
		t.Fatal("expected non-nil client")
	}
}
