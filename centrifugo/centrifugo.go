package centrifugo

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/centrifugal/gocent/v3"
)

// Connection configuration
type Option struct {
	Addr   string
	ApiKey string
}

// Serve structure
type Serve struct {
	Cent   *gocent.Client
	Secret string
}

var ctx = context.Background()

// Create Centrifugal Serve instance
func NewServe(option *Option) *Serve {
	client := gocent.New(gocent.Config{
		Addr: option.Addr,
		Key:  option.ApiKey,
	})
	return &Serve{Cent: client}
}

// Set token key
func (s *Serve) SetSecret(secret string) *Serve {
	s.Secret = secret
	return s
}

// Generate connection JWT
func (s *Serve) GenConnToken(uid string, exp int, info, channels interface{}) string {
	header := map[string]string{
		"typ": "JWT",
		"alg": "HS256",
	}
	payload := map[string]interface{}{
		"sub": uid,
	}

	if exp > 0 {
		payload["exp"] = exp
	}
	if info != nil {
		payload["info"] = info
	}
	if channels != nil {
		payload["channels"] = channels
	}

	var segments []string
	headerJson, _ := json.Marshal(header)
	headerEncodeStr := urlSafeB64Encode(headerJson)
	segments = append(segments, headerEncodeStr)

	payloadJson, _ := json.Marshal(payload)
	payloadEncodeStr := urlSafeB64Encode(payloadJson)
	segments = append(segments, payloadEncodeStr)

	signingInput := strings.Join(segments, ".")
	signature := sign(signingInput, []byte(s.Secret))
	signatureStr := urlSafeB64Encode(signature)
	segments = append(segments, signatureStr)

	return strings.Join(segments, ".")
}

// Generate private channel JWT
func (s *Serve) GenPrivateChannelToken(client, channel string, exp int, info interface{}) string {
	header := map[string]string{
		"typ": "JWT",
		"alg": "HS256",
	}
	payload := map[string]interface{}{
		"client":  client,
		"channel": channel,
	}

	if exp > 0 {
		payload["exp"] = exp
	}
	if info != nil {
		payload["info"] = info
	}

	var segments []string
	headerJson, _ := json.Marshal(header)
	headerEncodeStr := urlSafeB64Encode(headerJson)
	segments = append(segments, headerEncodeStr)

	payloadJson, _ := json.Marshal(payload)
	payloadEncodeStr := urlSafeB64Encode(payloadJson)
	segments = append(segments, payloadEncodeStr)

	signingInput := strings.Join(segments, ".")
	signature := sign(signingInput, []byte(s.Secret))
	signatureStr := urlSafeB64Encode(signature)
	segments = append(segments, signatureStr)

	return strings.Join(segments, ".")
}

// Publish message to channel
func (s *Serve) Publish(channel string, body []byte) (gocent.PublishResult, error) {
	result, err := s.Cent.Publish(ctx, channel, body)
	if err != nil {
		return gocent.PublishResult{}, err
	}
	return result, nil
}

// Publish message to pipe
func (s *Serve) PipePublish(channel string, data [][]byte) ([]gocent.Reply, error) {
	pipe := s.Cent.Pipe()

	for _, v := range data {
		_ = pipe.AddPublish(channel, v)
	}
	replies, err := s.Cent.SendPipe(ctx, pipe)
	if err != nil {
		return nil, err
	}
	return replies, nil
}

func urlSafeB64Encode(input []byte) string {
	str := base64.URLEncoding.EncodeToString(input)
	str = strings.Replace(str, "+/", "-_", -1)
	str = strings.Replace(str, "=", "", -1)
	return str
}

func sign(cipherText string, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(cipherText))
	return h.Sum(nil)
}
