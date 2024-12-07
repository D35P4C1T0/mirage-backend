package utils

import (
	"bytes"
	"errors"
	"image"
	"log"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"golang.org/x/image/draw"
)

// ScaleAndConvertToWebPBytes scales down an image from a byte slice to 1600x1200 and encodes it to WebP format.
// It allows customization of the WebP quality parameter (0-100).
func ScaleAndConvertToWebPBytes(imageData []byte, quality int) ([]byte, error) {
	if quality < 0 || quality > 100 {
		return nil, errors.New("quality must be between 0 and 100")
	}

	// Decode the input image
	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}

	log.Println("Image format:", format)

	// Get original image dimensions
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Calculate target dimensions while preserving aspect ratio
	targetWidth, targetHeight := 1600, 1200
	ratio := float64(width) / float64(height)
	if ratio > float64(targetWidth)/float64(targetHeight) {
		targetHeight = int(float64(targetWidth) / ratio)
	} else {
		targetWidth = int(float64(targetHeight) * ratio)
	}

	// Scale down the image
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	// Prepare WebP encoding options
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, float32(quality))
	if err != nil {
		log.Fatalf("Failed to create encoder options: %v", err)
	}

	// Encode scaled image to WebP format
	var output bytes.Buffer
	if err := webp.Encode(&output, dst, options); err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

// CompareImageFileSizes compares the sizes of two images.
func CompareImageFileSizes(original []byte, compressed []byte) (float64, error) {
	originalSize := len(original)
	compressedSize := len(compressed)

	if originalSize == 0 {
		return 0, errors.New("original image size is zero")
	}

	return float64(compressedSize) / float64(originalSize), nil
}
