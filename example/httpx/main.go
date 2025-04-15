package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/wenhao26/toolbox/httpx"
)

func main() {
	client := httpx.NewClient(httpx.Config{
		Timeout: 3 * time.Second,
		Retries: 2,
		Backoff: nil,
		Middleware: []httpx.MiddlewareFunc{
			func(req *http.Request) {
				req.Header.Set("X-App-Version", "1.0.0")
			},
		},
		EnableTrace: true,
	})

	// 测试 JSON 请求
	/*var jsonResp map[string]interface{}
	err := client.Post("https://httpbin.org/post").
		WithJSON(map[string]interface{}{
			"username": "admin",
			"token":    "123456",
		}).
		ParseJSON(&jsonResp)
	if err != nil {
		fmt.Println("❌ JSON 请求失败:", err)
	} else {
		fmt.Println("✅ JSON 响应:", jsonResp)
	}*/

	// 测试 GET 请求
	/*var textResp string
	err = client.Get("https://httpbin.org/get").
		WithQuery(map[string]string{
			"id": "999",
		}).
		ParseText(&textResp)

	if err != nil {
		fmt.Println("❌ GET 请求失败:", err)
	} else {
		fmt.Println("✅ GET 响应:", textResp)
	}*/

	// 测试超时机制
	/*var textResp string
	err := client.Get("https://httpstat.us/200?sleep=5000").ParseText(&textResp)
	if err != nil {
		fmt.Println("❌ GET 请求失败:", err)
	} else {
		fmt.Println("✅ GET 响应:", textResp)
	}*/

	// 测试重试机制
	var textResp string
	err := client.Get("https://httpstat.us/500").ParseText(&textResp)
	if err != nil {
		fmt.Println("❌ GET 请求失败:", err)
	} else {
		fmt.Println("✅ GET 响应:", textResp)
	}

}
