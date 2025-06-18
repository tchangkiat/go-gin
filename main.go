package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"go-gin/routes"

	"github.com/aws/aws-xray-sdk-go/v2/awsplugins/ec2"
	"github.com/aws/aws-xray-sdk-go/v2/xray"
	"github.com/gin-gonic/gin"
)

func main() {
	// -----------------------------
	// AWS X-Ray
	// -----------------------------

	// conditionally load plugin
	if os.Getenv("ENVIRONMENT") == "production" {
		ec2.Init()
	}

	xray.Configure(xray.Config{
		ServiceVersion: "1.2.3",
	})

	// -----------------------------

	r := gin.Default()
	r.Use(handleTracingAndError)

	routes.Base(r)

	r.Run(":8000")
}

func handleTracingAndError(c *gin.Context) {
	// Create a segment for tracing in AWS X-Ray
	_, seg := xray.BeginSegment(context.Background(), "go-gin")
	c.Next()
	// Close the segment after processing the request
	seg.Close(nil)

	for _, err := range c.Errors {
		log.Println(err)
	}

	c.AbortWithStatus(http.StatusInternalServerError)
}
