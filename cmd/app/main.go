package main

import (
	"github.com/gin-gonic/gin"
	"image-resizer/controllers"
	"image-resizer/services"
	"log"
	"net/http"
	"os"
)

func main() {
	filepath := os.Getenv("IMG_PATH")
	if filepath == "" {
		filepath = "./img"
	}

	_ = os.Mkdir(filepath, os.ModePerm)
	photoService := services.NewPhotoService(filepath)

	amqpService := services.NewAMQPService("base")
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
