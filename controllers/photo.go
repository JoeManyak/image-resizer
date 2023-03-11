package controllers

import (
	"github.com/gin-gonic/gin"
	"image-resizer/services"
	"log"
	"net/http"
)

type PhotoController interface {
	Upload(ctx *gin.Context)
}

func NewPhotoController(service services.PhotoService) PhotoController {
	return &photoController{service}
}

type photoController struct {
	photoService services.PhotoService
}

func (p *photoController) Upload(ctx *gin.Context) {
	raw, err := ctx.GetRawData()
	if err != nil {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}

	num, err := p.photoService.SaveFiles(raw)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		log.Fatalln(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": num,
	})
}
