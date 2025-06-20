package main

import (
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

	if os.Getenv("AWS_XRAY_SDK_DISABLED") == "FALSE" {
		// Conditionally load plugin
		if os.Getenv("GIN_MODE") == "release" {
			ec2.Init()
		}

		xrayDaemonAddr := "127.0.0.1:2100"
		if os.Getenv("AWS_XRAY_DAEMON_ADDRESS") != "" {
			xrayDaemonAddr = os.Getenv("AWS_XRAY_DAEMON_ADDRESS")
		}
		xray.Configure(xray.Config{
			DaemonAddr:     xrayDaemonAddr,
			ServiceVersion: "1.0.0",
		})
	}

	// -----------------------------

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
