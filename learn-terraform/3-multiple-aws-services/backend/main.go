package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/joho/godotenv"
)

var tableName string

// FileRecord is one item in the DynamoDB table. "ID" is the hash_key.
type FileRecord struct {
	ID   string `dynamodbav:"ID" json:"id"`
	Name string `dynamodbav:"Name" json:"name"`
	Path string `dynamodbav:"Path" json:"path"` // e.g. s3://bucket/uploads/<id>/<name>
}

// Server holds the AWS clients so handlers can use them.
type Server struct {
	s3     *s3.Client
	db     *dynamodb.Client
	bucket string
}

func main() {
	_ = godotenv.Load() // ignore the error if there is no .env file

	bucket := os.Getenv("BUCKET_NAME")
	if bucket == "" {
		log.Fatal("please set the BUCKET_NAME environment variable")
	}

	tableName = os.Getenv("TABLE_NAME")
	if tableName == "" {
		log.Fatal("please set the TABLE_NAME environment variable")
	}

	// Loads credentials and region from `aws configure` / environment variables
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	srv := &Server{
		s3: s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.UsePathStyle = os.Getenv("AWS_ENDPOINT_URL") != ""
		}),
		db:     dynamodb.NewFromConfig(cfg),
		bucket: bucket,
	}

	r := gin.Default()
	r.POST("/files", srv.uploadFile) // upload a file
	r.GET("/files", srv.listFiles)   // list all saved records
	r.GET("/files/:id", srv.getFile) // get one record by ID

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(r.Run(":" + port))
}

// POST /files  (multipart form, field name "file")
func (s *Server) uploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form field 'file' is required"})
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open uploaded file"})
		return
	}
	defer src.Close()

	id := uuid.NewString()
	name := filepath.Base(fileHeader.Filename) // strips any folder parts
	key := fmt.Sprintf("uploads/%s/%s", id, name)

	// 1. Upload to S3
	_, err = s.s3.PutObject(c.Request.Context(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        src,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		log.Printf("s3 upload error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload to S3 failed"})
		return
	}

	// 2. Save metadata to DynamoDB
	rec := FileRecord{
		ID:   id,
		Name: name,
		Path: fmt.Sprintf("s3://%s/%s", s.bucket, key),
	}
	item, err := attributevalue.MarshalMap(rec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot build record"})
		return
	}
	_, err = s.db.PutItem(c.Request.Context(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		log.Printf("dynamodb put error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save to DynamoDB failed"})
		return
	}

	c.JSON(http.StatusCreated, rec)
}

// GET /files/:id
func (s *Server) getFile(c *gin.Context) {
	out, err := s.db.GetItem(c.Request.Context(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: c.Param("id")},
		},
	})
	if err != nil {
		log.Printf("dynamodb get error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "read from DynamoDB failed"})
		return
	}
	if out.Item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	var rec FileRecord
	if err := attributevalue.UnmarshalMap(out.Item, &rec); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot read record"})
		return
	}
	c.JSON(http.StatusOK, rec)
}

// GET /files
// Uses Scan, which reads the whole table. Fine for small tables and learning.
func (s *Server) listFiles(c *gin.Context) {
	records := []FileRecord{}

	pager := dynamodb.NewScanPaginator(s.db, &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	for pager.HasMorePages() {
		page, err := pager.NextPage(c.Request.Context())
		if err != nil {
			log.Printf("dynamodb scan error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "list from DynamoDB failed"})
			return
		}

		var batch []FileRecord
		if err := attributevalue.UnmarshalListOfMaps(page.Items, &batch); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot read records"})
			return
		}
		records = append(records, batch...)
	}

	c.JSON(http.StatusOK, records)
}