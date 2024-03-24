package zmiddleware

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zzjcool/goutils"
	"github.com/zzjcool/goutils/zhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func LoggerMiddleware(isDebug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(zhttp.RequestIDKey)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		log := zap.S().Named(requestID)
		c.Set(zhttp.RequestIDKey, requestID)
		c.Set(zhttp.RequestIDLogKey, log)

		if isDebug {
			data, err := c.GetRawData()
			if err != nil {
				log.Error(err)
			}
			body := string(data)

			c.Request.Body = io.NopCloser(bytes.NewBuffer(data)) // copy body

			blw := &ResponseWriterWrapper{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw

			start := time.Now()
			defer func() {
				data, _ := io.ReadAll(blw.Body)
				log.Debugf(`
IP:         %s
Header:     %s
Method:     %s
RequestURI: %s
Comsuming:  %s
HttpStatus: %d
Body:       %s
Return:     %s`,
					goutils.RemoteIp(c.Request),
					c.Request.Header,
					c.Request.URL.Port(),
					c.Request.URL.Scheme,
					c.Request.Method,
					c.Request.RequestURI,
					time.Since(start),
					c.Writer.Status(),
					body,
					string(data))

			}()
		}

		c.Next()
	}
}

func UnaryServerInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
	resp interface{}, err error) {

	md, ok := metadata.FromIncomingContext(ctx)
	rawId, ok2 := md[zhttp.GRPCRequestIDKey]
	var requestID string
	if !ok || !ok2 {
		requestID = uuid.NewString()
	} else {
		requestID = rawId[0]
	}

	log := zap.S().Named(requestID)

	ctx = zhttp.SetCtxLog(ctx, log)

	remote, _ := peer.FromContext(ctx)
	remoteAddr := remote.Addr.String()

	start := time.Now()
	defer func() {
		log.Debug(fmt.Sprintf("%s %s",
			remoteAddr,
			time.Since(start),
		))
	}()

	header := metadata.Pairs(zhttp.GRPCRequestIDKey, requestID)
	if err := grpc.SendHeader(ctx, header); err != nil {
		log.Error(err.Error())
	}

	return handler(ctx, req)
}

type ResponseWriterWrapper struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w ResponseWriterWrapper) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w ResponseWriterWrapper) WriteString(s string) (int, error) {
	w.Body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
