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
		memPerc := memUsed / memTotal * 100
		diskFreeBytes, _ := strconv.ParseFloat(values[6], 64)
		diskFreeMB := diskFreeBytes / (1024 * 1024)
		netUsage, _ := strconv.ParseFloat(values[1], 64)
		netMbit := netUsage / (1024 * 1024)

		if loadAvg > 30 {
			fmt.Printf("Load Average is too high: %.2f\n", loadAvg)
		}
		if memPerc > 80 {
			fmt.Printf("Memory usage too high: %d%%\n", int(memPerc))
		}
		if diskFreeMB < 32000 {
			fmt.Printf("Free disk space is too low: %d Mb left\n", int(diskFreeMB))
		}
		if netMbit > 150 {
			fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", int(netMbit))
		}
	}
}
