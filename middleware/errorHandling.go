package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandling(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		log.Println(err)
	}

	c.AbortWithStatus(http.StatusInternalServerError)
}
