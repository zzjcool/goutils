package zlog

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/zzjcool/goutils"
	"github.com/zzjcool/goutils/defaults"

	"time"
)

var level zapcore.Level

var logger *zap.Logger

func New(conf *LogConf) (l *zap.SugaredLogger) {
	defer func() {
		zap.ReplaceGlobals(logger)
	}()

	if ok := goutils.FileIsExist(conf.Director); !ok { // 判断是否有Director文件夹
		_ = os.Mkdir(conf.Director, os.ModePerm)
	}

	switch conf.Level { // 初始化配置文件的Level
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	logger = zap.New(getEncoderCore(conf), zap.AddStacktrace(zap.ErrorLevel))

	if conf.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger.Sugar()
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig(conf *LogConf) (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  conf.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder(conf),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   CustomCallerEncoder(conf),
	}
	switch {
	case conf.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case conf.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case conf.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case conf.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder(conf *LogConf) zapcore.Encoder {
	if conf.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig(conf))
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig(conf))
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore(conf *LogConf) (core zapcore.Core) {
	writer, err := goutils.GetWriteSyncer(conf.Director, conf.LogInConsole) // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}
	return zapcore.NewCore(getEncoder(conf), writer, level)
}

// 自定义日志输出时间格式
func CustomTimeEncoder(conf *LogConf) func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		if conf.Prefix != "" {
			enc.AppendString(conf.Prefix + ":")
		}
		enc.AppendString(t.Format("20060102-15:04:05.000"))
	}
}

// 自定义日志路径
func CustomCallerEncoder(conf *LogConf) func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	return func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		if conf.Trimmed {
			enc.AppendString(caller.TrimmedPath())
		} else {
			enc.AppendString(caller.String())
		}
	}
}

func NewTest()  (l *zap.SugaredLogger){
	conf := new(LogConf)
	err := defaults.Apply(conf)
	if err != nil {
		panic(err)
	}
	conf.Director = ""
	return New(conf)
}
