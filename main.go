package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	errorCount := 0
	
	for {
		resp, err := http.Get("http://srv.msk01.gigacorp.local/_stats")
		if err != nil || resp.StatusCode != 200 {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
				return
			}
			time.Sleep(5 * time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
				return
			}
			time.Sleep(5 * time.Second)
			continue
		}

		values := strings.Split(strings.TrimSpace(string(body)), ",")
		if len(values) != 7 {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
				return
			}
			time.Sleep(5 * time.Second)
			continue
		}

		errorCount = 0

		loadAvg, _ := strconv.ParseFloat(values[0], 64)
		totalMem, _ := strconv.ParseUint(values[1], 10, 64)
		usedMem, _ := strconv.ParseUint(values[2], 10, 64)
		totalDisk, _ := strconv.ParseUint(values[3], 10, 64)
		usedDisk, _ := strconv.ParseUint(values[4], 10, 64)
		totalNet, _ := strconv.ParseUint(values[5], 10, 64)
		usedNet, _ := strconv.ParseUint(values[6], 10, 64)

		messages := []string{}

		// Load Average > 30
		if loadAvg > 30 {
			messages = append(messages, fmt.Sprintf("Load Average is too high: %.0f", loadAvg))
		}

		// Memory usage > 80%
		if totalMem > 0 {
			memoryUsage := float64(usedMem) / float64(totalMem) * 100
			if memoryUsage > 80 {
				// Используем math.Round для точного округления
				messages = append(messages, fmt.Sprintf("Memory usage too high: %.0f%%", math.Round(memoryUsage)))
			}
		}

		// Disk usage > 90%
		if totalDisk > 0 {
			diskUsage := float64(usedDisk) / float64(totalDisk) * 100
			if diskUsage > 90 {
				freeDiskMB := (totalDisk - usedDisk) / (1024 * 1024)
				messages = append(messages, fmt.Sprintf("Free disk space is too low: %d Mb left", freeDiskMB))
			}
		}

		// Network usage > 90%
		if totalNet > 0 {
			netUsage := float64(usedNet) / float64(totalNet) * 100
			if netUsage > 90 {
				availableNetMbit := (totalNet - usedNet) / 1000000
				messages = append(messages, fmt.Sprintf("Network bandwidth usage high: %d Mbit/s available", availableNetMbit))
			}
		}

		for _, msg := range messages {
			fmt.Println(msg)
		}

		time.Sleep(5 * time.Second)
	}
}