package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/wenhao26/toolbox/centrifugo"
)

func main() {
	cent := centrifugo.NewServe(&centrifugo.Option{
		Addr: "http://127.0.0.1:9501/api",
		Key:  "99e2a129-e85d-4cfc-bdd4-250df0525165",
	})

	channel := "test-channel"
	for {
		body := map[string]string{
			"time": time.Now().String(),
		}
		data, _ := json.Marshal(body)
		publishResult, err := cent.Publish(channel, data)
		if err != nil {
			log.Println("发布消息异常:", err)
		}
		fmt.Printf("发布到频道 %s 成功, 流位置 {offset: %d, epoch: %s} \n", channel, publishResult.Offset, publishResult.Epoch)
		time.Sleep(100 * time.Millisecond)
	}
}
