package zhttp

import (
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/zzjcool/goutils/ferr"
)

type baseResult struct {
	Time      int64  `json:"time,omitempty"`
	RequestID string `json:"requestId,omitempty"`
	Code      int    `json:"code"`
}

type okResult struct {
	baseResult
	Data interface{} `json:"data,omitempty"`
}

type errorResult struct {
	baseResult
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Result(c *gin.Context, data interface{}, err error) {

	if c.IsAborted() {
		return
	}
	requestID := ""
	if reqID, exists := c.Get(RequestIDKey); exists {
		requestID = reqID.(string)
	}
	log := GetLogFromGinContext(c)

	if err != nil {
		e, ok := err.(StatusItf)
		if !ok {
			fe := NewE(err, ErrUnknown)
			e = fe
			log.Debug(fe.TraceStack())
		} else {
			fe := ferr.Convert(err)
			log.Debug(fe.TraceStack())
		}
		result := errorResult{
			baseResult{
				Time:      time.Now().Unix(),
				RequestID: requestID,
				Code:      e.Code(),
			},
			e.Message(),
			data,
		}
		c.JSON(e.HttpStatusCode(), result)
		c.Abort()

	} else {
		result := okResult{baseResult{
			Time:      time.Now().Unix(),
			RequestID: requestID,
			Code:      0}, data}
		c.JSON(http.StatusOK, result)
		c.Abort()
	}
}
