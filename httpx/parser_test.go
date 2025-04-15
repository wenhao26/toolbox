package httpx

import (
	"testing"
)

func TestParseJSON(t *testing.T) {
	client := NewClient(Config{})
	var result map[string]interface{}

	err := client.Post("https://httpbin.org/post").
		WithJSON(map[string]interface{}{
			"name": "tester",
		}).
		ParseJSON(&result)

	if err != nil {
		t.Fatalf("parse json failed: %v", err)
	}

	if _, ok := result["json"]; !ok {
		t.Error("expected key `json` in response")
	}
}
