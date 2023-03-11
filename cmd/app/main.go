package main

import (
	"github.com/gin-gonic/gin"
	"image-resizer/controllers"
	"image-resizer/services"
	"net/http"
	"os"
)

func main() {
	// todo change to env path
	os.Mkdir("./img", os.ModePerm)
	photoService := services.NewPhotoService("./img")
	photoController := controllers.NewPhotoController(photoService)

	r := gin.Default()
	r.GET("/lifecheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.POST("/upload", func(c *gin.Context) {
		photoController.Upload(c)
	})

	r.Run()
}
