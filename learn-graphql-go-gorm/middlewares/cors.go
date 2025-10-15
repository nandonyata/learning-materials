package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func CORS(c *gin.Context) {
	host := c.Request.Header.Get("Origin")
	fmt.Printf("CORS middleware : host: %v", host)

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, "+
		"X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Content-Language, auth-token, "+
		"Upload-Length, Access-Control-Allow-Origin, Location, Tus-Resumable, Upload-Metadata, Referer, "+
		"Organization-username")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}
