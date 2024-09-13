package utils

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitializeLogger(logLevel zapcore.Level, outputPaths []string) error {
	config := zap.NewProductionConfig()
	config.Level.SetLevel(logLevel)
	config.OutputPaths = outputPaths

	var err error
	logger, err = config.Build()
	return err
}

func GetLogger() *zap.Logger {
	if logger == nil {
		_ = InitializeLogger(zapcore.InfoLevel, []string{"stdout"})
	}
	return logger
}

func LogRequest(r *http.Request, traceID string) {
	logger.Info("Incoming request",
		zap.String("method", r.Method),
		zap.String("url", r.URL.Path),
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("trace_id", traceID),
		zap.Time("timestamp", time.Now()),
	)
}

func LogError(r *http.Request, traceID, msg string, err error) {
	logger.Error("Request error",
		zap.String("method", r.Method),
		zap.String("url", r.URL.Path),
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("trace_id", traceID),
		zap.String("message", msg),
		zap.Error(err),
		zap.Time("timestamp", time.Now()),
	)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}
