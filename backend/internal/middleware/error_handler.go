package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("request error: %v", err)
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}
	}
}
