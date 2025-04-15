package httpx

import (
	"net/http"
	"time"
)

// Config 客户端配置结构体
type Config struct {
	Timeout     time.Duration                 // 请求超时时间
	Retries     int                           // 最大重试次数（建议不超过3次）
	Backoff     func(retry int) time.Duration // 自定义的重试退避策略
	Middleware  []MiddlewareFunc              // 全局中间件
	EnableTrace bool                          // 是否启用 X-trace-ID
}

// Client HTTP 客户端结构体
type Client struct {
	config Config
	client *http.Client
}

// NewClient 创建一个新的 HTTP 客户端
func NewClient(cfg Config) *Client {
	// 默认退避策略：线性递增
	if cfg.Backoff == nil {
		cfg.Backoff = DefaultBackoff()
	}

	return &Client{
		config: cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}
