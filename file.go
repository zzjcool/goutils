/*
 * @Author: ccchieh
 * @Github: https://github.com/ccchieh
 * @Email: email@zzj.cool
 * @Date: 2021-01-18 16:10:17
 * @LastEditors: ccchieh
 * @LastEditTime: 2021-04-13 17:15:39
 */
package goutils

import "os"

// 判断所给路径文件/文件夹是否存在
func FileIsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
