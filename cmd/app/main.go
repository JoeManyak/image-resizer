package main

import (
	"github.com/gin-gonic/gin"
	"image-resizer/config"
	"image-resizer/controllers"
	"image-resizer/services"
	"log"
	"net/http"
	"os"
)

func main() {
	config.Setup()

	_ = os.Mkdir(config.MainConfig.ImagePath, os.ModePerm)
	photoService := services.NewPhotoService(config.MainConfig.ImagePath)

	amqpService := services.NewAMQPService(config.MainConfig.AMQPConfig.QueueName)
	err := amqpService.Setup()
	if err != nil {
		log.Fatalln(err.Error())
	}

	photoController := controllers.NewPhotoController(photoService, amqpService)

	go photoController.ConsumeAndResize()

	r := gin.Default()
	r.GET("/lifecheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.POST("/upload", func(c *gin.Context) {
		photoController.Upload(c)
	})

	if err := r.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}
