package md5_helper

import (
	"crypto/md5"
	"encoding/hex"
)

// md5加密
func Md5EncodingOnly(content string) (string, error) {
	hash := md5.New()
	_, err := hash.Write([]byte(content))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// md5加密加盐
func Md5EncodingWithSalt(content string, salt string) (string, error) {
	if salt == "" {
		return content, nil
	}
	return Md5EncodingOnly(content + "{" + salt + "}")
}
