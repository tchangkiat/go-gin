package main

import (
	"log"
	"net/http"

	"os"

	"go-gin/routes"

	"github.com/aws/aws-xray-sdk-go/v2/awsplugins/ec2"
	"github.com/aws/aws-xray-sdk-go/v2/header"
	"github.com/aws/aws-xray-sdk-go/v2/xray"
	"github.com/gin-gonic/gin"
)

func main() {
	// -----------------------------
	// AWS X-Ray Configuration
	// -----------------------------

	if os.Getenv("AWS_XRAY_SDK_DISABLED") == "FALSE" {
		// Conditionally load plugin
		if os.Getenv("GIN_MODE") == "release" {
			ec2.Init()
		}
		xray.Configure(xray.Config{
			ServiceVersion: "1.0.0",
		})
	}

	// -----------------------------

	r := gin.Default()
	r.Use(handleTracingAndError)

	routes.Base(r)

	r.Run(":8000")
}

func handleTracingAndError(c *gin.Context) {
	if os.Getenv("AWS_XRAY_SDK_DISABLED") == "FALSE" {
		// Create a segment for tracing in AWS X-Ray
		traceHeader := header.FromString(c.Request.Header.Get("x-amzn-trace-id"))
		xrayCtx, seg := xray.NewSegmentFromHeader(c.Request.Context(), "web-app", c.Request, traceHeader)
		c.Request = c.Request.WithContext(xrayCtx)
		c.Next()
		// Close the segment after processing the request
		seg.Close(nil)
	} else {
		c.Next()
	}

	for _, err := range c.Errors {
		log.Println(err)
	}

	c.AbortWithStatus(http.StatusInternalServerError)
}
