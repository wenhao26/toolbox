package centrifugo

import (
	"context"

	"github.com/centrifugal/gocent/v3"
)

type Option struct {
	Addr string
	Key  string
}

type Serve struct {
	Cent *gocent.Client
}

var ctx = context.Background()

func NewServe(option *Option) *Serve {
	client := gocent.New(gocent.Config{
		Addr: option.Addr,
		Key:  option.Key,
	})
	return &Serve{Cent: client}
}

func (s *Serve) Publish(channel string, body []byte) (gocent.PublishResult, error) {
	result, err := s.Cent.Publish(ctx, channel, body)
	if err != nil {
		return gocent.PublishResult{}, err
	}
	return result, nil
}
