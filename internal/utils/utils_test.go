package utils

import "testing"

func TestGenerateShort(t *testing.T) {
	s := GenerateShort()
	t.Log(s)
}

func TestGenerateUid(t *testing.T) {
	s := GenerateUid()
	t.Log(s)
}

func TestMD5(t *testing.T) {
	s := MD5("Yourion")
	t.Log(s)
}
