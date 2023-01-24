package mysql

import (
	"context"
	"errors"
	"fmt"
	xlogger "github.com/go-slark/slark/logger"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type FuncOpts func(l *customizedLogger)

func WithLogLevel(level logger.LogLevel) FuncOpts {
	return func(l *customizedLogger) {
		l.LogLevel = level
	}
}

func WithColor(color bool) FuncOpts {
	return func(l *customizedLogger) {
		l.Colorful = color
	}
}

func WithSlowThreshold(tm time.Duration) FuncOpts {
	return func(l *customizedLogger) {
		l.SlowThreshold = tm
	}
}

func WithRecordNotFound(i bool) FuncOpts {
	return func(l *customizedLogger) {
		l.IgnoreRecordNotFoundError = i
	}
}

func WithLogger(logger xlogger.Logger) FuncOpts {
	return func(l *customizedLogger) {
		l.Logger = logger
	}
}

type customizedLogger struct {
	logger.Config
	xlogger.Logger
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func newCustomizedLogger(opts ...FuncOpts) logger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	cl := &customizedLogger{
		Config: logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: false,
			Colorful:                  false,
		},
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
	for _, opt := range opts {
		opt(cl)
	}

	if cl.Colorful {
		cl.infoStr = logger.Green + "%s\n" + logger.Reset + logger.Green + "[info] " + logger.Reset
		cl.warnStr = logger.BlueBold + "%s\n" + logger.Reset + logger.Magenta + "[warn] " + logger.Reset
		cl.errStr = logger.Magenta + "%s\n" + logger.Reset + logger.Red + "[error] " + logger.Reset
		cl.traceStr = logger.Green + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
		cl.traceWarnStr = logger.Green + "%s " + logger.Yellow + "%s\n" + logger.Reset + logger.RedBold + "[%.3fms] " + logger.Yellow + "[rows:%v]" + logger.Magenta + " %s" + logger.Reset
		cl.traceErrStr = logger.RedBold + "%s " + logger.MagentaBold + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
	}
	return cl
}

func (l *customizedLogger) LogMode(level logger.LogLevel) logger.Interface {
	nl := *l
	nl.LogLevel = level
	return &nl
}

func (l customizedLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel < logger.Info {
		return
	}
	l.Log(ctx, xlogger.InfoLevel, map[string]interface{}{}, fmt.Sprintf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
}

func (l customizedLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel < logger.Warn {
		return
	}
	l.Log(ctx, xlogger.InfoLevel, map[string]interface{}{}, fmt.Sprintf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
}

func (l customizedLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel < logger.Error {
		return
	}
	l.Log(ctx, xlogger.InfoLevel, map[string]interface{}{}, fmt.Sprintf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
}

func (l customizedLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	var param interface{}
	if rows == -1 {
		param = "-"
	} else {
		param = rows
	}
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		l.Log(ctx, xlogger.InfoLevel, map[string]interface{}{}, fmt.Sprintf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, param, sql))

	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		l.Log(ctx, xlogger.InfoLevel, map[string]interface{}{}, fmt.Sprintf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, param, sql))

	case l.LogLevel == logger.Info:
		l.Log(ctx, xlogger.InfoLevel, map[string]interface{}{}, fmt.Sprintf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, param, sql))
	}
}
