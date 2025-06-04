package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		cpuInfo, _ := cpu.Info()
		memInfo, _ := mem.VirtualMemory()
		hostInfo, _ := host.Info()
		memFilteredInfo := gin.H{
			"total":       memInfo.Total,
			"available":   memInfo.Available,
			"used":        memInfo.Used,
			"usedPercent": memInfo.UsedPercent,
		}

		c.JSON(http.StatusOK, gin.H{
			"cpu":  cpuInfo,
			"mem":  memFilteredInfo,
			"host": hostInfo,
		})
	})

	r.GET("/fib", func(c *gin.Context) {
		n_str := c.Query("n")
		n, _ := strconv.Atoi(n_str)

		start := time.Now()
		answer := fibonacci(n)
		elapsed := time.Since(start)

		c.JSON(http.StatusOK, gin.H{
			"fibonacci-" + n_str: answer,
			"timeTaken":          elapsed.String(),
		})
	})

	r.Run(":8000")
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}
