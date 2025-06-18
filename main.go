package main

import (
	//"context"
	"log"
	"net/http"

	//"os"

	"go-gin/routes"

	//"github.com/aws/aws-xray-sdk-go/v2/awsplugins/ec2"
	//"github.com/aws/aws-xray-sdk-go/v2/xray"
	"github.com/gin-gonic/gin"
)

func main() {
	// -----------------------------
	// AWS X-Ray
	// -----------------------------

	// if os.Getenv("TRACING") == "true" {
	// 	// conditionally load plugin
	// 	if os.Getenv("GIN_MODE") == "release" {
	// 		ec2.Init()
	// 	}

	// 	xrayDaemonAddr := "127.0.0.1:2000"
	// 	if os.Getenv("AWS_XRAY_DAEMON_ADDRESS") != "" {
	// 		xrayDaemonAddr = os.Getenv("AWS_XRAY_DAEMON_ADDRESS")
	// 	}
	// 	xray.Configure(xray.Config{
	// 		DaemonAddr:     xrayDaemonAddr,
	// 		ServiceVersion: "1.2.3",
	// 	})
	// }

	// -----------------------------

	r := gin.Default()
	r.Use(handleTracingAndError)

	routes.Base(r)

	r.Run(":8000")
}

func handleTracingAndError(c *gin.Context) {
	// if os.Getenv("TRACING") == "true" {
	// 	// Create a segment for tracing in AWS X-Ray
	// 	_, seg := xray.BeginSegment(context.Background(), "go-gin")
	// 	c.Next()
	// 	// Close the segment after processing the request
	// 	seg.Close(nil)
	// } else {
	// 	c.Next()
	// }
	c.Next()

	for _, err := range c.Errors {
		log.Println(err)
	}

	c.AbortWithStatus(http.StatusInternalServerError)
}
