package main

import (
	"fmt"
	"time"

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

	for {
		msg := time.Now().String()
		_, _ = serve.Publish(msg)
		time.Sleep(100 * time.Millisecond)
		fmt.Println(msg)
	}
}
