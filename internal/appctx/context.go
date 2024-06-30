package appctx

import (
	"context"
	"go-worker-rabbitmq-kafka-event/internal/constants"

	"go.uber.org/zap"
)

var (
	LoggerKey = "logger"
)

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	logger := ctx.Value(LoggerKey).(*zap.Logger)
	return logger
}

func SetCorrelationId(ctx context.Context, correlationId string) context.Context {
	return context.WithValue(ctx, constants.CorrelationIdHeader, correlationId)
}

func GetCorrelationId(ctx context.Context) string {
	correlationId := ctx.Value(constants.CorrelationIdHeader)
	if correlationId != nil {
		return correlationId.(string)
	}
	return ""
}
