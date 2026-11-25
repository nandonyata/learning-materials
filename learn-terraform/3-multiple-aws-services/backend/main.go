package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
    AWS_ACCESS_KEY_ID     string = ""
    AWS_SECRET_ACCESS_KEY string = ""
    AWS_REGION            string = "us-east-1"
)


type Files struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(255)"`
	Path string `gorm:"type:varchar(255)"`
}

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
	log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&Files{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	router.POST("/upload", func(c *gin.Context) {
	// single file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(file.Filename)

	// Upload the file to specific dst.
	dst := filepath.Join("./files/", filepath.Base(file.Filename))
	c.SaveUploadedFile(file, dst)

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	router.Run(":8080")
}