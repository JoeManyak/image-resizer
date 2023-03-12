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

	// Should try to create directory for images if it is not exists
	_ = os.Mkdir(config.MainConfig.ImagePath, os.ModePerm)

	photoService, err := services.NewPhotoService(config.MainConfig.ImagePath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	amqpService := services.NewAMQPService(config.MainConfig.AMQPConfig.QueueName)
	err = amqpService.Setup()
	if err != nil {
		log.Fatalln(err.Error())
	}

	photoController := controllers.NewPhotoController(photoService, amqpService)

	go photoController.ConsumeAndResize()

	r := gin.Default()
	r.GET("/lifecheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/upload", func(ctx *gin.Context) {
		photoController.Upload(ctx)
	})

	r.GET("/download/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		quality := ctx.DefaultQuery("quality", "100")

		photoController.DownloadFromDisk(ctx, id, quality)
	})

	r.GET("/image/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		quality := ctx.DefaultQuery("quality", "100")

		photoController.ShowFromDisk(ctx, id, quality)
	})

	if err := r.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}
