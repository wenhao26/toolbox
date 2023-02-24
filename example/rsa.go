package main

import (
	"encoding/base64"
	"fmt"

	"github.com/wenhao26/toolbox/rsa"
)

func main() {
	pubKeyFile := "../keys/public_key.pem"
	priKeyFile := "../keys/private_key.pem"
	cipher, err := rsa.NewRSACipher(pubKeyFile, priKeyFile)
	if err != nil {
		panic(err)
	}

	text := []byte("Key ID:888888899999ABC#$*`")
	encrypted := cipher.PublicKeyEncrypt(text)
	base64Str := base64.StdEncoding.EncodeToString(encrypted)
	fmt.Println("公钥加密(base64):", base64Str)

	decrypted := cipher.PrivateKeyDecrypt(base64Str)
	fmt.Println("私钥解密结果:", string(decrypted))
}
