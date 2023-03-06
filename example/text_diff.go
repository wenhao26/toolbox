package main

import (
	"fmt"

	"github.com/wenhao26/toolbox/utils"
)

func main() {
	text1 := "错误指的是可能出现问题的地方出现了问题。比如打开一个文件时失败，这种情况在人们的意料之中 。"
	text2 := "出现问题的地方比如引用了空指针出现了问题。这种情况在人们的意料之中 。测试一下吧"
	result := utils.TextCompare(text1, text2)
	fmt.Println(result)
}
