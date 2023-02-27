package main

import (
	"github.com/wenhao26/toolbox/mq/beans"
)

func main() {
	serve, err := beans.NewServe(&beans.Option{Addr: "127.0.0.1:11300"})
	if err != nil {
		panic(err)
	}
	serve.Receive("T1")
}
