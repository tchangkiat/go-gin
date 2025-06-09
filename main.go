package main

import (
	"log"
	"net/http"

	"go-gin/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(handleError)

	routes.Base(r)

	r.Run(":8000")
}

func handleError(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		log.Println(err)
	}

	c.AbortWithStatus(http.StatusInternalServerError)
}
