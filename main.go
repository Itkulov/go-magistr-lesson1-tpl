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
		netBytes, _ := strconv.ParseFloat(values[1], 64)
		memTotal, _ := strconv.ParseFloat(values[2], 64)
		memFree, _ := strconv.ParseFloat(values[4], 64)
		diskFree, _ := strconv.ParseFloat(values[6], 64)

		memUsed := memTotal - memFree
		memUsage := float64(memUsed) / float64(memTotal) * 100
		diskFreeMb := diskFree / (1024 * 1024)
		availableNetMbit := netBytes * 8 / 1000000

		if memUsage > 80 {
			fmt.Printf("Memory usage too high: %d%%\n", int(memUsage+0.5))
		}
		if diskFreeMb < 32000 {
			fmt.Printf("Free disk space is too low: %d Mb left\n", int(diskFreeMb+0.5))
		}
		if availableNetMbit < 20 {
			fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", int(availableNetMbit+0.5))
		}
		if loadAvg > 30 {
			fmt.Printf("Load Average is too high: %.0f\n", loadAvg)
		}
	}
}
