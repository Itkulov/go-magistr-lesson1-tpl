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
		memTotal, _ := strconv.ParseFloat(values[1], 64)
		memUsed, _ := strconv.ParseFloat(values[2], 64)
		diskTotal, _ := strconv.ParseFloat(values[3], 64)
		diskUsed, _ := strconv.ParseFloat(values[4], 64)
		networkBandwidth, _ := strconv.ParseFloat(values[5], 64)
		networkUsage, _ := strconv.ParseFloat(values[6], 64)

		memUsagePercent := (memUsed / memTotal) * 100
		diskFreeBytes := diskTotal - diskUsed
		diskFreeMB := diskFreeBytes / 1e6
		networkUsagePercent := (networkUsage / networkBandwidth) * 100
		networkFreeMbit := ((networkBandwidth - networkUsage) * 8) / 1e6

		if loadAvg > 30 {
			fmt.Printf("Load Average is too high: %.0f\n", loadAvg)
		}
		if memUsagePercent > 80 {
			fmt.Printf("Memory usage too high: %d%%\n", int(memUsagePercent))
		}
		if diskFreeBytes < diskTotal*0.1 {
			fmt.Printf("Free disk space is too low: %d Mb left\n", int(diskFreeMB))
		}
		if networkUsagePercent > 90 {
			fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", int(networkFreeMbit))
		}
	}
}
