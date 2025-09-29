package routes

import (
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/host"
)

// Sequential matrix multiplication using three nested loops
func MultiplySequential(a, b [][]int) [][]int {
	n := len(a)
	result := make([][]int, n)
	for i := range result {
		result[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return result
}

// Parallel matrix multiplication that distributes rows across CPU cores using goroutines
// Uses sync.WaitGroup and runtime.NumCPU() to distribute work across available CPU cores
func MultiplyParallel(a, b [][]int) [][]int {
	n := len(a)
	result := make([][]int, n)
	for i := range result {
		result[i] = make([]int, n)
	}

	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	rowsPerWorker := n / numWorkers

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			end := start + rowsPerWorker
			if end > n {
				end = n
			}
			for i := start; i < end; i++ {
				for j := 0; j < n; j++ {
					for k := 0; k < n; k++ {
						result[i][j] += a[i][k] * b[k][j]
					}
				}
			}
		}(w * rowsPerWorker)
	}
	wg.Wait()
	return result
}

func MatrixMultiplication(c *gin.Context) {
	size_str := c.Query("size")

	// Matrix size
	size, _ := strconv.Atoi(size_str)

	// Create matrix with given size
	a := make([][]int, size)
	for i := range a {
		a[i] = make([]int, size)
		for j := range a[i] {
			a[i][j] = i + j
		}
	}

	// Single thread
	start := time.Now()
	MultiplySequential(a, a)
	seqTime := time.Since(start)

	// Multi thread
	start = time.Now()
	MultiplyParallel(a, a)
	parallelTime := time.Since(start)

	hostInfo, _ := host.Info()
	c.JSON(http.StatusOK, gin.H{"arch": hostInfo.KernelArch, "matrixMultiplication": gin.H{"size": size, "timeTaken": gin.H{"sequential": seqTime.String(), "parallel": parallelTime.String()}}})
}

func MM(router *gin.Engine) {
	// Add paths with prefixes. Use case: handle traffic from multiple load balancer paths but routing to the same service
	pathPrefixes := []string{""}
	if os.Getenv("PATH_PREFIXES") != "" {
		pathPrefixes = append(strings.Split(os.Getenv("PATH_PREFIXES"), `,`), pathPrefixes...)
	}
	for _, pathPrefix := range pathPrefixes {
		base := router.Group(pathPrefix + "/")
		{
			base.GET("/mm", MatrixMultiplication)
		}
	}
}
