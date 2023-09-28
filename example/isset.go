package main

import (
	"fmt"
	"strings"
)

func main() {
	// 判断元素是否在切片中，类似PHP的 $map['title'] ?? ''
	str := "CN,16777728,16778239,1.0.2.0,1.0.3.255"
	arr := strings.Split(str, ",")

	info := map[int]string{}
	for k, v := range arr {
		info[k] = v
	}
	fmt.Println(info)

	if _, ok := info[8]; ok {
		fmt.Println("IN")
	} else {
		fmt.Println("NOT IN")
	}

}
