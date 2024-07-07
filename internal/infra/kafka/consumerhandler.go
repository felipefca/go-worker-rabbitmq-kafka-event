package kafka

import (
	"context"
	"go-worker-rabbitmq-kafka-event/internal/appctx"

	"github.com/IBM/sarama"
)

type ConsumerHandler struct {
	Ready chan bool
	ctx   context.Context
}

func (c *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	close(c.Ready)
	return nil
}

func (c *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	logger := appctx.FromContext(c.ctx)

	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				logger.Info("messa channel was closed!")
				return nil
			}

			//TODO: FAZER O HANDLER

			logger.Info("kafka message received")
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
