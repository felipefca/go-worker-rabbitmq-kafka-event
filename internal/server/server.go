package server

import (
	"context"
	"go-worker-rabbitmq-kafka-event/configs"
	"go-worker-rabbitmq-kafka-event/internal/appctx"
	"go-worker-rabbitmq-kafka-event/internal/consumer"
	"go-worker-rabbitmq-kafka-event/internal/db/mongodb"
	"go-worker-rabbitmq-kafka-event/internal/db/redisdb"
	"go-worker-rabbitmq-kafka-event/internal/services"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Server interface {
	Start()
}

type ServerOptions struct {
	Logger    *zap.Logger
	Context   context.Context
	AmqpConn  *amqp.Connection
	MongoConn *mongo.Client
	RedisConn *redis.Client
}

type server struct {
	ServerOptions
}

func NewServer(opt ServerOptions) Server {
	return server{
		ServerOptions: opt,
	}
}

func (s server) Start() {
	logger := appctx.FromContext(s.ServerOptions.Context)

	eventMongoDB := mongodb.NewEventDB(s.MongoConn)
	eventRedisDB := redisdb.NewEventRedisDB(s.RedisConn)
	service := services.NewService(eventMongoDB, eventRedisDB)

	logger.Info("Starting consumer...")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := s.startRabbitConsumer(s.Context, service); err != nil {
			logger.Error(err.Error())
			panic(err)
		}
	}()

	wg.Wait()
}

func (s server) startRabbitConsumer(ctx context.Context, service services.Service) error {
	channel, err := s.ServerOptions.AmqpConn.Channel()
	if err != nil {
		return err
	}
	defer s.ServerOptions.AmqpConn.Close()
	defer channel.Close()

	mq, err := consumer.NewRabbitMQConsumer(channel, configs.GetConfig().RabbitMQ, service)
	if err != nil {
		return err
	}

	return mq.Listen(ctx)
}
