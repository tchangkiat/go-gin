package routes

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go-gin/middleware"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-xray-sdk-go/v2/xray"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"golang.org/x/net/context/ctxhttp"
)

func Base(router *gin.Engine) {
	app_name := "web-app"
	if os.Getenv("AWS_XRAY_APP_NAME") != "" {
		app_name = os.Getenv("AWS_XRAY_APP_NAME")
	}
	// Add paths with prefixes. Use case: handle traffic from multiple load balancer paths but routing to the same service
	pathPrefixes := []string{"/"}
	if os.Getenv("PATH_PREFIXES") != "" {
		pathPrefixes = append(strings.Split(os.Getenv("PATH_PREFIXES"), `,`), pathPrefixes...)
	}
	for _, pathPrefix := range pathPrefixes {
		base := router.Group(pathPrefix, middleware.Trace(xray.NewFixedSegmentNamer(app_name)))
		{
			base.GET("/req", proxyRequest)
			base.GET("/", getSysInfo)
		}
	}
}

// Get system information
func getSysInfo(c *gin.Context) {
	hostInfo, _ := host.Info()
	cpuInfo, _ := cpu.Info()
	memInfo, _ := mem.VirtualMemory()
	netInterfaceInfo, _ := net.Interfaces()

	hostFilteredInfo := gin.H{
		"hostname":        hostInfo.Hostname,
		"os":              hostInfo.OS,
		"platform":        hostInfo.Platform,
		"platformVersion": hostInfo.PlatformVersion,
		"kernelVersion":   hostInfo.KernelVersion,
		"kernelArch":      hostInfo.KernelArch,
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

	memFilteredInfo := gin.H{
		"totalInGb":   fmt.Sprintf("%.2f", float64(memInfo.Total)/1000000000),
		"usedInGb":    fmt.Sprintf("%.2f", float64(memInfo.Used)/1000000000),
		"usedPercent": fmt.Sprintf("%.2f", memInfo.UsedPercent),
	}

	netInterfaceFilteredInfo := make([]gin.H, 0, len(netInterfaceInfo))
	for _, netInterface := range netInterfaceInfo {
		// Append network interface info only if there are IPv4 / IPv6 addresses and the IP addresses are not internal host loopback address ranges
		if len(netInterface.Addrs) > 0 && netInterface.Addrs[0].Addr != "127.0.0.1/8" {
			netInterfaceFilteredInfo = append(netInterfaceFilteredInfo, gin.H{
				"name":  netInterface.Name,
				"addrs": netInterface.Addrs,
				"mtu":   netInterface.MTU,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"host":    hostFilteredInfo,
		"cpu":     cpuFilteredInfo,
		"memory":  memFilteredInfo,
		"network": netInterfaceFilteredInfo,
	})
}

// Proxy request to a URL
func proxyRequest(c *gin.Context) {
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
	url := protocol + "://" + host + ":" + port + path
	resp := &http.Response{}
	if os.Getenv("AWS_XRAY_SDK_DISABLED") == "FALSE" {
		// -----------------------------
		// AWS X-Ray
		// -----------------------------
		resp, _ = ctxhttp.Get(c.Request.Context(), xray.Client(nil), url)
	} else {
		resp, _ = http.Get(url)
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	var jsonResp map[string]interface{}

	for i := 0; scanner.Scan() && i < 5; i++ {
		json.Unmarshal([]byte(scanner.Text()), &jsonResp)
	}
	c.JSON(http.StatusOK, gin.H{"url": url, "response": jsonResp})
}
