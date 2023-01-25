package dao

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateUUID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}

func GenerateRandomString(size int) string {
	str := make([]byte, size)
	_, err := rand.Read(str)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(str)
}

func GenerateZeros(size int) string {
	return string(make([]byte, size)[:])
}
