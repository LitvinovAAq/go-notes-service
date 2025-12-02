package midleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "request_id"

func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rid := ctx.GetHeader("X-Request-ID")
		if rid == "" {
			rid = uuid.NewString()
		}

		ctx.Set(RequestIDKey, rid)
		ctx.Writer.Header().Set("X-Request-ID", rid)

		ctx.Next()
	}
}
