package centrifugo

import (
	"context"

	"github.com/centrifugal/gocent/v3"
)

// Connection configuration
type Option struct {
	Addr string
	Key  string
}

// Serve structure
type Serve struct {
	Cent *gocent.Client
}

var ctx = context.Background()

// Create Centrifugal Serve instance
func NewServe(option *Option) *Serve {
	client := gocent.New(gocent.Config{
		Addr: option.Addr,
		Key:  option.Key,
	})
	return &Serve{Cent: client}
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
