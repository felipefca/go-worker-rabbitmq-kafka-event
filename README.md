# go-worker-rabbitmq-kafka
A Go worker to consume messages from RabbitMQ and Apache Kafka and save them events in Redis/MongoDB.

[![LinkedIn][linkedin-shield]][linkedin-url]


<!-- ABOUT THE PROJECT -->
## About the Project

A Go worker to consume messages with 2 go routines from RabbitMQ and Apache Kafka (with Avro) and save them events in Redis/MongoDB.

### Related Projects
- https://github.com/felipefca/go-api-event

- https://github.com/felipefca/go-job-aws-s3-kafka-event

![Screenshot_3](https://github.com/felipefca/go-api-event/assets/21323326/691c6cfe-2bce-48bc-b437-b60964738db4)

### Using



* [![Go][Go-badge]][Go-url]
* [![Redis](https://img.shields.io/badge/Redis-v6.0-red.svg)](https://redis.io/)
* [![MongoDB](https://img.shields.io/badge/MongoDB-v5.0-green.svg)](https://www.mongodb.com/)
* [![Apache Kafka](https://img.shields.io/badge/Apache%20Kafka-v3.0-orange.svg)](https://kafka.apache.org/)
* [![RabbitMQ][RabbitMQ-badge]][RabbitMQ-url]
* [![Docker][Docker-badge]][Docker-url]

<!-- GETTING STARTED -->
## Getting Started

Instructions for running the application

### Prerequisites

Run the command to initialize RedisDB, MongoDB, RabbitMQ, Kafka (zookeeper, schema registry...) and the application on the selected port
* docker
  ```sh
  docker-compose up -d
  ```

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/felipefca/go-worker-rabbitmq-kafka-event.git
   ```

2. Exec
   ```sh
   go mod tidy
   ```
   
   ```sh
   cp .env.example .env
   ```
      
   ```sh
   go run ./cmd/main.go
   ```


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://www.linkedin.com/in/felipe-fernandes-fca/
[Go-url]: https://golang.org/
[Go-badge]: https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white
[RabbitMQ-badge]: https://img.shields.io/badge/rabbitmq-%23ff6600.svg?style=flat&logo=rabbitmq&logoColor=white
[RabbitMQ-url]: https://www.rabbitmq.com/
[Docker-badge]: https://img.shields.io/badge/docker-%230db7ed.svg?style=flat&logo=docker&logoColor=white
[Docker-url]: https://www.docker.com/
