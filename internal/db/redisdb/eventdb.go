package redisdb

import (
	"context"
	"encoding/json"
	"go-worker-rabbitmq-kafka-event/configs"
	"go-worker-rabbitmq-kafka-event/internal/appctx"
	"go-worker-rabbitmq-kafka-event/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type EventRedisDB interface {
	AddEvent(ctx context.Context, event models.Event) error
}

type eventRedisDB struct {
	redisConn *redis.Client
}

func NewEventRedisDB(redisConn *redis.Client) EventRedisDB {
	return &eventRedisDB{
		redisConn: redisConn,
	}
}

func (db eventRedisDB) AddEvent(ctx context.Context, event models.Event) error {
	logger := appctx.FromContext(ctx)
	cfg := configs.GetConfig().RedisDB

	key := cfg.KeyPrefix + event.Message

	json, err := json.Marshal(event)
	if err != nil {
		logger.Warn("Fail to insert event")
		return nil
	}

	err = db.redisConn.Set(ctx, key, json, time.Minute*5).Err()
	if err != nil {
		return err
	}

	logger.Info("Event successfully inserted in Redis!")
	return nil
}
