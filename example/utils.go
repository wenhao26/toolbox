package main

import (
	"fmt"

	"github.com/wenhao26/toolbox/utils"
)

func main() {
	/*password := "123456"
	hash := utils.GenPassword(password)
	fmt.Println(hash)
	fmt.Println(utils.VerifyPassword(hash, password))*/

	/*params := map[string]interface{}{
		"app-version": "1.0.0",
		"app-key":     "12345678",
		"timestamp":   1677226735,
		"nonce":       "1234",
		"client":      "PC",
		"uuid":        "11-22-33-44",
	}
	secret := "yourkey"
	sign := utils.GenSimpleSignature(params, secret, "")
	fmt.Println(sign)*/

	/*keys, err := utils.GenRSAKey()
	if err != nil {
		panic(err)
	}
	for _, key := range keys {
		fmt.Println(string(key))
	}
	fmt.Println(keys)*/

	dates := utils.GetDatesWithinTimeRange("2023-12-01", "2024-02-04")
	for _, date := range dates {
		fmt.Println(date)
	}

}
