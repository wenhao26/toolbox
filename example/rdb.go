package main

import (
	"context"
	"fmt"

	"github.com/wenhao26/toolbox/storage/rdb"
)

func main() {
	storage, err := rdb.NewStorage(&rdb.Option{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	if err != nil {
		panic(err)
	}

	err = storage.Client.Set(context.TODO(), "test:string-example", "8888888", 0).Err()
	if err != nil {
		panic(err)
	}
	res, err := storage.Client.Get(context.TODO(), "test:string-example").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
