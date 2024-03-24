package zlog_test

import (
	"testing"

	"github.com/zzjcool/goutils/zlog"
)

func TestLogger(t *testing.T) {
	log := zlog.NewTest()
	log.Info("Info")
	log.Debug("Debug")
	log.Error("Error")
	defer func() {
		if err := recover(); err != nil {
			return
		}
		t.Fatal("should panic")
	}()
	log.Panic("Panic")
}
