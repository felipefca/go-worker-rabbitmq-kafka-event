package kafka

import (
	"context"
	"fmt"
	"go-worker-rabbitmq-kafka-event/configs"
	"go-worker-rabbitmq-kafka-event/internal/appctx"

	"github.com/IBM/sarama"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"golang.org/x/sync/errgroup"
)

type KafkaConsumer interface {
	Listen(ctx context.Context) error
}

type kafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup
	config        configs.Kafka
	schemaClient  schemaregistry.Client
}

func NewKafkaConsumer(consumerGroup sarama.ConsumerGroup, config configs.Kafka) (KafkaConsumer, error) {
	schemaRegistry := schemaregistry.NewConfig(config.SchemaRegistryUrl)

	schemaClient, err := schemaregistry.NewClient(schemaRegistry)
	if err != nil {
		return nil, err
	}

	return &kafkaConsumer{
		consumerGroup: consumerGroup,
		config:        config,
		schemaClient:  schemaClient,
	}, nil
}

func (k *kafkaConsumer) Listen(ctx context.Context) error {
	logger := appctx.FromContext(ctx)

	// partitionConsumer, err := consumer.ConsumePartition(config.TopicName, 0, sarama.OffsetNewest)
	// if err != nil {
	// 	return nil, err
	// }

	consumerHandler := &ConsumerHandler{
		Ready: make(chan bool),
		ctx:   ctx,
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		for {
			err := k.consumerGroup.Consume(ctx, []string{k.config.TopicName}, consumerHandler)
			if err != nil {
				logger.Error(err.Error())
				return fmt.Errorf("failed to consume kafka: %d", err)
			}

			if err = ctx.Err(); err != nil {
				logger.Error(err.Error())
				return fmt.Errorf("context error on kafka consume: %d", err)
			}

			consumerHandler.Ready = make(chan bool)
		}
	})

	<-consumerHandler.Ready

	logger.Info("Consumer kafka started!")

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
