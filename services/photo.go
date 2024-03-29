package services

import (
	"fmt"
	"github.com/h2non/bimg"
	"image-resizer/config"
	"os"
	"strconv"
	"strings"
)

type PhotoService interface {
	SaveFilesSequence(b []byte) (int, error)
	SaveResized(b *bimg.Image, percents, width, height int) error
	ResizePercentage(b *bimg.Image, width, height, percents int) ([]byte, error)
	SaveFile(filename string, image []byte) error
	GetFile(id, quality string) (string, string, error)
}

func NewPhotoService(path string) (PhotoService, error) {
	initialNumber, err := getInitialNumber()
	if err != nil {
		return nil, fmt.Errorf("couldn't initiate photo service: %w", err)
	}

	return &photoService{
		number: initialNumber,
		path:   path,
	}, nil
}

type photoService struct {
	number int
	path   string
}

// SaveFilesSequence creates 100%, 75%, 50% and 25% quality images into image directory
func (p *photoService) SaveFilesSequence(b []byte) (int, error) {
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

// SaveResized Resizing image and saving it into image directory
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

// ResizePercentage resizing bimg.Image with provided image width, height and percents to resize
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

// SaveFile creates empty file by filename and write there provided bytes
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

// GetFile seeks file from image directory by ID and quality
func (p *photoService) GetFile(id, quality string) (string, string, error) {
	filename := fmt.Sprintf("%s-%s.png", id, quality)
	filePath := fmt.Sprintf("%s/%s", config.MainConfig.ImagePath, filename)

	dir, err := os.ReadDir(config.MainConfig.ImagePath)
	if err != nil {
		return "", "", err
	}

	fileFound := false
	for i := range dir {
		if dir[i].Name() == filename {
			fileFound = true
		}
	}

	if !fileFound {
		return "", "", err
	}

	return filename, filePath, nil
}

// getInitialNumber looks for highest ID in image directory, to avoid collisions. If directory is empty, returns 0
func getInitialNumber() (int, error) {
	dir, err := os.ReadDir(config.MainConfig.ImagePath)
	if err != nil {
		return 0, fmt.Errorf("unable to read dir: %w", err)
	}

	maxNum := 0
	for i := range dir {
		num, err := strconv.Atoi(strings.Split(dir[i].Name(), "-")[0])
		if err != nil {
			continue
		}

		if num > maxNum {
			maxNum = num
		}
	}

	return maxNum, nil
}
