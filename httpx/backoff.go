package httpx

import (
	"time"
)

// DefaultBackoff 返回默认指数退避策略
func DefaultBackoff() func(int) time.Duration {
	return func(retry int) time.Duration {
		return time.Duration(retry+1) * 300 * time.Millisecond
	}
}
