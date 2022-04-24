package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	gl "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"

	"github.com/2754github/gorm-with-logrus/gorm"
)

// Replace gorm logger with logrus.
// https://github.com/go-gorm/gorm/blob/master/logger/logger.go
// https://github.com/onrik/gorm-logrus/blob/master/logger.go
type logger struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	LogLevel                  LogLevel
	infoStr                   string
	warnStr                   string
	errStr                    string
	traceStr                  string
	traceWarnStr              string
	traceErrStr               string
}

type LogLevel int

const (
	Silent LogLevel = iota + 1
	Error
	Warn
	Info
)

func New(slowThreshold time.Duration, ignoreRecordNotFoundError bool, logLevel LogLevel) *logger {
	return &logger{
		SlowThreshold:             slowThreshold,
		IgnoreRecordNotFoundError: ignoreRecordNotFoundError,
		LogLevel:                  logLevel,
		infoStr:                   "%s\n[info] ",
		warnStr:                   "%s\n[warn] ",
		errStr:                    "%s\n[error] ",
		traceStr:                  "%s\n[%.3fms] [rows:%v] %s",
		traceWarnStr:              "%s %s\n[%.3fms] [rows:%v] %s",
		traceErrStr:               "%s %s\n[%.3fms] [rows:%v] %s",
	}
}

// Do nothing.
func (l *logger) LogMode(level gl.LogLevel) gl.Interface {
	return l
}

func (l logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= Info {
		logrus.Infof(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (l logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= Warn {
		logrus.Warnf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (l logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= Error {
		logrus.Errorf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			logrus.Errorf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logrus.Errorf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			logrus.Warnf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logrus.Warnf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == Info:
		sql, rows := fc()
		if rows == -1 {
			logrus.Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logrus.Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
