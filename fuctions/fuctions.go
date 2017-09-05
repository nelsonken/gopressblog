package functions

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5Encode md5 encode
func MD5Encode(s string) string {
	encoder := md5.New()
	encoder.Write([]byte(s))
	b := encoder.Sum(nil)

	return hex.EncodeToString(b)
}

