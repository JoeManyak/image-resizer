package services

import (
	"fmt"
	"github.com/h2non/bimg"
	_ "github.com/h2non/bimg"
	"os"
)

type PhotoService interface {
	SaveFilesSequence(b []byte) (int64, error)
	SaveResized(b *bimg.Image, percents, width, height int) error
	ResizePercentage(b *bimg.Image, width, height, percents int) ([]byte, error)
	SaveFile(filename string, image []byte) error
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

func (p *photoService) SaveFilesSequence(b []byte) (int64, error) {
	image := bimg.NewImage(b)
	p.number++

	size, err := image.Size()
	if err != nil {
		return 0, fmt.Errorf("get size files sequence: %w", err)
	}

	err = p.SaveResized(image, 100, size.Width, size.Height)
	if err != nil {
		return 0, fmt.Errorf("resize in sequence 100: %w", err)
	}

	err = p.SaveResized(image, 75, size.Width, size.Height)
	if err != nil {
		return 0, fmt.Errorf("resize in sequence 75: %w", err)
	}

	err = p.SaveResized(image, 50, size.Width, size.Height)
	if err != nil {
		return 0, fmt.Errorf("resize in sequence 50: %w", err)
	}

	err = p.SaveResized(image, 25, size.Width, size.Height)
	if err != nil {
		return 0, fmt.Errorf("resize in sequence 25: %w", err)
	}

	return p.number, nil
}

func (p *photoService) SaveResized(b *bimg.Image, percents, width, height int) error {
	resized, err := p.ResizePercentage(b, width, height, percents)
	if err != nil {
		return fmt.Errorf("save resized resizing: %w", err)
	}

	filename := fmt.Sprintf("%s/%d-%d.png", p.path, p.number, percents)
	err = p.SaveFile(filename, resized)
	if err != nil {
		return fmt.Errorf("save resized saving: %w", err)
	}
	return nil
}

func (p *photoService) ResizePercentage(b *bimg.Image, width, height, percents int) ([]byte, error) {
	image, err := b.Resize(
		(width*percents)/100,
		(height*percents)/100,
	)
	if err != nil {
		return nil, fmt.Errorf("resize percentage: %w", err)
	}

	return image, nil
}

func (p *photoService) SaveFile(filename string, image []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file during saving: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("close file during saving: %w", err)
	}

	err = bimg.Write(filename, image)
	if err != nil {
		return fmt.Errorf("unable to write to file: %w", err)
	}
	return nil
}
