package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres"
)

type HealthCheckResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	System    SystemInformation `json:"system"`
	Component ComponentInfo     `json:"component"`
}
type SystemInformation struct {
	Version          string `json:"version"`
	GoroutinesCount  int    `json:"goroutines_count"`
	TotalAlloc       string `json:"total_alloc"`
	HeapObjectsCount uint64 `json:"heap_objects_count"`
	Alloc            string `json:"alloc"`
	HeapAlloc        string `json:"heap_alloc"`
}
type ComponentInfo struct {
	ServiceName     string `json:"service_name"`
	DBTotalConns    int32  `json:"db_total_conns"`
	DBIdleConns     int32  `json:"db_idle_conns"`
	DBAcquiredConns int32  `json:"db_acquired_conns"`
}

func HealthCheck(pg *postgres.Postgres, serviceName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		systemInfo := getSystemInformation()
		stat := pg.Pool().Stat()
		response := HealthCheckResponse{
			Status:    http.StatusText(http.StatusOK),
			Timestamp: time.Now().Format(time.RFC3339),
			System:    systemInfo,
			Component: ComponentInfo{
				ServiceName:     serviceName,
				DBTotalConns:    stat.TotalConns(),
				DBIdleConns:     stat.IdleConns(),
				DBAcquiredConns: stat.AcquiredConns(),
			},
		}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to serialize JSON response")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(jsonResponse)
	})
}

func getSystemInformation() SystemInformation {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return SystemInformation{
		Version:          runtime.Version(),
		GoroutinesCount:  runtime.NumGoroutine(),
		TotalAlloc:       formatBytesToMB(memStats.TotalAlloc),
		HeapObjectsCount: memStats.HeapObjects,
		Alloc:            formatBytesToMB(memStats.Alloc),
		HeapAlloc:        formatBytesToMB(memStats.HeapAlloc),
	}
}

func formatBytesToMB(bytes uint64) string {
	const mb = 1024 * 1024
	return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
}
