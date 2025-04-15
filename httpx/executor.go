package httpx

import (
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

// do 执行 HTTP 请求（包含重试逻辑、中间件、日志）
func (r *Request) do() ([]byte, error) {
	urlWithQuery := r.url
	if len(r.query) > 0 {
		urlWithQuery += "?" + r.query.Encode()
	}

	var lastErr error

	for i := 0; i < r.client.config.Retries; i++ {
		req, _ := http.NewRequestWithContext(r.ctx, r.method, urlWithQuery, r.body)

		// 注入 Trace-ID
		if r.client.config.EnableTrace {
			req.Header.Set("X-Trace-ID", GenerateTraceID())
		}

		// 添加 Header
		for k, values := range r.headers {
			for _, v := range values {
				req.Header.Add(k, v)
			}
		}

		// 执行中间件（顺序执行）
		for _, m := range r.middleware {
			m(req)
		}

		start := time.Now()
		resp, err := r.client.client.Do(req)
		cost := time.Since(start)

		if err != nil {
			lastErr = err
			log.Printf("[HTTPX] %s %s fail: %v cost=%s", r.method, r.url, err, cost)
		} else {
			body, readErr := func() ([]byte, error) {
				defer resp.Body.Close() // 现在 defer 在函数体中，不会延迟到函数外
				return io.ReadAll(resp.Body)
			}()

			if readErr != nil {
				lastErr = readErr
			} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return body, nil
			} else {
				lastErr = errors.New(string(body))
			}

			log.Printf("[HTTPX] %s %s status=%d cost=%s", r.method, r.url, resp.StatusCode, cost)
		}

		time.Sleep(r.client.config.Backoff(i))
	}

	return nil, lastErr
}
