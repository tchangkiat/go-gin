package routes

// import (
// 	"fmt"
// 	"math/rand"
// 	"net/http"
// 	"os"
// 	"runtime"
// 	"strconv"
// 	"strings"
// 	"sync"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/shirou/gopsutil/v4/host"
// )

// // Endpoints for operations that are optimized
// func Perf(router *gin.Engine) {
// 	// Add paths with prefixes. Use case: handle traffic from multiple load balancer paths but routing to the same service
// 	pathPrefixes := []string{""}
// 	if os.Getenv("PATH_PREFIXES") != "" {
// 		pathPrefixes = append(strings.Split(os.Getenv("PATH_PREFIXES"), `,`), pathPrefixes...)
// 	}
// 	for _, pathPrefix := range pathPrefixes {
// 		base := router.Group(pathPrefix + "/perf")
// 		{
// 			base.GET("/mm", matrixMultiplication)
// 		}
// 	}
// }

// // Matrix represents a 2D matrix with contiguous memory layout
// type Matrix struct {
// 	Rows, Cols int
// 	Data       []float64
// }

// // Get returns matrix element at (i,j)
// func (m *Matrix) Get(i, j int) float64 {
// 	return m.Data[i*m.Cols+j]
// }

// // Set sets matrix element at (i,j)
// func (m *Matrix) Set(i, j int, val float64) {
// 	m.Data[i*m.Cols+j] = val
// }

// // NewMatrix creates a new matrix with contiguous memory layout
// func NewMatrix(rows, cols int) *Matrix {
// 	return &Matrix{
// 		Rows: rows,
// 		Cols: cols,
// 		Data: make([]float64, rows*cols),
// 	}
// }

// // FillRandom fills the matrix with random values between 0 and 1
// func (m *Matrix) FillRandom() {
// 	rand.Seed(time.Now().UnixNano())
// 	for i := 0; i < len(m.Data); i++ {
// 		m.Data[i] = rand.Float64()
// 	}
// }

// // SingleThreadMultiply performs matrix multiplication using a single thread with Graviton optimizations
// func SingleThreadMultiply(a, b *Matrix) (*Matrix, error) {
// 	if a.Cols != b.Rows {
// 		return nil, fmt.Errorf("matrix dimensions don't match for multiplication: a(%dx%d) * b(%dx%d)",
// 			a.Rows, a.Cols, b.Rows, b.Cols)
// 	}

// 	result := NewMatrix(a.Rows, b.Cols)

// 	// Block size optimized for Graviton's 64KB L1 cache
// 	blockSize := 64
// 	if a.Rows < 64 {
// 		blockSize = a.Rows
// 	}

// 	// Blocked matrix multiplication for better cache utilization
// 	for ii := 0; ii < a.Rows; ii += blockSize {
// 		for kk := 0; kk < a.Cols; kk += blockSize {
// 			for jj := 0; jj < b.Cols; jj += blockSize {
// 				// Process block
// 				iEnd := ii + blockSize
// 				if iEnd > a.Rows {
// 					iEnd = a.Rows
// 				}
// 				kEnd := kk + blockSize
// 				if kEnd > a.Cols {
// 					kEnd = a.Cols
// 				}
// 				jEnd := jj + blockSize
// 				if jEnd > b.Cols {
// 					jEnd = b.Cols
// 				}

// 				for i := ii; i < iEnd; i++ {
// 					for k := kk; k < kEnd; k++ {
// 						aVal := a.Data[i*a.Cols+k]
// 						for j := jj; j < jEnd; j++ {
// 							result.Data[i*result.Cols+j] += aVal * b.Data[k*b.Cols+j]
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return result, nil
// }

// // MultiThreadMultiply performs matrix multiplication using multiple threads optimized for Graviton
// func MultiThreadMultiply(a, b *Matrix) (*Matrix, error) {
// 	if a.Cols != b.Rows {
// 		return nil, fmt.Errorf("matrix dimensions don't match for multiplication: a(%dx%d) * b(%dx%d)",
// 			a.Rows, a.Cols, b.Rows, b.Cols)
// 	}

// 	numCPU := runtime.NumCPU()
// 	result := NewMatrix(a.Rows, b.Cols)
// 	var wg sync.WaitGroup

// 	// Block size optimized for Graviton's cache hierarchy
// 	blockSize := 64
// 	if a.Rows < 64 {
// 		blockSize = a.Rows
// 	}

// 	// Calculate blocks per thread for better load balancing
// 	totalBlocks := (a.Rows + blockSize - 1) / blockSize
// 	blocksPerThread := (totalBlocks + numCPU - 1) / numCPU

// 	worker := func(threadID int) {
// 		defer wg.Done()
// 		startBlock := threadID * blocksPerThread
// 		endBlock := startBlock + blocksPerThread
// 		if endBlock > totalBlocks {
// 			endBlock = totalBlocks
// 		}

// 		for blockI := startBlock; blockI < endBlock; blockI++ {
// 			ii := blockI * blockSize
// 			iEnd := ii + blockSize
// 			if iEnd > a.Rows {
// 				iEnd = a.Rows
// 			}

// 			for kk := 0; kk < a.Cols; kk += blockSize {
// 				for jj := 0; jj < b.Cols; jj += blockSize {
// 					kEnd := kk + blockSize
// 					if kEnd > a.Cols {
// 						kEnd = a.Cols
// 					}
// 					jEnd := jj + blockSize
// 					if jEnd > b.Cols {
// 						jEnd = b.Cols
// 					}

// 					for i := ii; i < iEnd; i++ {
// 						for k := kk; k < kEnd; k++ {
// 							aVal := a.Data[i*a.Cols+k]
// 							for j := jj; j < jEnd; j++ {
// 								result.Data[i*result.Cols+j] += aVal * b.Data[k*b.Cols+j]
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	for t := 0; t < numCPU; t++ {
// 		wg.Add(1)
// 		go worker(t)
// 	}

// 	wg.Wait()
// 	return result, nil
// }

// func matrixMultiplication(c *gin.Context) {
// 	size_str := c.Query("size")

// 	// Matrix size
// 	size, _ := strconv.Atoi(size_str)

// 	a := NewMatrix(size, size)
// 	b := NewMatrix(size, size)

// 	a.FillRandom()
// 	b.FillRandom()

// 	// Single thread
// 	start := time.Now()
// 	_, err := SingleThreadMultiply(a, b)
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 	}
// 	singleTime := time.Since(start)

// 	// Multi thread
// 	start = time.Now()
// 	_, err = MultiThreadMultiply(a, b)
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 	}
// 	multiTime := time.Since(start)

// 	hostInfo, _ := host.Info()
// 	c.JSON(http.StatusOK, gin.H{"arch": hostInfo.KernelArch, "matrixMultiplication": gin.H{"size": size, "timeTaken": gin.H{"singleThreaded": singleTime.String(), "multiThreaded": multiTime.String()}}})
// }
