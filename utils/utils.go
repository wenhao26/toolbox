package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
	"github.com/rs/xid"
	"github.com/sergi/go-diff/diffmatchpatch"
	"golang.org/x/crypto/bcrypt"
)

// Compare Text
func CompareText(text1, text2 string) map[string]interface{} {
	result := map[string]interface{}{}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text1, text2, false)
	if len(diffs) == 0 {
		return result
	}

	result["diff"] = diffs

	// Calculate similarity
	equalCount := 0
	for _, diff := range diffs {
		// Type description:
		//
		// -1 DiffDelete item represents a delete diff.
		//  1 DiffInsert item represents an insert diff.
		//  0 DiffEqual item represents an equal diff.
		if diff.Type == 0 {
			equalCount++
		}
	}
	similarity, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(equalCount)/float64(len(diffs))), 64)
	result["similarity"] = similarity

	return result
}

// Generate RSA Public key & private key
func GenRSAKey() (map[string][]byte, error) {
	// private-key
	privateKey, err := rsa.GenerateKey(cryptoRand.Reader, 1024)
	if err != nil {
		return nil, errors.New("private key error")
	}

	derPrvStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derPrvStream,
	}
	prvKey := pem.EncodeToMemory(block)

	// public-key
	publicKey := &privateKey.PublicKey
	derPubStream, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, errors.New("public key error")
	}

	block = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPubStream,
	}
	pubKey := pem.EncodeToMemory(block)

	return map[string][]byte{
		"private_key": prvKey,
		"public_key":  pubKey,
	}, nil
}

// Generate AES key
func GenAESKey() string {
	alphabet := "BCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	key := ""
	rand.Seed(time.Now().Unix())
	for n := 0; n < 32; n++ {
		i := int(math.Floor(rand.Float64() * float64(len(alphabet))))
		key = key + alphabet[i:i+1]
	}
	return key
}

// Generate signature
func GenSimpleSignature(params map[string]interface{}, secret, signField string) string {
	// The signature field to be verified needs to be deleted
	// No signature processing
	if signField != "" {
		if _, ok := params[signField]; ok {
			delete(params, signField)
		}
	}

	// Map is sorted from small to large by ASCII code
	// Store keys in slices in sort order
	var fields []string
	for field := range params {
		fields = append(fields, field)
	}
	sort.Strings(fields)

	// Splice string
	var buf bytes.Buffer
	for i, field := range fields {
		if field != "" && params[field] != "" {
			val := ""
			switch params[field].(type) {
			case int:
				val = strconv.Itoa(params[field].(int))
			case string:
				val = params[field].(string)
			}

			if i != (len(fields) - 1) {
				val += "&"
			}
			buf.WriteString(field + "=" + val)
		}
	}

	// Sha256 encryption
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(buf.String()))
	result := hex.EncodeToString(hash.Sum(nil))

	return strings.ToUpper(result)
}

// Generate password
func GenPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

// Verify password
func VerifyPassword(encryptPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptPassword), []byte(password))
	return !(err != nil)
}

// MD5
func MD5(v string) string {
	h := md5.New()
	h.Write([]byte(v))
	return hex.EncodeToString(h.Sum(nil))
}

// Generate UUID
func GenGoogleUUID() string {
	return uuid.New().String()
}

// Generate short UUID
func GenShortUUID() string {
	return shortuuid.New()
}

// Generate short UUID with namespace
func GenShortUUIDWithNamespace(name string) string {
	return shortuuid.NewWithNamespace(name)
}

// Generate short UUID with alphabet
func GenShortUUIDWithAlphabet(alphabet string) string {
	if alphabet == "" {
		alphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxy="
	}
	return shortuuid.NewWithAlphabet(alphabet)
}

// Generate Security unique ID
func GenXID() string {
	return xid.New().String()
}

// Generate snowflake ID
func GenSnowflakeID(n int64) (int64, error) {
	node, err := snowflake.NewNode(n)
	if err != nil {
		return 0, err
	}
	return node.Generate().Int64(), nil
}
