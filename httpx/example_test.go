package httpx

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestHttpxProduction(t *testing.T) {
	client := NewClient(Config{
		Timeout:     5 * time.Second,
		Retries:     2,
		EnableTrace: true,
		Middleware: []MiddlewareFunc{
			func(req *http.Request) {
				req.Header.Set("X-App-Version", "1.0.0")
			},
		},
	})

	var result map[string]interface{}
	err := client.Post("https://httpbin.org/post").
		WithJSON(map[string]string{"username": "admin"}).
		ParseJSON(&result)

	if err != nil {
		t.Fatal("请求失败:", err)
	}

	fmt.Println("✅ 请求结果：", result)
}
