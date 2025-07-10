package main

import (
	"os"

	"go-gin/middleware"
	"go-gin/routes"

	"github.com/gin-contrib/pprof"

	"github.com/aws/aws-xray-sdk-go/v2/awsplugins/ec2"
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
	r.Use(middleware.ErrorHandling)

	routes.Base(r)
	routes.Perf(r)
	routes.BadPerf(r)

	pprof.Register(r)

	r.Run(":8000")
}
