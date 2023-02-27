package main

import (
	"github.com/wenhao26/toolbox/mq/kafka"
)

func main() {
	serve, err := kafka.NewConsumer(&kafka.Option{Addr: []string{"localhost:9092"}})
	if err != nil {
		panic(err)
	}

	serve.Receive("local-topic1")
}
