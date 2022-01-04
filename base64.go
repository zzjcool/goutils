/*
 * @Author: ccchieh
 * @Github: https://github.com/ccchieh
 * @Email: email@zzj.cool
 * @Date: 2021-04-20 19:09:02
 * @LastEditors: ccchieh
 * @LastEditTime: 2021-04-20 19:14:19
 */
package goutils

import "encoding/base64"

// n is the input len
func Base64Encode(n int, input []byte) []byte {
	base64Len := base64.URLEncoding.EncodedLen(n)
	base64Output := make([]byte, base64Len)
	base64.URLEncoding.Encode(base64Output, input[:n])
	return base64Output
}
