package zconf_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/zzjcool/goutils/zconf"
	"gotest.tools/assert"
)

func ExampleLoad() {
	conf := new(ConfigTest)
	err := zconf.Load(conf)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf.Setting)
}

func CreateTmpYamlFile(content, filePath string) func() {
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		fmt.Println("Error writing config file:", err)
		panic(err)
	}

	return func() {
		os.Remove(filePath)
	}
}

type ConfigTest struct {
	DefaultFild bool `default:"true"`
	Setting     int  `yaml:"setting"`
}

func TestConfig(t *testing.T) {

	t.Run("Test default value", func(t *testing.T) {

		defer CreateTmpYamlFile(`setting: 123`, "config.yml")()

		conf := new(ConfigTest)
		err := zconf.Load(conf)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, conf.DefaultFild, true)
		assert.Equal(t, conf.Setting, 123)

	})

	t.Run("Test error yaml", func(t *testing.T) {

		defer CreateTmpYamlFile(`setting: eee`, "config.yml")()
		conf := new(ConfigTest)
		err := zconf.Load(conf)
		assert.Equal(t, errors.Is(err, zconf.ErrUnmarshalConfig), true)
	})
	t.Run("Test yaml formatted", func(t *testing.T) {

		defer CreateTmpYamlFile(`setting:eee`, "config.yml")()
		conf := new(ConfigTest)
		err := zconf.Load(conf)
		assert.Equal(t, errors.Is(err, zconf.ErrInvalidConfigFile), true)
	})

	t.Run("Test yaml not exist", func(t *testing.T) {
		conf := new(ConfigTest)
		err := zconf.Load(conf)
		assert.Equal(t, errors.Is(err, os.ErrNotExist), true)
		assert.Equal(t, errors.Is(err, zconf.ErrOpenConfigFile), true)
	})
}
