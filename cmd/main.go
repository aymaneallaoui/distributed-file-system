package main

import (
	"distributed-file-system/config"
	"distributed-file-system/pkg/filesystem"
	"distributed-file-system/pkg/storage"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var fs *filesystem.FileSystem
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	config.LoadConfig()

	nodes := []*storage.Node{
		&storage.Node{ID: "node-1"},
		&storage.Node{ID: "node-2"},
	}

	replicationFactor, _ := strconv.Atoi(os.Getenv("REPLICATION_FACTOR"))
	fs = filesystem.NewFileSystem(nodes, replicationFactor)

	router := gin.Default()
	currWorkingDir, _ := os.Getwd()

	router.Static("/static", filepath.Join(currWorkingDir, "web", "static"))

	templateFilePath := filepath.Join(currWorkingDir, "web", "templates", "*")

	funcMap := template.FuncMap{
		"sanitize": func(name string) string {

			re := regexp.MustCompile(`[^\w-]`)
			return re.ReplaceAllString(name, "-")
		},
	}

	router.SetHTMLTemplate(template.Must(template.New("").Funcs(funcMap).ParseGlob(templateFilePath)))

	router.GET("/", listFiles)
	router.POST("/upload", uploadFile)
	router.GET("/download/:filename", downloadFile)
	router.GET("/list", listFiles)
	router.GET("/progress/:filename", progressHandler)

	router.Run(":8080")
}

func listFiles(c *gin.Context) {
	files := fs.ListFiles()
	c.HTML(http.StatusOK, "index.html", gin.H{
		"files": files,
	})
}

func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	tempPath := "./temp/" + file.Filename
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	outputPath := "./output/" + file.Filename
	err = os.Rename(tempPath, outputPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to move file to output directory"})
		return
	}

	if err := fs.UploadFile(outputPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to distributed FS"})
		return
	}

	c.HTML(http.StatusOK, "upload.html", gin.H{"filename": file.Filename})
}

func downloadFile(c *gin.Context) {
	filename := c.Param("filename")
	fmt.Printf("Requested filename: %s\n", filename)

	sanitizedFilename := sanitize(filename)
	fmt.Printf("Sanitized filename: %s\n", sanitizedFilename)

	fileMeta, exists := fs.GetFileMetadata(sanitizedFilename)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	outputFile := filepath.Join("output", sanitizedFilename)
	fmt.Printf("Constructed output file path: %s\n", outputFile)

	if err := fs.DownloadFile(fileMeta.ShardIDs, outputFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file"})
		return
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found on disk"})
		return
	}

	c.File(outputFile)
}
func progressHandler(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for i := 0; i <= 100; i += 20 {
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Downloading... %d%%", i)))
		time.Sleep(500 * time.Millisecond)
	}

	conn.WriteMessage(websocket.TextMessage, []byte("Download complete"))
}

func sanitize(name string) string {
	re := regexp.MustCompile(`[^\w-]`)
	return re.ReplaceAllString(name, "-")
}
