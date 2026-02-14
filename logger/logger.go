package logger

import (
	"io"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

type Option struct {
	Type       string // text, json
	Level      string // 日志级别
	LogFile    string // 日志路径
	MaxSize    int    // 日志文件最大(MB)
	MaxBackups int    // 归档文件最大保留数目
	MaxAge     int    // 归档文件最大保留时间(天)
	Compress   bool   // 是否压缩
}

func NewLogger(opt *Option) {
	var encoder zapcore.Encoder
	if opt.Type == "text" {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	} else {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}

	core := zapcore.NewCore(encoder, getLogWriter(opt), parseLogLevel(opt.Level))

	logger = zap.New(core)
	log.Printf("Logger init success[%s]...\n", opt.Type)
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
	logger := &lumberjack.Logger{
		Filename:   opt.LogFile,
		MaxSize:    opt.MaxSize,
		MaxBackups: opt.MaxBackups,
		MaxAge:     opt.MaxAge,
		Compress:   opt.Compress,
	}
	ws := io.MultiWriter(logger, os.Stdout)
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
