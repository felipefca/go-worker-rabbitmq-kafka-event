package consumer

import (
	"context"
	"fmt"
	"go-worker-rabbitmq-kafka-event/configs"
	"go-worker-rabbitmq-kafka-event/internal/appctx"
	"go-worker-rabbitmq-kafka-event/internal/middlewares"
	"go-worker-rabbitmq-kafka-event/internal/models"
	"go-worker-rabbitmq-kafka-event/internal/services"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	channel *amqp.Channel
	cfg     configs.RabbitMQ
	service services.Service
}

func NewRabbitMQConsumer(channel *amqp.Channel, cfg configs.RabbitMQ, service services.Service) (*RabbitMQConsumer, error) {
	consumer := RabbitMQConsumer{
		channel: channel,
		cfg:     cfg,
		service: service,
	}

	err := consumer.setup()
	if err != nil {
		return nil, err
	}

	return &consumer, nil
}

func (r *RabbitMQConsumer) Listen(ctx context.Context) error {
	logger := appctx.FromContext(ctx)

	r.channel.Qos(
		int(r.cfg.PrefetchCount),
		0,
		false,
	)

	messages, err := r.channel.Consume(r.cfg.QueueName, configs.GetConfig().Application.AppName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	chClosed := make(chan *amqp.Error, 1)
	r.channel.NotifyClose(chClosed)

	forever := make(chan bool)

	go func() {
		for {
			select {
			case amqpErr := <-chClosed:
				fmt.Println("AMQP Channel closed: %w", amqpErr)
				forever <- false

			case message := <-messages:
				if len(message.Body) == 0 {
					logger.Warn("message is empty!")
					message.Ack(false)
					continue
				}

				correlationId, _ := message.Headers["x-correlation-id"].(string)
				ctx = middlewares.CorrelationIdMiddleware(ctx, correlationId)

				event := models.BuildEvent(string(message.Body), "RabbitMQ")

				err = r.service.HandlerEvent(ctx, *event)
				err = fmt.Errorf("teste")
				if err != nil {
					logger.Error(err.Error())
					r.retryMessage(ctx, message)
					continue
				}

				message.Ack(false)
			}
		}
	}()

	logger.Info(fmt.Sprintf("waiting for messages in exchange %s and queue %s...", r.cfg.ExchangeName, r.cfg.QueueName))
	<-forever
	return nil
}

func (r *RabbitMQConsumer) retryMessage(ctx context.Context, message amqp.Delivery) {
	logger := appctx.FromContext(ctx)

	death, exists := message.Headers["x-death"].([]interface{})
	if exists {
		count := death[0].(amqp.Table)["count"].(int64)

		if count >= int64(r.cfg.MaxRetry) {
			logger.Info("Message has reached the number retries. Sending message to error queue")

			err := r.sendErrorMessage(ctx, message.Body)
			if err != nil {
				logger.Warn(fmt.Sprintf("Error while trying sending message to error queue: %s", err.Error()))
			} else {
				logger.Info("Messa was sent to error queue")
				message.Ack(false)
			}

			return
		}
	}

	logger.Info("Message sended to deadletter queue")
	message.Nack(false, false)
}

func (r *RabbitMQConsumer) sendErrorMessage(ctx context.Context, message []byte) error {
	err := r.channel.PublishWithContext(
		ctx,
		"",
		r.cfg.ErrorQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         message,
			DeliveryMode: amqp.Persistent,
			Headers: amqp.Table{
				"x-correlation-id": appctx.GetCorrelationId(ctx),
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQConsumer) setup() error {
	//Dead Letter Exchange
	err := r.channel.ExchangeDeclare(
		r.cfg.DeadLetterExchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	//Dead Letter Queue
	_, err = r.channel.QueueDeclare(
		r.cfg.DeadLetterQueueName,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    r.cfg.ExchangeName,
			"x-dead-letter-routing-key": r.cfg.RoutingKey,
			"x-message-ttl":             r.cfg.DeadLetterTTL,
		},
	)
	if err != nil {
		return err
	}

	//Bind Dead Letter Queue to Dead Letter Exchange
	r.channel.QueueBind(r.cfg.DeadLetterQueueName, "", r.cfg.DeadLetterExchangeName, false, nil)

	//Exchange
	err = r.channel.ExchangeDeclare(
		r.cfg.ExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	//Queue
	_, err = r.channel.QueueDeclare(
		r.cfg.QueueName,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange": r.cfg.DeadLetterExchangeName,
		},
	)
	if err != nil {
		return err
	}

	//Bind Queue to Exchange
	r.channel.QueueBind(r.cfg.QueueName, r.cfg.RoutingKey, r.cfg.ExchangeName, false, nil)

	//Error Queue
	_, err = r.channel.QueueDeclare(
		r.cfg.ErrorQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}
