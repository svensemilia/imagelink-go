package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/svensemilia/imagelink-go/aws"
)

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "Everything seems to be fine! Nice work"})
}

func upload(c *gin.Context) {
	jwt := c.Request.Header.Get("Authorization")
	userSub, err := aws.AuthenticateUser(jwt)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	
	m, _ := c.MultipartForm()
	
	album := ""
	if len(m.Value["album"]) > 0 {
		album = m.Value["album"][0]
	}

	var wg sync.WaitGroup
	
	for key, value := range m.File {
		fmt.Println("Partname", key)
		for i := range value {
			fileHeader := value[i]
			fmt.Println("Filename", fileHeader.Filename)
			//for each fileheader, get a handle to the actual file
			wg.Add(1)
			go aws.S3Upload(fileHeader, fileHeader.Filename, userSub, album, &wg)
		}
	}

	wg.Wait()

	c.JSON(200, gin.H{"msg": "Successfully Uploaded Files"})
}

func image(c *gin.Context) {
	album := c.DefaultQuery("album", "")
	imageKey := c.Query("key")

	if len(imageKey) == 0 {
		c.JSON(500, gin.H{"error": "No specific resource requested!"})
		return
	}

	jwt := c.Request.Header.Get("Authorization")
	userSub, err := aws.AuthenticateUser(jwt)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	image, err := aws.GetImage(album, imageKey, userSub)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, *image)
}

func images(c *gin.Context) {
	resolution := c.Query("resolution")
	album := c.DefaultQuery("album", "")
	continueToken := c.DefaultQuery("continue", "")
	fmt.Println(resolution, album, continueToken)
	jwt := c.Request.Header.Get("Authorization")
	userSub, err := aws.AuthenticateUser(jwt)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	resolutionInt, err := strconv.Atoi(resolution)
	if err != nil {
		c.JSON(500, gin.H{"error": "resolution param must be an int"})
		return
	}

	//aws.GetSubDirs(album, userSub)
	images, err := aws.GetImages(album, continueToken, userSub, resolutionInt)

	if err != nil {
		c.JSON(500, gin.H{"error": "Reading from S3 failed"})
		return
	}

	c.JSON(200, *images)
}

func main() {
	argsWithoutProg := os.Args[1:]
	fmt.Println(argsWithoutProg)

	router := gin.Default()
	router.Use(CORS())
	router.GET("/healthcheck", healthCheck)
	router.POST("/upload", upload)
	router.GET("/image", image)
	router.GET("/images", images)
	router.Run(fmt.Sprintf(":%s", "8080"))
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// export PATH=$PATH:/usr/local/go/bin
// /home/ec2-user/go/bin/imagelink-go
// ps -A (lists all processes)
// ps -A | grep imagelink-go (filter for imagelink)
// kill PID
// sudo kill $(pgrep imagelink-go)
// go run .\server.go .\jwtLib.go .\s3Service.go
