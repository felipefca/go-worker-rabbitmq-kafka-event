package server

import (
	"context"
	"go-worker-rabbitmq-kafka-event/configs"
	"go-worker-rabbitmq-kafka-event/internal/appctx"
	"go-worker-rabbitmq-kafka-event/internal/db/mongodb"
	"go-worker-rabbitmq-kafka-event/internal/db/redisdb"
	"go-worker-rabbitmq-kafka-event/internal/infra/kafka"
	"go-worker-rabbitmq-kafka-event/internal/infra/rabbitmq"
	"go-worker-rabbitmq-kafka-event/internal/services"
	"sync"

	"github.com/IBM/sarama"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Server interface {
	Start()
}

type ServerOptions struct {
	Logger        *zap.Logger
	Context       context.Context
	AmqpConn      *amqp.Connection
	MongoConn     *mongo.Client
	RedisConn     *redis.Client
	KafkaConsumer sarama.ConsumerGroup
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
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := s.startRabbitConsumer(s.Context, service); err != nil {
			logger.Error(err.Error())
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := s.startKafkaConsumer(s.Context, service); err != nil {
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

	mq, err := rabbitmq.NewRabbitMQConsumer(channel, configs.GetConfig().RabbitMQ, service)
	if err != nil {
		return err
	}

	return mq.Listen(ctx)
}

func (s server) startKafkaConsumer(ctx context.Context, service services.Service) error {
	consumer, err := kafka.NewKafkaConsumer(s.KafkaConsumer, configs.GetConfig().Kafka)
	if err != nil {
		return err
	}

	return consumer.Listen(ctx)
}
