package main

import (
	"encoding/base64"
	"encoding/hex"
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

	fmt.Println("\n-------------------------------------------\n")
	data := "神奇了~"
	signature, _ := rsa.SignWithSha256([]byte(data), cipher.PriKey)
	fmt.Println("私钥签名(hex):", hex.EncodeToString(signature))
	fmt.Println("私钥签名(base64):", base64.StdEncoding.EncodeToString(signature))

	result, _ := rsa.VerySignWithSha256([]byte(data), signature, cipher.PubKey)
	fmt.Println("公钥验签:", result)

}
