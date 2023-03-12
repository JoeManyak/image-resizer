# Image resize application

### Commands

| Command          | Description                                                    |
|:-----------------|----------------------------------------------------------------|
| ` make run `     | Start application (Note: you need to run rabbitmq to use that) |
| ` make build `   | Builds application in current directory                        |
| ` make amqp-up ` | Starts rabbitmq service in Docker                              |               
| ` make dc-up `   | Starts all services in Docker                                  |               
| ` make dc-down ` | Stops all services in Docker                                   |              

### Endpoints

| Route             | Description                          |
|:------------------|--------------------------------------|
| ` /lifecheck `    | Returns {"status": "ok"}             |
| ` /upload `       | Is used to upload images             |
| ` /download/{id}` | Is used for downloading images by ID |
| ` /image/{id}`    | Is used for image preview by ID      |

P.S. You can use quality query for download and preview endpoints to set 25%-50%-75%-100% resolution of images
(percent sign is not needed)

### Environmental variables

| Variable       | Description                        | Default       |
|:---------------|------------------------------------|---------------|
| ` IMG_PATH `   | Specifies path for image directory | ` ./img `     |
| ` AMQP_QUEUE ` | Specifies RabbitMQ queue name      | ` main `      |
| ` AMQP_QUEUE ` | Specifies RabbitMQ username        | ` guest `     |
| ` AMQP_QUEUE ` | Specifies RabbitMQ password        | ` guest `     |
| ` AMQP_URL `   | Specifies RabbitMQ URL             | ` localhost ` |
| ` AMQP_PORT `  | Specifies RabbitMQ port            | ` 5672 `      |
