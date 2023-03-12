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
	DownloadFromDisk(ctx *gin.Context, id, quality string)
	ShowFromDisk(ctx *gin.Context, id, quality string)
}

func NewPhotoController(photoService services.PhotoService, amqpService services.AMQPService) PhotoController {
	return &photoController{photoService, amqpService}
}

type photoController struct {
	photoService services.PhotoService
	amqpService  services.AMQPService
}

// Upload is used for parsing images and sending it to RabbitMQ
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

// ConsumeAndResize is listening to consumer channel and resizing all images from there and saving them.
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

// DownloadFromDisk is used for serving images to download them
func (p *photoController) DownloadFromDisk(ctx *gin.Context, id, quality string) {
	filename, filePath, err := p.photoService.GetFile(id, quality)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Writer.Header().Set("Content-Type", "multipart/form-data")
	ctx.Writer.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment; filename= %s", filename),
	)
	ctx.File(filePath)
}

// ShowFromDisk is used for serving images to view them
func (p *photoController) ShowFromDisk(ctx *gin.Context, id, quality string) {
	_, filePath, err := p.photoService.GetFile(id, quality)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.File(filePath)
}
