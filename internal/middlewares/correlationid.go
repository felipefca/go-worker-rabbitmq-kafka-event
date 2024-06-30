package middlewares

import (
	"context"
	"go-worker-rabbitmq-kafka-event/internal/appctx"
	"strings"

	"github.com/google/uuid"
)

func CorrelationIdMiddleware(ctx context.Context, correlationId string) context.Context {
	if strings.TrimSpace(correlationId) == "" {
		correlationId = uuid.NewString()
	}

	return appctx.SetCorrelationId(ctx, correlationId)
}
