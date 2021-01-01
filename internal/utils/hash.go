package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 hash 散列
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
