package routes

import (
	"bufio"
	"encoding/json"
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
		base.GET("/req", proxy_request)
		base.GET("/", getSysInfo)
	}
}

// Get system information
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

	cpuFilteredInfo := make([]gin.H, 0, len(cpuInfo))
	for _, cpuStat := range cpuInfo {
		// Append CPU info only if any of the fields are non-empty
		if cpuStat.VendorID != "" || cpuStat.ModelName != "" {
			cpuFilteredInfo = append(cpuFilteredInfo, gin.H{
				"vendorId": cpuStat.VendorID,
				"model":    cpuStat.ModelName,
			})
		}
	}

	hostFilteredInfo := gin.H{
		"hostname":        hostInfo.Hostname,
		"os":              hostInfo.OS,
		"platform":        hostInfo.Platform,
		"platformVersion": hostInfo.PlatformVersion,
		"kernelVersion":   hostInfo.KernelVersion,
		"kernelArch":      hostInfo.KernelArch,
	}

	c.JSON(http.StatusOK, gin.H{
		"cpu":    cpuFilteredInfo,
		"memory": memFilteredInfo,
		"host":   hostFilteredInfo,
	})
}

// Fibonacci sequence without memoization
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

// Helper function for Fibonacci sequence
func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

// Proxy request to a URL
func proxy_request(c *gin.Context) {
	protocol := c.Query("protocol")
	if protocol == "" {
		protocol = "http"
	}
	host := c.Query("host")
	if host == "" {
		host = "localhost"
	}
	port := c.Query("port")
	if port == "" {
		port = "8000"
	}
	path := c.Query("path")
	if path == "" {
		path = "/"
	}
	url := protocol + "://" + hostname + ":" + port + path
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	var jsonResp map[string]interface{}

	for i := 0; scanner.Scan() && i < 5; i++ {
		json.Unmarshal([]byte(scanner.Text()), &jsonResp)
	}
	c.JSON(http.StatusOK, gin.H{"url": url, "response": jsonResp})
}
