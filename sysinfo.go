package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

func getCPUSample() (uint64, uint64) {
	var idle, total uint64

	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		log.Fatalf("unable to read /proc/stat: %s", err)
	}

	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)

			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}

				total += val
				if i == 4 {
					idle = val
				}
			}

			return idle, total
		}
	}

	return idle, total
}

func GetCPUPercentage() float64 {
	idle0, total0 := getCPUSample()
	time.Sleep(time.Second * 1)
	idle1, total1 := getCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)

	result := 100 * (totalTicks - idleTicks) / totalTicks

	if math.IsNaN(result) {
		return 0.0
	}

	return result
}

func GetRAMUsage() (memUsed float32, memTotal float32) {
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		log.Fatalf("unable to read /proc/stat: %s", err)
	}

	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)

		if fields[0] == "MemTotal:" {
			raw, err := strconv.Atoi(fields[1])
			if err != nil {
				log.Fatal(err)
			}

			memTotal = convertKBtoGB(raw)
		}

		if fields[0] == "MemFree:" {
			raw, err := strconv.Atoi(fields[1])
			if err != nil {
				log.Fatal(err)
			}

			memFree := convertKBtoGB(raw)
			memUsed = memTotal - memFree

			return
		}
	}

	return
}