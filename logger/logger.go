package logger

import (
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

type Option struct {
	Format   string // text, json
	Level    string // 日志级别
	LogFile  string // 日志路径
	Maxage   int    // 归档文件最大保留天数
	Compress bool   // 是否压缩归档文件
}

func NewLogger(opt *Option) {
	var encoder zapcore.Encoder
	if opt.Format == "text" {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	} else {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}

	core := zapcore.NewCore(encoder, getLogWriter(opt), parseLogLevel(opt.Level))

	logger = zap.New(core)
	logger.Info("Logger init success", zap.String("format", opt.Format))
}

func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	}
	return zapcore.InfoLevel
}

func getLogWriter(opt *Option) zapcore.WriteSyncer {
	writer, err := rotatelogs.New(
		opt.LogFile+".%Y%m%d",
		rotatelogs.WithLinkName(opt.LogFile),                          // 软链接
		rotatelogs.WithRotationTime(24*time.Hour),                     // 按天切割日志
		rotatelogs.WithMaxAge(time.Duration(opt.Maxage)*24*time.Hour), // 最大保留时间
	)
	if err != nil {
		panic(err)
	}
	ws := io.MultiWriter(writer, os.Stdout)
	return zapcore.AddSync(ws)
}

func GetLogger() *zap.Logger {
	return logger
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Sync() {
	logger.Sync()
}
