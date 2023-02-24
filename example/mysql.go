package main

import (
	"encoding/json"
	"fmt"
	"log"

	"toolbox/storage/mysql"
)

type User struct {
	//gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	option := mysql.Option{
		Username:    "root",
		Password:    "root",
		Host:        "127.0.0.1",
		Port:        3306,
		Dbname:      "ginblog",
		MaxIdleConn: 10,
		MaxOpenConn: 100,
		MaxLifetime: 10e9,
	}
	storage, err := mysql.NewMySQLStorage(&option)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(storage)

	data, _ := json.Marshal(storage.SqlDB.Stats())
	log.Println(string(data))

	var user User
	storage.DB.Where("username=?", "wuwenhao").First(&user)
	fmt.Println(user)
}
