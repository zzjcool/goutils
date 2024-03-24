package zmiddleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/zzjcool/goutils/zhttp"
	"go.uber.org/zap"
)

// 宕机恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				zap.L().Error("panic recovered", zap.Any("error", err))
				zhttp.Result(c, err, errors.New("server panic"))
			}
		}()
		c.Next()
	}
}
