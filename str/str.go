package str

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
)

func GetRandomString(n uint) string {
	randBytes := make([]byte, n/2)
	_, _ = rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

// BytesUint16 获取bytes对应的uint16数值
func BytesUint16(bs []byte) uint16 {
	ret := uint16(0)
	for _, b := range bs {
		ret <<= 8
		ret |= uint16(b)
	}
	return ret
}

// Uint16Bytes 获取uint16数值对应的byte数组
func Uint16Bytes(u uint16) []byte {
	return []byte{byte(u >> 8), byte(u)}
}

// Neat 整洁的输出
func Neat(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
