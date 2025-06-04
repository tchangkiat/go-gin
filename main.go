package main

import (
	"net/http"

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

	r.Run(":8000")
}
