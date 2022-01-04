/*
 * @Author: ccchieh
 * @Github: https://github.com/ccchieh
 * @Email: email@zzj.cool
 * @Date: 2021-01-18 16:10:17
 * @LastEditors: ccchieh
 * @LastEditTime: 2021-05-24 18:43:43
 */
package goutils

import (
	"os"
	"path"
	"time"

	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap/zapcore"
)

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: GetWriteSyncer
//@description: zap logger中加入file-rotatelogs
//@return: zapcore.WriteSyncer, error

func GetWriteSyncer(director string, logInConsole bool) (zapcore.WriteSyncer, error) {
	fileWriter, err := zaprotatelogs.New(
		path.Join(director, "%Y-%m-%d.log"),
		zaprotatelogs.WithMaxAge(7*24*time.Hour),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	)
	if director == "" && logInConsole {
		return zapcore.AddSync(os.Stdout), err
	}
	if director != "" && logInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}
