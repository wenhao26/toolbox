package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Cipher structure
type Cipher struct {
	CB  cipher.Block
	Key []byte
	IV  []byte
}

// Create AES instance
func NewAESCipher(key, iv string) (*Cipher, error) {
	cb, err := aes.NewCipher([]byte(key))
	if err != nil {
		return &Cipher{}, nil
	}
	return &Cipher{CB: cb, Key: []byte(key), IV: []byte(iv)}, nil
}

// CBC encrypt
func (c *Cipher) CBCEncrypt(data []byte) []byte {
	blockSize := c.CB.BlockSize()
	data = pkcs5Padding(data, blockSize)
	blockMode := cipher.NewCBCEncrypter(c.CB, c.IV[:blockSize])

	encrypted := make([]byte, len(data))
	blockMode.CryptBlocks(encrypted, data)
	return encrypted
}

// CBC decrypt
func (c *Cipher) CBCDecrypt(data []byte) []byte {
	blockSize := c.CB.BlockSize()
	blockMode := cipher.NewCBCDecrypter(c.CB, c.IV[:blockSize])

	decrypted := make([]byte, len(data))
	blockMode.CryptBlocks(decrypted, data)
	decrypted = pkcs5UnPadding(decrypted)
	return decrypted
}

// ECB encrypt
func (c *Cipher) ECBEncrypt(data []byte) []byte {
	length := (len(data) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, data)
	pad := byte(len(plain) - len(data))
	for i := len(data); i < len(plain); i++ {
		plain[i] = pad
	}

	encrypted := make([]byte, len(plain))
	for bs, be := 0, c.CB.BlockSize(); bs <= len(data); bs, be = bs+c.CB.BlockSize(), be+c.CB.BlockSize() {
		c.CB.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return encrypted
}

// ECB decrypt
func (c *Cipher) ECBDecrypt(data []byte) []byte {
	decrypted := make([]byte, len(data))
	for bs, be := 0, c.CB.BlockSize(); bs < len(data); bs, be = bs+c.CB.BlockSize(), be+c.CB.BlockSize() {
		c.CB.Decrypt(decrypted[bs:be], data[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	return decrypted[:trim]
}

// CFB encrypt
func (c *Cipher) CFBEncrypt(data []byte) ([]byte, error) {
	encrypted := make([]byte, aes.BlockSize+len(data))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(c.CB, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], data)
	return encrypted, nil
}

// CFB decrypt
func (c *Cipher) CFBDecrypt(data []byte) ([]byte, error) {
	if len(data) < aes.BlockSize {
		return nil, errors.New("cipher text too short")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(c.CB, iv)
	stream.XORKeyStream(data, data)
	return data, nil
}

// Perform pkcs5 padding
func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// Cancel pkcs5 padding
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
