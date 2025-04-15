package httpx

import (
	"encoding/json"
)

// ParseJSON 将响应体解析为结构体
func (r *Request) ParseJSON(out interface{}) error {
	body, err := r.do()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, out)
}

// ParseText 将响应体解析为字符串
func (r *Request) ParseText(out *string) error {
	body, err := r.do()
	if err != nil {
		return err
	}
	*out = string(body)
	return nil
}

// ParseBytes 返回原始响应体
func (r *Request) ParseBytes() ([]byte, error) {
	return r.do()
}
