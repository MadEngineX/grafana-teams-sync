package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/google/uuid"
)

func GenerateLoggerID() string {

	id := uuid.New()
	hash := sha1.New()
	hash.Write([]byte(id.String()))
	hashBytes := hash.Sum(nil)
	shortSHA := hex.EncodeToString(hashBytes)[:6]
	return shortSHA
}
