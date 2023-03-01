package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/wenhao26/toolbox/storage/mongodb"
)

type Community24h struct {
	Number    int   `json:"number"`
	CreatedAt int64 `json:"created_at"`
}

var ctx = context.Background()

func main() {
	mgo, err := mongodb.NewStorage(&mongodb.Option{
		Uri:      "mongodb://127.0.0.1:27017/?connect=direct",
		Timeout:  5 * time.Second,
		PoolSize: 10,
	})
	if err != nil {
		panic(err)
	}

	db := mgo.Db("test")
	collection := db.Collection("community_24h")

	// 插入文档
	//document := Community24h{Number: 888, CreatedAt: time.Now().Unix()}
	//insertResult, err := collection.InsertOne(ctx, document)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(insertResult.InsertedID)

	// 批次插入文档
	//documents := []interface{}{
	//	Community24h{Number: 888, CreatedAt: time.Now().Unix()},
	//	Community24h{Number: 777, CreatedAt: time.Now().Unix()},
	//	Community24h{Number: 666, CreatedAt: time.Now().Unix()},
	//}
	//insertResult, err := collection.InsertMany(ctx, documents)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(insertResult.InsertedIDs)

	// 更新文档
	filter := bson.D{{"number", 888}} // 过滤条件
	// 更新条件
	update := bson.D{
		{
			"$set",
			bson.D{{"created_at", time.Now().Unix()}},
		},
	}
	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		panic(err)
	}
	fmt.Println(updateResult.ModifiedCount)

	// todo more...

}
