package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	client := &http.Client{Timeout: 5 * time.Second}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		resp, err := client.Get("http://srv.msk01.gigacorp.local/_stats")
		if err != nil {
			fmt.Printf("Error fetching stats: %v\n", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			continue
		}

		values := strings.Split(strings.TrimSpace(string(body)), ",")
		if len(values) < 7 {
			fmt.Printf("Invalid response format\n")
			continue
		}

		loadAvg, _ := strconv.ParseFloat(values[0], 64)
		memTotal, _ := strconv.ParseFloat(values[2], 64)
		memUsed, _ := strconv.ParseFloat(values[4], 64)
		diskFree, _ := strconv.ParseFloat(values[5], 64)

		if loadAvg > 30 {
			fmt.Printf("Load Average is too high: %.2f\n", loadAvg)
		}

		if memUsed > memTotal*0.9 {
			fmt.Printf("Memory usage is too high: %.2f%%\n", (memUsed/memTotal)*100)
		}

		if diskFree < 1e9 {
			fmt.Printf("Disk space is low: %.2f GB available\n", diskFree/(1024*1024*1024))
		}
	}
}
