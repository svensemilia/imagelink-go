package image

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"

	"github.com/nfnt/resize"
)

func ScaleImage(imageData []byte, size int) []byte {
	imageDec, format, _ := image.Decode(bytes.NewReader(imageData))

	rect := imageDec.Bounds()
	width := rect.Max.X
	height := rect.Max.Y
	fmt.Println("width, height", width, height)

	var newImage image.Image
	if width >= height {
		newImage = resize.Resize(0, uint(size), imageDec, resize.NearestNeighbor)
		width = newImage.Bounds().Max.X
		sub := width - size
		subLR := sub / 2

		newImage = newImage.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(image.Rect(subLR, 0, width-subLR, size))

	} else {
		newImage = resize.Resize(uint(size), 0, imageDec, resize.NearestNeighbor)
		height = newImage.Bounds().Max.Y
		sub := height - size
		subLR := sub / 2

		newImage = newImage.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(image.Rect(0, subLR, size, height-subLR))
	}

	var buffer *bytes.Buffer
	imgBuffer := make([]byte, 0, 200000)
	buffer = bytes.NewBuffer(imgBuffer)

	// write new image to file
	if format == "jpeg" {
		jpeg.Encode(buffer, newImage, nil)
	} else if format == "png" {
		png.Encode(buffer, newImage)
	} else if format == "gif" {
		gif.Encode(buffer, newImage, nil)
	} else {
		fmt.Println("Unknown format:", format)
	}
	return buffer.Bytes()
}

func GetContentType(data []byte) string {
	return http.DetectContentType(data)
}
