package main

import (
	"go-gin/middleware"
	"go-gin/routes"

	"github.com/gin-contrib/pprof"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(middleware.ErrorHandling)

	routes.Base(r)
	routes.Perf(r)
	routes.BadPerf(r)

	pprof.Register(r)

	r.Run(":8000")
}
