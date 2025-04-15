package httpx

import (
	"testing"
)

// go test -v -run TestRequestBuilding
func TestRequestBuilding(t *testing.T) {
	client := NewClient(Config{})
	req := client.Get("https://httpbin.org/get").
		WithQuery(map[string]string{"id": "123"}).
		WithHeader("X-Test", "true")

	if req == nil {
		t.Fatal("expected non-nil request")
	}

	if req.Method() != "GET" {
		t.Errorf("expected GET, got %s", req.Method())
	}
}
