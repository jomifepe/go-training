package logging

import (
	"context"
	"fmt"
	"github.com/jomifepe/gin_api/util"
	"github.com/sirupsen/logrus"
	gormLogger "gorm.io/gorm/logger"
	"strings"
	"time"
)

type CustomGORMLogger struct {
	LogLevel      gormLogger.LogLevel
	SlowThreshold time.Duration
}

func NewGORMLogger(logrusLevel string) gormLogger.Interface {
	var level gormLogger.LogLevel

	switch strings.ToLower(logrusLevel) {
	case "panic", "fatal", "error":
		level = gormLogger.Error
	case "warn", "warning", "info":
		level = gormLogger.Warn
	case "debug":
		fallthrough
	default:
		level = gormLogger.Info
	}

	return &CustomGORMLogger{
		LogLevel:      level,
		SlowThreshold: 100 * time.Millisecond,
	}
}

// LogMode log mode
func (l *CustomGORMLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info print info
func (l CustomGORMLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Info {
		Logger.Infof("[DB|GORM] " + msg, data...)
	}
}

// Warn print warn messages
func (l CustomGORMLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		Logger.Warnf("[DB|GORM] " + msg, data...)
	}
}

// Error print error messages
func (l CustomGORMLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		Logger.Errorf("[DB|GORM] " + msg, data...)
	}
}

// Trace print sql message
func (l CustomGORMLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= gormLogger.Error:
			sql, rows := fc()
			Logger.WithFields(util.OmitEmptyFields(logrus.Fields{
				"error": err,
				"time": fmt.Sprintf("%.2fs", float64(elapsed.Nanoseconds()) / 1e6),
				"rows": rows,
			})).Errorln("[DB|GORM]", sql)
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLogger.Warn:
			sql, rows := fc()
			Logger.WithFields(util.OmitEmptyFields(logrus.Fields{
				"error": err,
				"time": fmt.Sprintf("%.2fs", float64(elapsed.Nanoseconds()) / 1e6),
				"rows": rows,
			})).Warnln("[DB|GORM]", sql)
		case l.LogLevel >= gormLogger.Info:
			sql, rows := fc()
			Logger.WithFields(util.OmitEmptyFields(logrus.Fields{
				"error": err,
				"time": fmt.Sprintf("%.2fs", float64(elapsed.Nanoseconds()) / 1e6),
				"rows": rows,
			})).Infoln("[DB|GORM]", sql)
		}
	}
}

func getFields(data ...interface{}) logrus.Fields {
	switch data[0] {
	case "sql":
		return logrus.Fields{
			"module":  "gorm",
			"type":    "sql",
			"rows":    data[5],
			"src_ref": data[1],
			"values":  data[4],
		}
	}

	return logrus.Fields{
		"module": "gorm",
		"type": "log",
	}
}