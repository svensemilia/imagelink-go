package image

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"

	"github.com/nfnt/resize"
)

func ScaleImage(imageData []byte, size int) {
	fmt.Println("hello")
	imageDec, _, err := image.Decode(bytes.NewReader(imageData))

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

	out, err := os.Create("test_resized.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, newImage, nil)
}

func ScaleImage2(imageData []byte, size int) []byte {
	fmt.Println("hello")
	imageDec, _, _ := image.Decode(bytes.NewReader(imageData))

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
	jpeg.Encode(buffer, newImage, nil)
	return buffer.Bytes()
}
