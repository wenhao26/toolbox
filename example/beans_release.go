package main

import (
	"fmt"

	"github.com/wenhao26/toolbox/mq/beans"
)

func main() {
	serve, err := beans.NewServe(&beans.Option{Addr: "127.0.0.1:11300"})
	if err != nil {
		panic(err)
	}

	tubeInfos, err := serve.WatchTubes()
	if err != nil {
		panic(err)
	}
	fmt.Println(tubeInfos)

	/*jodID, err := serve.Release("T1", []byte(time.Now().String()), 1, 3e9, 5e9)
	if err != nil {
		panic(err)
	}
	fmt.Println("jod-id:", jodID)*/
}
