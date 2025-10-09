package midleware

import (
	"myproject/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)
		status := ctx.Writer.Status()
		method := ctx.Request.Method
		path := ctx.FullPath()

		ridAny, _ := ctx.Get(RequestIDKey)
		rid, _ := ridAny.(string)

		// log.Printf("request_id=%s method=%s path=%s status=%d duration=%v", rid, method, path,status, duration)
		logger.Infof("request_id=%s method=%s path=%s status=%d duration=%v", rid, method, path,status, duration)
	}
}