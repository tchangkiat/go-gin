package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

func Base(router *gin.Engine) {
	base := router.Group("/")
	{
		base.GET("/fib", fibonacci)
		base.GET("/", getSysInfo)
	}
}

func getSysInfo(c *gin.Context) {
	cpuInfo, _ := cpu.Info()
	memInfo, _ := mem.VirtualMemory()
	hostInfo, _ := host.Info()

	memFilteredInfo := gin.H{
		"total":       memInfo.Total,
		"available":   memInfo.Available,
		"used":        memInfo.Used,
		"usedPercent": memInfo.UsedPercent,
	}

	cpuFilteredInfo := make([]gin.H, len(cpuInfo))
	for _, cpu := range cpuInfo {
		cpuFilteredInfo = append(cpuFilteredInfo, gin.H{
			"vendorId":  cpu.VendorID,
			"family":    cpu.Family,
			"model":     cpu.Model,
			"modelName": cpu.ModelName,
			"cores":     cpu.Cores,
			"mhz":       cpu.Mhz,
			"cacheSize": cpu.CacheSize,
		})
	}

	hostFilteredInfo := gin.H{
		"hostname":        hostInfo.Hostname,
		"os":              hostInfo.OS,
		"platform":        hostInfo.Platform,
		"platformFamily":  hostInfo.PlatformFamily,
		"platformVersion": hostInfo.PlatformVersion,
		"kernelVersion":   hostInfo.KernelVersion,
		"kernelArch":      hostInfo.KernelArch,
	}

	c.JSON(http.StatusOK, gin.H{
		"cpu":  cpuFilteredInfo,
		"mem":  memFilteredInfo,
		"host": hostFilteredInfo,
	})
}

func fibonacci(c *gin.Context) {
	n_str := c.Query("n")
	n, _ := strconv.Atoi(n_str)
	hostInfo, _ := host.Info()

	start := time.Now()
	fib_result := fib(n)
	elapsed := time.Since(start)

	c.JSON(http.StatusOK, gin.H{
		"architecture": hostInfo.KernelArch,
		"fibonacci":    gin.H{"n": n_str, "result": fib_result},
		"timeTaken":    elapsed.String(),
	})
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
