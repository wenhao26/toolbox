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
		Addr:   "http://127.0.0.1:9501/api",
		ApiKey: "99e2a129-e85d-4cfc-bdd4-250df0525165",
	})

	// 生成连接JWT
	//tokenSecretKey := "5c1b094c-dc02-43e4-8cd8-0a3d08531112"
	//token1 := cent.SetSecret(tokenSecretKey).GenConnToken("1688", 0, nil, nil)
	//fmt.Println(token1)
	//os.Exit(0)
	////token2 := cent.SetSecret(tokenSecretKey).GenPrivateChannelToken("1688", "test-channel", 0, nil)
	////fmt.Println(token2)
	//os.Exit(0)

	// 发布消息
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

	// Mock - 使用ants协程池，模拟数据并发推送
	//wg := sync.WaitGroup{}
	//publisher := func() {
	//	for i := 0; i < 10000; i++ {
	//		body := map[string]string{
	//			"time": time.Now().String(),
	//		}
	//		data, _ := json.Marshal(body)
	//		publishResult, err := cent.Publish(channel, data)
	//		if err != nil {
	//			log.Println("发布消息异常:", err)
	//		}
	//		fmt.Printf("发布到频道 %s 成功, 流位置 {offset: %d, epoch: %s} \n", channel, publishResult.Offset, publishResult.Epoch)
	//		//time.Sleep(100 * time.Millisecond)
	//	}
	//}
	//
	//p, _ := ants.NewPool(5)
	//defer p.Release()
	//for i := 0; i < 5; i++ {
	//	wg.Add(1)
	//	_ = p.Submit(func() {
	//		publisher()
	//		wg.Done()
	//	})
	//}
	//wg.Wait()
	//
	//fmt.Println("Running Number:", p.Running())
}
