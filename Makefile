run:
	go run ./cmd/app/

build:
	go build -o image-resizer ./cmd/app/

amqp-up:
	docker-compose up rabbitmq

dc-up:
	docker-compose up

dc-down:
	docker-compose down