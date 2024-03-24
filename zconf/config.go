package zconf

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/zzjcool/goutils/defaults"
	"go.uber.org/zap"
)

// defaultConfig 获取默认配置
func defaultConfig[T any](conf T) error {

	err := defaults.Apply(conf)
	if err != nil {
		return errors.Join(err, ErrSetDefaultValue)
	}
	return nil
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
}

var log Logger = zap.S()

func SetLogger(l Logger) {
	log = l
}

// Load 加载配置
func Load[T any](conf T) error {
	return LoadCustom(conf, "", "config.yml","", NewDeepFs(0))
}

func LoadWithDir[T any](conf T, dir string) error {
	return LoadCustom(conf, "", "config.yml", dir, NewDeepFs(0))
}

func LoadWithEnv[T any](conf T, env string) error {
	return LoadCustom(conf, env, "config.yml", "",NewDeepFs(0))
}

func LoadCustom[T any](conf T, env, configFile, dir string, cfs fs.FS) error {
	defaultConfig(conf)
	return load(conf, env, configFile, dir, cfs)
}

func load[T any](conf T, env, configFile, dir string, confFs fs.FS) error {
	var vp *viper.Viper

	if err := loadConfig(vp, configFile, dir, conf, confFs); err != nil {
		return err
	}
	if env != "" {
		return loadConfig(vp, env+"."+configFile, dir, conf, confFs)
	}
	return nil
}

func loadConfig[T any](v *viper.Viper, config string, dir string, conf T, confFs fs.FS) error {
	log.Debugf("CONFIG File:%v\n", config)
	if v == nil {
		v = viper.New()
	}
	// 设置CONF前缀的env
	v.SetEnvPrefix("CONF")
	v.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	v.SetConfigFile(config)
	rawConf, err := confFs.Open(path.Join(dir, config))

	if err != nil {
		log.Debug(err)
		log.Debug("The config file is not loaded:", config)
		return errors.Join(err, ErrOpenConfigFile)
	}
	err = v.ReadConfig(rawConf)
	if err != nil {
		log.Debugf("read config error:%v", err)
		return errors.Join(err, ErrInvalidConfigFile)
	}
	if err := v.Unmarshal(conf); err != nil {
		log.Debug(err)
		return errors.Join(err, ErrUnmarshalConfig)
	}

	return nil
}

// NewDeepFs ...
func NewDeepFs(deep int) fs.FS {
	return &DeepFs{Deep: deep}
}

// DeepFs 配置
type DeepFs struct {
	Deep int // 表示查找配置的目录深度，例如为1时，当前目录为/code/config/,这个时候配置的搜寻目录从/code/开始
}

// Open 打开配置文件
func (fs *DeepFs) Open(name string) (fs.File, error) {
	path, _ := os.Getwd()
	deep := fs.Deep
	for deep > 0 {
		path = filepath.Dir(path)
		deep--
	}
	file := filepath.Join(path, name)
	log.Debugf("try to open file:", file)
	return os.Open(file)
}

