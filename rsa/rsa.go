package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

// Cipher structure
type Cipher struct {
	PubKey []byte
	PriKey []byte
}

// Create RSA instance
func NewRSACipher(pubKeyFile, priKeyFile string) (*Cipher, error) {
	pubKey, err := readKey(pubKeyFile)
	if err != nil {
		return nil, errors.New("The public key file cannot be read")
	}
	priKey, err := readKey(priKeyFile)
	if err != nil {
		return nil, errors.New("The private key file cannot be read")
	}

	return &Cipher{PubKey: pubKey, PriKey: priKey}, nil
}

// Public Key encrypt
func (c *Cipher) PublicKeyEncrypt(data []byte) []byte {
	block, _ := pem.Decode(c.PubKey)
	pubInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)
	pubKey := pubInterface.(*rsa.PublicKey)
	encrypted, _ := rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
	return encrypted
}

// Private Key decrypt
func (c *Cipher) PrivateKeyDecrypt(text string) []byte {
	block, _ := pem.Decode(c.PriKey)
	decodeText, _ := base64.StdEncoding.DecodeString(text)
	priKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	decrypted, _ := rsa.DecryptPKCS1v15(rand.Reader, priKey, decodeText)
	return decrypted
}

// Read public and private keys
func readKey(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}
