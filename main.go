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

		// ВЫВОДИМ ВСЕ СООБЩЕНИЯ БЕЗ ПРИОРИТЕТА
		hasLoadAvg := loadAvg > 30
		hasMemory := false
		hasDisk := false  
		hasNetwork := false

		if totalMem > 0 {
			memoryUsage := float64(usedMem) / float64(totalMem) * 100
			hasMemory = memoryUsage > 80
		}

		if totalDisk > 0 {
			diskUsage := float64(usedDisk) / float64(totalDisk) * 100
			hasDisk = diskUsage > 90
		}

		if totalNet > 0 {
			netUsage := float64(usedNet) / float64(totalNet) * 100
			hasNetwork = netUsage > 90
		}

		// Выводим ВСЕ проблемы
		if hasLoadAvg {
			fmt.Printf("Load Average is too high: %.0f\n", loadAvg)
		}
		if hasMemory {
			memoryUsage := float64(usedMem) / float64(totalMem) * 100
			fmt.Printf("Memory usage too high: %d%%\n", int(memoryUsage))
		}
		if hasDisk {
			freeDiskMB := (totalDisk - usedDisk) / (1024 * 1024)
			fmt.Printf("Free disk space is too low: %d Mb left\n", freeDiskMB)
		}
		if hasNetwork {
			availableNetMbit := (totalNet - usedNet) / 1000000
			fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", availableNetMbit)
		}

		time.Sleep(5 * time.Second)
	}
}