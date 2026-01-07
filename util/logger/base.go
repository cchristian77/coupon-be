package logger

import (
	"base_project/util/constant"
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func Initialise() *zap.Logger {
	fileLoggerConfig := zap.NewProductionEncoderConfig()
	fileLoggerConfig.MessageKey = "message"
	fileLoggerConfig.LevelKey = "level"
	fileLoggerConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	fileLoggerConfig.TimeKey = "timestamp"
	fileLoggerConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileLoggerConfig.CallerKey = "caller"
	fileLoggerConfig.EncodeCaller = zapcore.ShortCallerEncoder
	fileLoggerConfig.FunctionKey = "func"
	logFile, _ := os.OpenFile("logs/errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	core := zapcore.NewTee(
		// logger to record in warn level (including errors) to errors.log
		zapcore.NewCore(
			zapcore.NewJSONEncoder(fileLoggerConfig),
			zapcore.AddSync(logFile),
			zapcore.WarnLevel,
		),
		// logger to record in debug level in terminal
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		),
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger
}

// L call logger instance without the need for context
func L() *zap.Logger {
	return logger
}

func Fatal(ctx context.Context, message string, args ...any) {
	logger.Fatal(fmt.Sprintf(message, args...), getZapFieldsFromCtx(ctx)...)
}

func Error(ctx context.Context, message string, args ...any) {
	logger.Error(fmt.Sprintf(message, args...), getZapFieldsFromCtx(ctx)...)
}

func Warn(ctx context.Context, message string, args ...any) {
	logger.Warn(fmt.Sprintf(message, args...), getZapFieldsFromCtx(ctx)...)
}

func Info(ctx context.Context, message string, args ...any) {
	logger.Info(fmt.Sprintf(message, args...), getZapFieldsFromCtx(ctx)...)
}

func Debug(ctx context.Context, message string, args ...any) {
	logger.Debug(fmt.Sprintf(message, args...), getZapFieldsFromCtx(ctx)...)
}

func getZapFieldsFromCtx(ctx context.Context) []zapcore.Field {
	correlationID := constant.CorrelationIDFromCtx(ctx)
	authUser := constant.AuthUserFromCtx(ctx)

	var userID uint64
	if authUser != nil {
		userID = authUser.ID
	}

	var fields []zapcore.Field

	if correlationID != "" {
		fields = append(fields, zap.String(constant.XCorrelationIDKey, correlationID))
	}

	if userID != 0 {
		fields = append(fields, zap.Uint64(constant.XUserIDKey, userID))
	}

	return fields
}
