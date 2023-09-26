// Package util 一般不会直接将上传的图片名暴露出来，因此我们对图片名进行 MD5 来达到这个效果
// 重名的结果会一样哦 图片尽量不要重名 会被覆盖!
package util

import (
	"crypto/md5"
	"encoding/hex"
)

func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}
