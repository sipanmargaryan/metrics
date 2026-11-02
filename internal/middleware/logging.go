package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"

	"metrics/internal/logging"
)

// RequestLogger must wrap recovery and error handling to record the final status.
func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id [16]byte
		_, _ = rand.Read(id[:]) // crypto/rand.Read terminates on an entropy source failure.
		requestID := hex.EncodeToString(id[:])
		requestLogger := logger.With("request_id", requestID, "method", c.Request.Method)
		c.Request = c.Request.WithContext(logging.AddLoggerToContext(c.Request.Context(), requestLogger))
		c.Header("X-Request-ID", requestID)
		start := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		level := slog.LevelInfo
		if c.Writer.Status() >= 500 {
			level = slog.LevelError
		} else if c.Writer.Status() >= 400 {
			level = slog.LevelWarn
		}
		requestLogger.LogAttrs(c.Request.Context(), level, "http request",
			slog.String("route", route),
			slog.Int("status", c.Writer.Status()),
			slog.Float64("duration_ms", float64(time.Since(start).Microseconds())/1000),
			slog.Int("response_bytes", c.Writer.Size()))
	}
}
