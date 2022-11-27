package str

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
)

func GetRandomString(n uint) string {
	randBytes := make([]byte, n/2)
	_, _ = rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

func GetRandomBytes(n uint) []byte {
	randBytes := make([]byte, n)
	_, _ = rand.Read(randBytes)
	return randBytes
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

// SortAndCompareStrs 对比两个字符串数组，并且排序
func SortAndCompareStrs(strs1, strs2 []string) bool {
	if len(strs1) != len(strs2) {
		return false
	}
	sort.Strings(strs1)
	sort.Strings(strs2)
	for i, str1 := range strs1 {
		str2 := strs2[i]
		if str1 != str2 {
			return false
		}
	}
	return true
}

// TempFile 将data输出到临时文件中，并且返回文件名
func TempFile(data string) string {
	tempDir := "/tmp"
	filename := GetRandomString(16)

	filePath := path.Join(tempDir, filename)
	if err := os.WriteFile(filePath, []byte(data), 0666); err != nil {
		log.Fatal(err)
	}

	return filePath
}
