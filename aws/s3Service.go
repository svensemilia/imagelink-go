package aws

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/svensemilia/imagelink-go/constants"
	"github.com/svensemilia/imagelink-go/image"
)

func S3Upload(fileUp *multipart.FileHeader, filename string, userSub string, album string, collector chan int) {
	file, err := fileUp.Open()
	defer file.Close()
	if err != nil {
		fmt.Println("Error occured while opening file", err)
		collector <- 0
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err != nil {
		fmt.Println("Error occured while creating a AWS session", err)
		collector <- 0
		return
	}

	bucket := "imagelink-version-3-upload-bucket"
	var stringBuilder strings.Builder
	stringBuilder.WriteString(userSub)
	stringBuilder.WriteString("/")
	if len(strings.Trim(album, " ")) > 0 {
		stringBuilder.WriteString(album)
		stringBuilder.WriteString("/")
	}
	stringBuilder.WriteString(filename)
	key := stringBuilder.String()

	fmt.Println(key)
	uploader := s3manager.NewUploader(sess)

	content := "image/*"
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: &content,
	})
	if err != nil {
		// Print the error and exit.
		fmt.Println("Unable to upload:", err)
		collector <- 0
		return
	}

	fmt.Printf("Successfully uploaded %q to %q\n", key, bucket)
	collector <- 1
}

func S3Download(objectKey, userSub string) []byte {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err != nil {
		fmt.Println("Error occured while creating a AWS session", err)
		return nil
	}

	downloader := s3manager.NewDownloader(sess)

	bucket := "imagelink-version-3-upload-bucket"
	var stringBuilder strings.Builder
	stringBuilder.WriteString(userSub)
	stringBuilder.WriteString(objectKey)
	key := stringBuilder.String()

	var buffer *aws.WriteAtBuffer
	init := make([]byte, 0, 200000)
	fmt.Println(len(init))
	buffer = aws.NewWriteAtBuffer(init)

	fmt.Println("S3 reading key", key)
	_, err = downloader.Download(buffer,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		fmt.Println("Error occured while reading from Bucket:", err)
		return nil
	}

	return buffer.Bytes()
}

func GetImages(folder, userSub string) ([][]byte, error) {
	maxElem := 10
	bucket := "imagelink-version-3-upload-bucket"
	var bytes = make([][]byte, 0, maxElem)

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	svc := s3.New(sess)
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(int64(maxElem)),
	}

	result, err := svc.ListObjectsV2(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return bytes, err
	}

	downloader := s3manager.NewDownloader(sess)
	var imgBuffer []byte
	var buffer *aws.WriteAtBuffer
	for _, content := range result.Contents {
		imgBuffer = make([]byte, 0, int(*content.Size))
		buffer = aws.NewWriteAtBuffer(imgBuffer)

		_, err = downloader.Download(buffer,
			&s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(*content.Key),
			})
		if err != nil {
			fmt.Println("Error occured while reading from Bucket:", err)
			return bytes, err
		}
		bytes = append(bytes, image.ScaleImage2(buffer.Bytes(), constants.ImageSize))
	}
	return bytes, nil
}
