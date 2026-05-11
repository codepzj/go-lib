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

func NewLogger(opt *Option) {
	opt = withDefault(opt)

	var encoder zapcore.Encoder
	if opt.Format == "text" {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	} else {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}

	core := zapcore.NewCore(encoder, getLogWriter(opt), parseLogLevel(opt.Level))

	logger = zap.New(core)
	logger.Info("Logger init success", zap.String("format", string(opt.Format)))
}

func parseLogLevel(level LogLevel) zapcore.Level {
	switch level {
	case DEBUG:
		return zapcore.DebugLevel
	case INFO:
		return zapcore.InfoLevel
	case WARN:
		return zapcore.WarnLevel
	case ERROR:
		return zapcore.ErrorLevel
	}
	return zapcore.InfoLevel
}

func getLogWriter(opt *Option) zapcore.WriteSyncer {
	if !opt.Output.EnableFile {
		return zapcore.AddSync(os.Stdout)
	}

	writer, err := rotatelogs.New(
		opt.Output.FilePath+".%Y%m%d",
		rotatelogs.WithLinkName(opt.Output.FilePath),                         // 软链接
		rotatelogs.WithRotationTime(24*time.Hour),                            // 按天切割日志
		rotatelogs.WithMaxAge(time.Duration(opt.Output.MaxAge)*24*time.Hour), // 最大保留天数
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
