package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/wenhao26/toolbox/aes"
)

func main() {
	key := "eTd4ZY3UISRk5tziVikYmwC+o/D4uL40"
	iv := "BEv2g8EfAuZAeRVy"
	cipher, err := aes.NewAESCipher(key, iv)
	if err != nil {
		panic(err)
	}

	text := []byte("The differential cryptanalysis is introduced. " +
		"The properties of construction and difference on the S boxes of the date encryption standard ( DES) are discussed. " +
		"The local property of S boxes is extended to the entire cipher structure through F function. " +
		"\n\n介绍了差分密码分析，讨论了数据加密标准（DES）的S－盒的结构与差分特性，然后通过F－函数将S－盒的局部特性扩展到整个密码结构")

	log.Println("----CBC模式----")
	encrypted := cipher.CBCEncrypt(text)
	fmt.Println("密文(hex):", hex.EncodeToString(encrypted))
	fmt.Println("密文(base64):", base64.StdEncoding.EncodeToString(encrypted))
	decrypted := cipher.CBCDecrypt(encrypted)
	fmt.Println("解密结果:", string(decrypted))
	fmt.Println("\n")

	log.Println("----ECB模式----")
	encrypted2 := cipher.ECBEncrypt(text)
	fmt.Println("密文(hex):", hex.EncodeToString(encrypted2))
	fmt.Println("密文(base64):", base64.StdEncoding.EncodeToString(encrypted2))
	decrypted2 := cipher.ECBDecrypt(encrypted2)
	fmt.Println("解密结果:", string(decrypted2))
	fmt.Println("\n")

	log.Println("----CFB模式----")
	encrypted3, _ := cipher.CFBEncrypt(text)
	fmt.Println("密文(hex):", hex.EncodeToString(encrypted3))
	fmt.Println("密文(base64):", base64.StdEncoding.EncodeToString(encrypted3))
	decrypted3, _ := cipher.CFBDecrypt(encrypted3)
	fmt.Println("解密结果:", string(decrypted3))
	fmt.Println("\n")
}
