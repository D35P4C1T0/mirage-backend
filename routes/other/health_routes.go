package other

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheckResponse represents the structure of the health check response
type HealthCheckResponse struct {
	Status      string        `json:"status"`
	Uptime      string        `json:"uptime"`
	Version     string        `json:"version"`
	Environment string        `json:"environment"`
	System      SystemInfo    `json:"system"`
	Metrics     ServerMetrics `json:"metrics"`
}

// SystemInfo contains information about the system
type SystemInfo struct {
	GoVersion    string `json:"goVersion"`
	NumCPU       int    `json:"numCPU"`
	NumGoRoutine int    `json:"numGoRoutine"`
	MemAlloc     string `json:"memAlloc"`
}

// ServerMetrics contains runtime performance metrics
type ServerMetrics struct {
	Timestamp int64 `json:"timestamp"`
	StartTime int64 `json:"startTime"`
}

var (
	appStartTime time.Time
	appVersion   = "0.1.0"
	appEnv       = "development"
)

func init() {
	appStartTime = time.Now()
}

// SetupHealthRoutes sets up the health check routes
func SetupHealthRoutes(rg *gin.RouterGroup) {
	rg.GET("/health", HealthCheck)
}

// HealthCheck provides a comprehensive health status of the application
func HealthCheck(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	response := HealthCheckResponse{
		Status:      "healthy",
		Uptime:      formatUptime(time.Since(appStartTime)),
		Version:     appVersion,
		Environment: appEnv,
		System: SystemInfo{
			GoVersion:    runtime.Version(),
			NumCPU:       runtime.NumCPU(),
			NumGoRoutine: runtime.NumGoroutine(),
			MemAlloc:     fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024),
		},
		Metrics: ServerMetrics{
			Timestamp: time.Now().Unix(),
			StartTime: appStartTime.Unix(),
		},
	}

	c.JSON(http.StatusOK, response)
}

// formatUptime converts duration to a human-readable uptime string
func formatUptime(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	return fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds",
		days, hours, minutes, seconds)
}
