package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/svensemilia/imagelink-go/aws"
	"github.com/svensemilia/imagelink-go/image"
)

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "Everything seems to be fine! Nice work"})
}

func androidUpload(c *gin.Context) {

	jwt := c.Request.Header.Get("Authorization")
	userSub := aws.ExtractSub(jwt)
	album := "android"
	fmt.Println("User sub", userSub)

	fmt.Println("Files Upload Endpoint Hit")

	m, _ := c.MultipartForm()

	album = m.Value["album"][0]

	count := 0
	collector := make(chan int, 10)

	for key, value := range m.File {
		count += len(value)
		fmt.Println("Partname", key)
		for i := range value {
			fileHeader := value[i]
			fmt.Println("Filename", fileHeader.Filename)
			//for each fileheader, get a handle to the actual file
			go aws.S3Upload(fileHeader, fileHeader.Filename, userSub, album, collector)
		}
	}

	counter := 0
	for range collector {
		fmt.Println("Upload completed")
		counter++
		if counter == count {
			close(collector)
		}
	}

	c.JSON(200, gin.H{"msg": "Successfully Uploaded Files"})
}

func uploadFiles(c *gin.Context) {

	jwt := c.Request.Header.Get("Authorization")
	userSub := aws.ExtractSub(jwt)
	var album string
	album = "backup"
	fmt.Println("User sub", userSub)

	m, _ := c.MultipartForm()

	//get the *fileheaders
	files := m.File["myFiles"]
	count := len(files)
	collector := make(chan int, count)

	for i := range files {
		//for each fileheader, get a handle to the actual file
		fileHeader := files[i]

		go aws.S3Upload(fileHeader, files[i].Filename, userSub, album, collector)
	}

	counter := 0
	for range collector {
		fmt.Println("Upload completed")
		counter++
		if counter == count {
			close(collector)
		}
	}

	c.JSON(200, gin.H{"msg": "Successfully Uploaded Files"})
}

func imageDownloadScaled(c *gin.Context) {

	resourceId := "/6823214-large.jpg"
	fmt.Println(resourceId)
	if len(resourceId) == 0 {
		c.JSON(500, gin.H{
			"error": "No specific resource requested!",
		})
		return
	} else {
		fmt.Println("The requested resource has the Id", resourceId)
	}

	jwt := c.Request.Header.Get("Authorization")
	userSub := aws.ExtractSub(jwt)
	byteArray := aws.S3Download(resourceId, userSub)
	fmt.Println("Laenge des ByteArrays", len(byteArray))
	if byteArray == nil {
		c.JSON(500, gin.H{
			"error": "Reading from S3 failed",
		})
		return
	}

	image.ScaleImage(byteArray, 256)

}

func imageDownload(c *gin.Context) {

	resourceId := "test" //aws.GetObjectKey(r, "image")
	fmt.Println(resourceId)
	if len(resourceId) == 0 {
		c.JSON(500, gin.H{
			"error": "No specific resource requested!",
		})
		return
	} else {
		fmt.Println("The requested resource has the Id", resourceId)
	}

	jwt := c.Request.Header.Get("Authorization")
	userSub := aws.ExtractSub(jwt)
	byteArray := aws.S3Download(resourceId, userSub)
	fmt.Println("Laenge des ByteArrays", len(byteArray))
	if byteArray == nil {
		c.JSON(500, gin.H{
			"error": "Reading from S3 failed",
		})
		return
	}

	ret := gin.H{}
	ret["image"] = byteArray
	c.JSON(200, ret)

	/*
		w.Header().Set("Content-Type", "image/jpeg")
		_, err := w.Write(byteArray)
		if err != nil {
			fmt.Println("Error?")
			fmt.Fprintf(w, "Writing to Response failed")
			return
		}
	*/
}

func images(c *gin.Context) {
	folder := c.Param("folder")
	jwt := c.Request.Header.Get("Authorization")
	userSub := aws.ExtractSub(jwt)
	images, err := aws.GetImages(folder, userSub)

	ret := gin.H{}
	if err != nil {
		ret["error"] = "Reading from S3 failed"
		c.JSON(500, ret)
		return
	}
	count := 1
	for _, img := range images {
		ret[fmt.Sprintf("img%d", count)] = img
		count++
	}
	c.JSON(200, ret)
}

/*
func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
*/

func main() {
	argsWithoutProg := os.Args[1:]
	fmt.Println(argsWithoutProg)

	router := gin.Default()
	router.GET("/healthcheck", healthCheck)
	router.POST("/androidUpload", androidUpload)
	router.POST("/upload", uploadFiles)
	router.GET("/image", imageDownload)
	router.GET("/images", images)
	router.GET("/images/:folder", images)
	router.GET("/scale", imageDownloadScaled)

	router.Use(cors.Default())
	router.Run(fmt.Sprintf(":%s", "8080"))
}

// export PATH=$PATH:/usr/local/go/bin
// go run .\server.go .\jwtLib.go .\s3Service.go