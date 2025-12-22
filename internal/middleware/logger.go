package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := uuid.NewString()
		ctx.Set("trace_id", id)
		ctx.Writer.Header().Set("X-Request-ID", id)
		ctx.Next()
	}
}

func Logger() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		zap.S().Infow("request",
			"method", ctx.Request.Method,
			"path", ctx.Request.URL.Path,
			"status", ctx.Writer.Status(),
			"latency", time.Since(start),
			"trace_id", ctx.GetString("trace_id"),
			"user_id", ctx.GetString("user_id"),
		)
	}
}