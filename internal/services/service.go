package services

import (
	"context"
	"go-worker-rabbitmq-kafka-event/internal/appctx"
	"go-worker-rabbitmq-kafka-event/internal/db/mongodb"
	"go-worker-rabbitmq-kafka-event/internal/db/redisdb"
	"go-worker-rabbitmq-kafka-event/internal/models"
)

type Service interface {
	HandlerEvent(ctx context.Context, event models.Event) error
}

type service struct {
	eventDB      mongodb.EventDB
	eventRedisDB redisdb.EventRedisDB
}

func NewService(eventDB mongodb.EventDB, eventRedisDB redisdb.EventRedisDB) Service {
	return &service{
		eventDB:      eventDB,
		eventRedisDB: eventRedisDB,
	}
}

func (s service) HandlerEvent(ctx context.Context, event models.Event) error {
	logger := appctx.FromContext(ctx)

	_, err := s.eventDB.AddEvent(ctx, event)
	if err != nil {
		return err
	}

	err = s.eventRedisDB.AddEvent(ctx, event)
	if err != nil {
		return err
	}

	logger.Info("Success to add event!")
	return nil
}
