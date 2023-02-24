package utils

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

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
