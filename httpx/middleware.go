package httpx

import (
	"net/http"
)

// MiddlewareFunc 定义请求中间件函数类型
type MiddlewareFunc func(req *http.Request)
