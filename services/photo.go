package services

import (
	"fmt"
	"github.com/h2non/bimg"
	_ "github.com/h2non/bimg"
	"log"
	"os"
)

type PhotoService interface {
	SaveFiles(b []byte) (int64, error)
}

func NewPhotoService(path string) PhotoService {
	return &photoService{
		number: 0,
		path:   path,
	}
}

type photoService struct {
	number int64
	path   string
}

func (p *photoService) SaveFiles(b []byte) (int64, error) {
	image := bimg.NewImage(b)
	p.number++

	size, err := image.Size()
	if err != nil {
		return 0, err
	}

	err = p.SaveResized(image, 100, size.Width, size.Height)
	if err != nil {
		return 0, err
	}

	err = p.SaveResized(image, 75, size.Width, size.Height)
	if err != nil {
		return 0, err
	}

	err = p.SaveResized(image, 50, size.Width, size.Height)
	if err != nil {
		return 0, err
	}

	err = p.SaveResized(image, 25, size.Width, size.Height)
	if err != nil {
		return 0, err
	}

	return p.number, nil
}

func (p *photoService) SaveResized(b *bimg.Image, percents, width, height int) error {
	resized, err := b.Resize(
		(width*percents)/100,
		(height*percents)/100,
	)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%s/%d-%d.png", p.path, p.number, percents)

	f, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = f.Close()

	if err != nil {
		log.Fatalln(err.Error())
	}
	return bimg.Write(filename, resized)
}

func (p *photoService) ResizePercentage(b *bimg.Image, percents int) {

}
