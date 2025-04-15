package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Request HTTP 请求构造过程结构体
type Request struct {
	ctx        context.Context  // 请求上下文
	client     *Client          // 关联的 httpx.Client
	method     string           // 请求方法
	url        string           // 请求地址
	body       io.Reader        // 请求体内容
	headers    http.Header      // 请求头
	query      url.Values       // 查询参数
	middleware []MiddlewareFunc // 局部中间件链
}

// WithContext 设置请求上下文（支持取消、超时）
func (r *Request) WithContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

// WithHeader 添加单个请求头字段
//func (r *Request) WithHeader(key, value string) *Request {
//	r.headers.Set(key, value)
//	return r
//}
func (r *Request) WithHeader(params map[string]string) *Request {
	for k, v := range params {
		r.headers.Set(k, v)
	}
	return r
}

// WithQuery 添加查询参数
func (r *Request) WithQuery(params map[string]string) *Request {
	for k, v := range params {
		r.query.Set(k, v)
	}
	return r
}

// WithJSON 设置 JSON 请求体
func (r *Request) WithJSON(data interface{}) *Request {
	b, _ := json.Marshal(data)
	r.body = bytes.NewReader(b)
	r.headers.Set("Content-Type", "application/json")
	return r
}

// WithForm 设置表单请求体
func (r *Request) WithForm(data map[string]string) *Request {
	values := url.Values{}
	for k, v := range data {
		values.Set(k, v)
	}
	r.body = strings.NewReader(values.Encode())
	r.headers.Set("Content-Type", "application/x-www-form-urlencoded")
	return r
}

// WithBody 设置原始自定义 Body
func (r *Request) WithBody(data []byte, contentType string) *Request {
	r.body = bytes.NewReader(data)
	if contentType != "" {
		r.headers.Set("Content-Type", contentType)
	}
	return r
}

// Use 为当前请求追加中间件
func (r *Request) Use(m MiddlewareFunc) *Request {
	r.middleware = append(r.middleware, m)
	return r
}

// Method 请求方法
func (r *Request) Method() string {
	return r.method
}

// Get 构造 GET 请求
func (c *Client) Get(url string) *Request {
	return newRequest(c, http.MethodGet, url)
}

// Get 构造 POST 请求
func (c *Client) Post(url string) *Request {
	return newRequest(c, http.MethodPost, url)
}

// newRequest 内部方法：创建请求对象
func newRequest(c *Client, method, rawURL string) *Request {
	return &Request{
		ctx:        context.Background(),
		client:     c,
		method:     method,
		url:        rawURL,
		headers:    http.Header{},
		query:      url.Values{},
		middleware: c.config.Middleware,
	}
}
