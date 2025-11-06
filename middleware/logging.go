package middleware

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(log *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqID := uuid.NewString()
			requestLogger := log.WithFields(logrus.Fields{
				"method":     c.Request().Method,
				"path":       c.Path(),
				"request_id": reqID,
			})

			c.Set("logger", requestLogger)
			c.Response().Header().Set("X-Request-ID", reqID)

			requestLogger.Info("Request started")

			start := time.Now()

			err := next(c)

			duration := time.Since(start)
			status := c.Response().Status

			logEntry := requestLogger.WithFields(logrus.Fields{
				"status_code": status,
				"duration_ms": duration.Milliseconds(),
			})

			if err != nil {
				logEntry.WithError(err).Error("Request failed")
			} else if status >= 400 && status < 500 {
				logEntry.Warn("Request completed (Client Error)")
			} else {
				logEntry.Info("Request finished successfully")
			}

			return err
		}
	}
}
