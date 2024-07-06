package configs

import (
	"time"

	"github.com/spf13/viper"
)

var cfg *config

type config struct {
	Application Application
	RabbitMQ    RabbitMQ
	MongoDB     MongoDB
	RedisDB     RedisDB
}

type Application struct {
	AppName string
}

type RabbitMQ struct {
	HostName               string
	VirtualHost            string
	Port                   int32
	UserName               string
	Password               string
	DeadLetterExchangeName string
	DeadLetterQueueName    string
	DeadLetterTTL          int32
	ExchangeName           string
	QueueName              string
	ErrorQueueName         string
	RoutingKey             string
	MaxRetry               int32
	PrefetchCount          int32
}

type MongoDB struct {
	ConnectionUrl string
}

type RedisDB struct {
	KeyPrefix     string
	Port          int32
	ConnectRetry  int
	Host          string
	TimeOut       time.Duration
	Password      string
	ConnectionUrl string
}

func GetConfig() config {
	return *cfg
}

func init() {
	viper.SetDefault("PORT", "8080")

	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()
	viper.ReadInConfig()

	cfg = &config{
		Application: Application{
			AppName: "go-worker-rabbitmq-kafka-event",
		},
		RabbitMQ: RabbitMQ{
			HostName:               viper.GetString("RABBITMQ_HOST"),
			VirtualHost:            viper.GetString("RABBITMQ_VHOST"),
			Port:                   viper.GetInt32("RABBITMQ_PORT"),
			UserName:               viper.GetString("RABBITMQ_USER"),
			Password:               viper.GetString("RABBITMQ_PASS"),
			DeadLetterExchangeName: viper.GetString("RABBITMQ_DLQ_EXCHANGE"),
			DeadLetterQueueName:    viper.GetString("RABBITMQ_DLQ_QUEUE"),
			DeadLetterTTL:          viper.GetInt32("RABBITMQ_TTL"),
			ExchangeName:           viper.GetString("RABBITMQ_EXCHANGE"),
			QueueName:              viper.GetString("RABBITMQ_QUEUE"),
			ErrorQueueName:         viper.GetString("RABBITMQ_ERROR_QUEUE"),
			RoutingKey:             viper.GetString("RABBITMQ_ROUTING_KEY"),
			MaxRetry:               viper.GetInt32("RABBITMQ_MAX_RETRY"),
			PrefetchCount:          viper.GetInt32("RABBITMQ_PREFETCH_COUNT"),
		},
		MongoDB: MongoDB{
			ConnectionUrl: viper.GetString("MONGODB_CONNECTION"),
		},
		RedisDB: RedisDB{
			KeyPrefix:    viper.GetString("REDIS_KEYPREFIX"),
			Port:         viper.GetInt32("REDIS_PORT"),
			ConnectRetry: viper.GetInt("REDIS_CONNECTRETRY"),
			Host:         viper.GetString("REDIS_HOST"),
			TimeOut:      viper.GetDuration("REDIS_TIMEOUT"),
			Password:     viper.GetString("REDIS_PASSWORD"),
		},
	}
}
