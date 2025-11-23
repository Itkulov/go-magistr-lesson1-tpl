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
	for {
		resp, err := http.Get("http://srv.msk01.gigacorp.local/_stats")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		values := strings.Split(strings.TrimSpace(string(body)), ",")
		if len(values) != 7 {
			fmt.Printf("Invalid data format\n")
			time.Sleep(5 * time.Second)
			continue
		}

		loadAvg, _ := strconv.ParseFloat(values[0], 64)
		totalMem, _ := strconv.ParseUint(values[1], 10, 64)
		usedMem, _ := strconv.ParseUint(values[2], 10, 64)
		totalDisk, _ := strconv.ParseUint(values[3], 10, 64)
		usedDisk, _ := strconv.ParseUint(values[4], 10, 64)
		totalNet, _ := strconv.ParseUint(values[5], 10, 64)
		usedNet, _ := strconv.ParseUint(values[6], 10, 64)

		var message string

		memoryUsage := float64(usedMem) / float64(totalMem) * 100
		if memoryUsage >= 98 {
			message = fmt.Sprintf("Memory usage too high: %.0f%%", memoryUsage)
		}

		if message == "" && loadAvg > 30 {
			message = fmt.Sprintf("Load Average is too high: %.0f", loadAvg)
		}

		if message == "" {
			freeDiskMB := (totalDisk - usedDisk) / (1024 * 1024)
			if freeDiskMB < 25000 {
				message = fmt.Sprintf("Free disk space is too low: %d Mb left", freeDiskMB)
			}
		}

		if message == "" {
			availableNetMbit := (totalNet - usedNet) / 1000000
			if availableNetMbit < 200 {
				message = fmt.Sprintf("Network bandwidth usage high: %d Mbit/s available", availableNetMbit)
			}
		}

		if message != "" {
			fmt.Println(message)
		}

		time.Sleep(5 * time.Second)
	}
}