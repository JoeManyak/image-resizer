package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image-resizer/services"
	"log"
	"net/http"
)

type PhotoController interface {
	Upload(ctx *gin.Context)
	ConsumeAndResize()
}

func NewPhotoController(photoService services.PhotoService, amqpService services.AMQPService) PhotoController {
	return &photoController{photoService, amqpService}
}

type photoController struct {
	photoService services.PhotoService
	amqpService  services.AMQPService
}

func (p *photoController) Upload(ctx *gin.Context) {
	raw, err := ctx.GetRawData()
	if err != nil {
		ctx.Status(http.StatusUnprocessableEntity)
		log.Println(fmt.Errorf("upload: %w", err))
		return
	}

	err = p.amqpService.Send(ctx, raw)
	if err != nil {
		ctx.Status(http.StatusUnprocessableEntity)
		log.Println(fmt.Errorf("upload: %w", err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (p *photoController) ConsumeAndResize() {
	consumer, err := p.amqpService.GetConsumer()
	if err != nil {
		return
	}

	for msg := range consumer {
		num, err := p.photoService.SaveFilesSequence(msg.Body)
		if err != nil {
			log.Printf("Failed to create sequence with ID=%d, error [%s]\n", num, err.Error())
			continue
		}

		log.Printf("Created sequence with ID=%d\n", num)
	}
}
