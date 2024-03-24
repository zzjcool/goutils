package zhttp

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetLogFromGinContext(c *gin.Context) *zap.SugaredLogger {
	r, ok := c.Get(RequestIDLogKey)
	if !ok {
		return zap.S()
	}

	return r.(*zap.SugaredLogger)
}

func GetGinCtxAny[T any](c *gin.Context,key string) (ret T) {
	r, ok := c.Get(key)
	if !ok {
		return ret
	}
	return r.(T)
}

type loggerContextKey struct{}

func SetCtxLog(ctx context.Context, log *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, log)
}

func GetCtxLog(ctx context.Context) *zap.SugaredLogger {
	log, ok := ctx.Value(loggerContextKey{}).(*zap.SugaredLogger)
	if !ok {
		return zap.S()
	}
	return log
}

func ErrPanic(err error) {
	if err != nil {
		zap.S().Panic(err)
	}
}
