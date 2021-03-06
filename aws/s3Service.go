package aws

import (
	"fmt"
	"mime/multipart"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/svensemilia/imagelink-go/constants"
	"github.com/svensemilia/imagelink-go/helper"
	"github.com/svensemilia/imagelink-go/image"
)

func S3Upload(fileUp *multipart.FileHeader, filename string, userSub string, album string, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := fileUp.Open()
	defer file.Close()
	if err != nil {
		fmt.Println("Error occured while opening file", err)
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err != nil {
		fmt.Println("Error occured while creating a AWS session", err)
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
		fmt.Println("Unable to upload:", err)
	} else {
		fmt.Printf("Successfully uploaded %q to %q\n", key, bucket)
	}
}

type ImageList struct {
	Images            []ImageData
	Dirs              []string
	ContinuationToken string
}

type ImageData struct {
	Name    string
	Content string
	Data    []byte
}

func GetImage(album, key, userSub string) (*ImageData, error) {
	imageData := ImageData{}

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	downloader := s3manager.NewDownloader(sess)
	imgBuffer := make([]byte, 0, 1)
	buffer := aws.NewWriteAtBuffer(imgBuffer)
	completeKey := helper.BuildObjectPathWithKey(userSub, album, key)

	bytes, contentType, err := download(completeKey, buffer, downloader)

	if err != nil {
		return &imageData, err
	}

	imageData.Content = contentType
	imageData.Name = key
	imageData.Data = bytes
	return &imageData, nil
}

func GetImages(album, continuation, userSub string, resolution int) (*ImageList, error) {
	bucket := "imagelink-version-3-upload-bucket"
	var bytes = make([]ImageData, 0, constants.MaxImageRequest)
	var imageList ImageList

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	svc := s3.New(sess)
	var token *string
	if continuation != "" {
		token = &continuation
	}
	deli := "/"
	input := &s3.ListObjectsV2Input{
		Bucket:            aws.String(bucket),
		MaxKeys:           aws.Int64(int64(constants.MaxImageRequest)),
		ContinuationToken: token,
		Prefix:            helper.BuildObjectPath(userSub, album),
		Delimiter:         &deli,
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
		return &imageList, err
	}

	downloader := s3manager.NewDownloader(sess)
	var imgBuffer []byte
	var buffer *aws.WriteAtBuffer

	for _, content := range result.Contents {
		if *content.Size == 0 {
			continue
		}
		imgBuffer = make([]byte, 0, int(*content.Size))
		buffer = aws.NewWriteAtBuffer(imgBuffer)

		imageBytes, contentType, err := download(*content.Key, buffer, downloader)

		if err != nil {
			fmt.Println("Error occured while reading from Bucket:", err)
			return &imageList, err
		}
		// use ffmpeg for video coding
		bytes = append(bytes, ImageData{Name: *content.Key, Content: contentType, Data: image.ScaleImage(imageBytes, resolution)})
	}
	fmt.Println("Con Token: ", result.NextContinuationToken)
	if result.NextContinuationToken != nil {
		imageList.ContinuationToken = *result.NextContinuationToken
	}
	imageList.Images = bytes
	dirs, err := GetSubDirs(album, userSub)
	if err == nil {
		imageList.Dirs = dirs
	}
	return &imageList, err
}

func GetSubDirs(album, userSub string) ([]string, error) {
	var dirNames = make([]string, 0, constants.MaxImageRequest)

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	svc := s3.New(sess)
	deli := "/"

	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(constants.Bucket),
		Delimiter: &deli,
		MaxKeys:   aws.Int64(int64(constants.MaxImageRequest)),
		Prefix:    helper.BuildObjectPath(userSub, album),
	}

	var trimmedDir string
	err := svc.ListObjectsV2Pages(input,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, content := range page.CommonPrefixes {
				trimmedDir = *content.Prefix
				if strings.HasSuffix(trimmedDir, "/") {
					trimmedDir = strings.TrimRight(trimmedDir, "/")
				}
				trimmedDir = trimmedDir[strings.LastIndex(trimmedDir, "/")+1 : len(trimmedDir)]
				dirNames = append(dirNames, trimmedDir)

			}
			return true
		})
	fmt.Println("Dirnames ", dirNames)
	return dirNames, err
}

func download(key string, buffer *aws.WriteAtBuffer, downloader *s3manager.Downloader) ([]byte, string, error) {

	var contentType string
	_, err := downloader.Download(buffer,
		&s3.GetObjectInput{
			Bucket: aws.String(constants.Bucket),
			Key:    aws.String(key),
		})

	if err != nil {
		return nil, "", err
	}

	contentType = image.GetContentType(buffer.Bytes())
	return buffer.Bytes(), contentType, nil
}
