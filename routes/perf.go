package routes

import (
	"fmt"
	"math/rand"
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

func Perf(router *gin.Engine) {
	// Add paths with prefixes. Use case: handle traffic from multiple load balancer paths but routing to the same service
	pathPrefixes := []string{""}
	if os.Getenv("PATH_PREFIXES") != "" {
		pathPrefixes = append(strings.Split(os.Getenv("PATH_PREFIXES"), `,`), pathPrefixes...)
	}
	for _, pathPrefix := range pathPrefixes {
		base := router.Group(pathPrefix + "/perf")
		{
			base.GET("/fib", fibonacci)
			base.GET("/mm", matrixMultiplication)
		}
	}
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
		"arch":      hostInfo.KernelArch,
		"fibonacci": gin.H{"n": n_str, "result": fib_result, "timeTaken": elapsed.String()},
	})
}

// Helper function for Fibonacci sequence
func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

// Matrix represents a 2D matrix
type Matrix struct {
	Rows, Cols int
	Data       [][]float64
}

// NewMatrix creates a new matrix with given dimensions
func NewMatrix(rows, cols int) *Matrix {
	m := &Matrix{
		Rows: rows,
		Cols: cols,
		Data: make([][]float64, rows),
	}
	for i := range m.Data {
		m.Data[i] = make([]float64, cols)
	}
	return m
}

// FillRandom fills the matrix with random values between 0 and 1
func (m *Matrix) FillRandom() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			m.Data[i][j] = rand.Float64()
		}
	}
}

// SingleThreadMultiply performs matrix multiplication using a single thread
func SingleThreadMultiply(a, b *Matrix) (*Matrix, error) {
	if a.Cols != b.Rows {
		return nil, fmt.Errorf("matrix dimensions don't match for multiplication: a(%dx%d) * b(%dx%d)",
			a.Rows, a.Cols, b.Rows, b.Cols)
	}

	result := NewMatrix(a.Rows, b.Cols)

	for i := 0; i < a.Rows; i++ {
		for k := 0; k < a.Cols; k++ {
			for j := 0; j < b.Cols; j++ {
				result.Data[i][j] += a.Data[i][k] * b.Data[k][j]
			}
		}
	}

	return result, nil
}

// MultiThreadMultiply performs matrix multiplication using multiple threads
func MultiThreadMultiply(a, b *Matrix) (*Matrix, error) {
	if a.Cols != b.Rows {
		return nil, fmt.Errorf("matrix dimensions don't match for multiplication: a(%dx%d) * b(%dx%d)",
			a.Rows, a.Cols, b.Rows, b.Cols)
	}

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	result := NewMatrix(a.Rows, b.Cols)
	var wg sync.WaitGroup

	// Calculate how many rows each goroutine should process
	rowsPerThread := a.Rows / numCPU
	if rowsPerThread == 0 {
		// If we have more CPUs than rows, just use one goroutine per row
		rowsPerThread = 1
		numCPU = a.Rows
	}

	// Function to process a subset of rows
	worker := func(startRow, endRow int) {
		defer wg.Done()
		for i := startRow; i < endRow; i++ {
			for j := 0; j < b.Cols; j++ {
				sum := 0.0
				for k := 0; k < a.Cols; k++ {
					sum += a.Data[i][k] * b.Data[k][j]
				}
				result.Data[i][j] = sum
			}
		}
	}

	// Launch workers
	for t := 0; t < numCPU; t++ {
		startRow := t * rowsPerThread
		endRow := startRow + rowsPerThread
		// Handle the case when rows don't divide evenly among threads
		if t == numCPU-1 {
			endRow = a.Rows
		}

		wg.Add(1)
		go worker(startRow, endRow)
	}

	// Wait for all calculations to complete
	wg.Wait()

	return result, nil
}

func matrixMultiplication(c *gin.Context) {
	size_str := c.Query("size")

	// Matrix size
	size, _ := strconv.Atoi(size_str)

	a := NewMatrix(size, size)
	b := NewMatrix(size, size)

	a.FillRandom()
	b.FillRandom()

	// Single thread
	start := time.Now()
	_, err := SingleThreadMultiply(a, b)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	singleTime := time.Since(start)

	// Multi thread
	start = time.Now()
	_, err = MultiThreadMultiply(a, b)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	multiTime := time.Since(start)

	hostInfo, _ := host.Info()
	c.JSON(http.StatusOK, gin.H{"arch": hostInfo.KernelArch, "matrixMultiplication": gin.H{"size": size, "timeTaken": gin.H{"singleThreaded": singleTime.String(), "multiThreaded": multiTime.String()}}})
}
