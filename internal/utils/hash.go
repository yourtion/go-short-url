package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(str string) string {
	// 实现hash散列
	h := md5.New()
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
