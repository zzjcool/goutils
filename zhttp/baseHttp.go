package zhttp

import (
	"github.com/gin-gonic/gin"
	"github.com/zzjcool/ginhelper"
	"github.com/zzjcool/goutils/defaults"
	"github.com/zzjcool/goutils/ferr"
)

type BaseParam struct {
	ginhelper.BaseParam
	RequestID string `header:"requestID"`
}

func (param *BaseParam) Bind(c *gin.Context, p ginhelper.Parameter) error {
	err := defaults.Apply(p)
	if err != nil {
		return NewE(ferr.Convert(err), ErrServerPanic)
	}
	if err := c.ShouldBind(p); err != nil {
		return NewE(ferr.Convert(err), ErrParameterMatch)
	}
	return nil
}

func (param *BaseParam) Result(c *gin.Context, data ginhelper.Data, err error) {
	Result(c, data, err)
}
