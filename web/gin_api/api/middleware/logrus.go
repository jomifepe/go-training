package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jomifepe/gin_api/util"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

// Logger is a logrus logging middleware for the gin framework
func Logger(logger logrus.FieldLogger) gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}

	return func(c *gin.Context) {
		var (
			path = c.Request.URL.Path
			start = time.Now()
		)
		c.Next()

		var (
			stop = time.Since(start)
			latency = stop.Nanoseconds()
			statusCode = c.Writer.Status()
			clientIP = c.ClientIP()
			clientUserAgent = c.Request.UserAgent()
			referer = c.Request.Referer()
			dataLength = c.Writer.Size()
		)

		entry := logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       path,
			"statusCode": statusCode,
			"hostname":   hostname,
			"clientIP":   clientIP,
			"latency":    util.GetTimeString(latency),
			"referer":    referer,
			"dataLength": util.GetSizeString(dataLength),
			"userAgent":  clientUserAgent,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
			return
		}

		//msg := fmt.Sprintf(
		//	"%s - %s [%s] \"%s %s\" %d %d \"%s\" \"%s\" (%dms)",
		//	clientIP, hostname, time.Now().Format(time.RFC3339), c.Request.Method,
		//	path, statusCode, dataLength, referer, clientUserAgent, latency,
		//)

		msg := "[API] Received HTTP request"

		if statusCode >= http.StatusInternalServerError {
			entry.Error(msg)
		} else if statusCode >= http.StatusBadRequest {
			entry.Warn(msg)
		} else {
			entry.Info(msg)
		}
	}
}