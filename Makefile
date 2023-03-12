run:
	go run ./cmd/app/main.go

build:
	go build -o image-resizer ./cmd/app/main.go

amqp-up:
	docker-compose up -d rabbitmq

dc-up:
	docker-compose up -d

dc-down:
	docker-compose down