package rpc

import (
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ginBodyLogger struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (g *ginBodyLogger) Write(b []byte) (int, error) {
	g.body.Write(b)
	return g.ResponseWriter.Write(b)
}

func RequestLoggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ginBodyLogger := &ginBodyLogger{
			body:           bytes.Buffer{},
			ResponseWriter: ctx.Writer,
		}
		ctx.Writer = ginBodyLogger
		ctx.Next()
		val, ok := ctx.Get("method")
		if !ok {
			logger.Infof("status: %d", ctx.Writer.Status())
		}
		color := "\033[42m"
		if ctx.Writer.Status() != 200 {
			color = "\033[41m"
		}
		logger.Infof("status: %s %d \033[0m | rpc_method: \033[100m %s \033[0m", color, ctx.Writer.Status(), val)
	}
}
