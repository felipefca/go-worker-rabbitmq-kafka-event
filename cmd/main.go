package main

import (
	"context"
	"fmt"
	"go-worker-rabbitmq-kafka-event/configs"
	"go-worker-rabbitmq-kafka-event/internal/appctx"
	"go-worker-rabbitmq-kafka-event/internal/server"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	logConfig := zap.NewProductionConfig()

	logger, err := logConfig.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	ctx = appctx.WithLogger(ctx, logger)

	amqpConn, err := connectRabbitMQ(ctx)
	if err != nil {
		panic(err)
	}
	defer amqpConn.Close()

	logger.Info("RabbitMQ Connected!")

	mongoConn, err := connectMongoDB(ctx)
	if err != nil {
		panic(err)
	}
	defer mongoConn.Disconnect(ctx)

	logger.Info("Mongo Connected!")

	redisConn, err := connectRedis(ctx)
	if err != nil {
		panic(err)
	}
	defer redisConn.Close()

	logger.Info("Redis Connected!")

	kafkaClient, err := connectKafka()
	if err != nil {
		panic(err)
	}
	defer kafkaClient.Close()

	logger.Info("Kafka Client connected!")

	s := server.NewServer(server.ServerOptions{
		Logger:        logger,
		Context:       ctx,
		AmqpConn:      amqpConn,
		MongoConn:     mongoConn,
		RedisConn:     redisConn,
		KafkaConsumer: kafkaClient,
	})
	s.Start()
}

func connectRabbitMQ(ctx context.Context) (*amqp.Connection, error) {
	logger := appctx.FromContext(ctx)
	countRetry := 0

	for {
		cfg := configs.GetConfig().RabbitMQ
		amqpUri := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", cfg.UserName, cfg.Password, cfg.HostName, cfg.Port, cfg.VirtualHost)

		conn, err := amqp.Dial(amqpUri)
		if err != nil {
			countRetry++
		} else {
			return conn, nil
		}

		if countRetry <= int(cfg.MaxRetry) {
			logger.Error(fmt.Sprintf("fail to connect RabbitMQ. Retry %d%d...", countRetry, cfg.MaxRetry))
			continue
		} else {
			return nil, fmt.Errorf("error to connect RabbitMQ: %w", err)
		}
	}
}

func connectMongoDB(ctx context.Context) (*mongo.Client, error) {
	mongoUri := configs.GetConfig().MongoDB.ConnectionUrl

	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		return nil, fmt.Errorf("error to connect MongoDb. %w", err)
	}

	err = conn.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func connectRedis(ctx context.Context) (*redis.Client, error) {
	cfg := configs.GetConfig().RedisDB

	client := redis.NewClient(&redis.Options{
		Addr:       cfg.Host + ":" + strconv.Itoa(int(cfg.Port)),
		Password:   cfg.Password,
		DB:         0,
		MaxRetries: cfg.ConnectRetry,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("error to connect Redis. %w", err)
	}

	return client, nil
}

func connectKafka() (sarama.ConsumerGroup, error) {
	cfg := configs.GetConfig().Kafka

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	brokers := []string{cfg.Broker}

	consumer, err := sarama.NewConsumerGroup(brokers, cfg.GroupId, config)
	if err != nil {
		panic(err)
	}

	return consumer, nil
}
