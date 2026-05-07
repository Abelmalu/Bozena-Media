package platform

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zapLogger *zap.Logger
}

func InitZapLogger() *Logger {

	var err error
	var logger *zap.Logger

	env := os.Getenv("APP_ENV")

	if env == "production" {

		logger, err = zap.NewProduction()

	} else {

		//logger, err = zap.NewDevelopment()
		config := zap.NewDevelopmentConfig() 
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		logger, err = config.Build()
		

	}

	if err != nil {

		panic(err)
	}

	return &Logger{
		zapLogger: logger,
	}

}

func (logger *Logger) Info(message string, fields ...zap.Field) {
	logger.zapLogger.Info(message, fields...)
}

func (logger *Logger) Error(message errors, fields ...zap.Field) {

	logger.zapLogger.Error(message, fields...)

}

func (logger *Logger) Warn(message string, fields ...zap.Field) {
	logger.zapLogger.Warn(message, fields...)
}

func (logger *Logger) Debug(message string, fields ...zap.Field) {
	logger.zapLogger.Debug(message, fields...)
}

func (logger *Logger) RequestStart(method, path, userAgent, sourceApp string, requestID string) {
	logger.zapLogger.Info("REQUEST_START",
		zap.String("method", method),
		zap.String("path", path),
		zap.String("user_agent", userAgent),
		zap.String("source_app", sourceApp),
		zap.String("request_id", requestID),
	)
}

func (l *Logger) RequestEnd(method, path string, statusCode int, duration time.Duration, requestID string) {
	l.zapLogger.Info("REQUEST_END",
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status_code", statusCode),
		zap.Duration("duration", duration),
		zap.String("request_id", requestID),
	)
}

// Sync flushes buffered logs (call this on app shutdown)
func (l *Logger) Sync() {
	_ = l.zapLogger.Sync()
}
