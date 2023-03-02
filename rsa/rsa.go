package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
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

// Private Key signature
func SignWithSha256(data []byte, priKeyBytes []byte) (signature []byte, err error) {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	block, _ := pem.Decode(priKeyBytes)
	if block == nil {
		return nil, errors.New("private key error")
	}

	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	signature, err = rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA256, hashed)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// Public Key Verify signature
func VerySignWithSha256(data, signature, pubKeyBytes []byte) (bool, error) {
	block, _ := pem.Decode(pubKeyBytes)
	if block == nil {
		return false, errors.New("public key error")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signature)
	if err != nil {
		return false, err
	}
	return true, nil
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
