package zhttp

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zzjcool/goutils"
	"github.com/zzjcool/goutils/defaults"
	"github.com/zzjcool/goutils/str"
	"github.com/zzjcool/goutils/zlog"
	"go.uber.org/zap"
)

type loadRouterFunc func(route *gin.Engine)

// HttpServerConfig 启动一个Http服务的配置
type HttpServerConfig struct {
	Name       string `default:"httpserver"`
	Port       uint   `default:"12321"`
	Https      bool
	Cert       string
	Key        string
	RouterFunc loadRouterFunc
}

var once sync.Once



var log zlog.Logger = zap.S()

func SetLogger(l zlog.Logger) {
	log = l
}
func Serve(c *HttpServerConfig) error {
	if err := defaults.Apply(c); err != nil {
		log.Error("apply default", zap.Error(err))
		return err
	}

	once.Do(func() {
		gin.SetMode(gin.ReleaseMode)
		gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
			log.Debug(httpMethod + "\t" + absolutePath)
		}
	})
	r := gin.New()
	c.RouterFunc(r)
	if c.Https {
		return serveTLS(c.Port, r, c.Cert, c.Key)
	}
	return serve(c.Port, r)
}

func serve(port uint, r http.Handler) error {
	if !goutils.LegalPort(port) {
		log.Panic("Port range error", zap.String("error", "Port range error"))
	}
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Info("Service running address: http://" + addr)
	s := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return s.ListenAndServe()
}

func serveTLS(port uint, r http.Handler, cert, key string) error {
	if !goutils.LegalPort(port) {
		log.Panic("Port range error", zap.String("error", "Port range error"))
	}
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Info("tls address: https://" + addr)
	s := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	certFile, keyFile := str.TempFile(cert), str.TempFile(key)

	log.Info("certFile: ", certFile, " keyFile: ", keyFile)
	return s.ListenAndServeTLS(certFile, keyFile)
}
