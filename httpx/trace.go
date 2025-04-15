package httpx

import (
	"crypto/rand"
	"encoding/hex"
)

// 生成16位的 Trace-ID
func GenerateTraceID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
