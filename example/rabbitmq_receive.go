package main

import (
	"github.com/wenhao26/toolbox/mq/rabbitmq"
)

func main() {
	serve, err := rabbitmq.NewServe(&rabbitmq.Option{
		"amqp://guest:guest@localhost:5672/",
		"test.queue_02",
		"",
		"",
	})
	if err != nil {
		panic(err)
	}

	serve.Receive()
}
