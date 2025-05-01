package image_shoshonify

import (
	"errors"
	"image"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

// LoadImage loads an image from file
func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

// LoadStickers loads all images from a directory
func LoadStickers(dir string) ([]image.Image, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var stickers []image.Image
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		path := filepath.Join(dir, file.Name())
		img, err := LoadImage(path)
		if err != nil {
			continue // Skip invalid images
		}
		stickers = append(stickers, img)
	}

	if len(stickers) == 0 {
		return nil, errors.New("no valid images found in directory")
	}
	return stickers, nil
}

// OverlayStickers overlays random stickers on a base image
func OverlayStickers(
	base image.Image,
	stickers []image.Image,
	count int,
	scale float64,
) (image.Image, error) {
	if count <= 0 {
		return nil, errors.New("count must be greater than 0")
	}
	if scale <= 0 {
		return nil, errors.New("scale must be greater than 0")
	}

	rand.Seed(time.Now().UnixNano())
	ctx := gg.NewContextForImage(base)

	baseBounds := base.Bounds()
	baseWidth := baseBounds.Dx()
	baseHeight := baseBounds.Dy()

	for i := 0; i < count; i++ {
		// Select random sticker
		sticker := stickers[rand.Intn(len(stickers))]

		// Calculate scaled dimensions
		originalWidth := sticker.Bounds().Dx()
		originalHeight := sticker.Bounds().Dy()
		scaledWidth := int(float64(originalWidth) * scale)
		scaledHeight := int(float64(originalHeight) * scale)

		// Resize sticker
		resized := imaging.Resize(sticker, scaledWidth, scaledHeight, imaging.Lanczos)

		// Calculate random position
		maxX := baseWidth - scaledWidth
		if maxX < 0 {
			maxX = 0
		}
		maxY := baseHeight - scaledHeight
		if maxY < 0 {
			maxY = 0
		}

		x := rand.Intn(maxX + 1)
		y := rand.Intn(maxY + 1)

		// Draw sticker
		ctx.DrawImage(resized, x, y)
	}

	return ctx.Image(), nil
}

// SaveImage saves image to file
func SaveImage(img image.Image, path string) error {
	return gg.SavePNG(path, img)
}
