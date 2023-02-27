package main

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/wenhao26/toolbox/mq/kafka"
)

type Message struct {
	MsgID string `json:"msg_id"`
	Body  string `json:"body"`
}

func main() {
	/*serve, err := kafka.NewSyncProducer(&kafka.Option{Addr: []string{"localhost:9092"}})
	if err != nil {
		panic(err)
	}
	for {
		message := Message{
			MsgID: strconv.FormatInt(time.Now().Unix(), 10),
			Body:  time.Now().String(),
		}
		data, _ := json.Marshal(message)
		result, err := serve.Release("local-topic1", data)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
		time.Sleep(100 * time.Millisecond)
	}*/

	serve, err := kafka.NewAsyncProducer(&kafka.Option{Addr: []string{"localhost:9092"}})
	if err != nil {
		panic(err)
	}
	for {
		message := Message{
			MsgID: strconv.FormatInt(time.Now().Unix(), 10),
			Body:  time.Now().String(),
		}
		data, _ := json.Marshal(message)
		serve.AsyncRelease("local-topic1", data)
		time.Sleep(10 * time.Millisecond)
	}

}
