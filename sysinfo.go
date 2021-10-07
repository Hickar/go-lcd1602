package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func GetCPUPercentage() float64 {
	idle0, total0 := GetCPUSample()
	time.Sleep(time.Second * 1)
	idle1, total1 := GetCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)

	usage := 100 * (totalTicks - idleTicks) / totalTicks

	if math.IsNaN(usage) {
		return 0.0
	}

	return usage
}

func GetCPUSample() (uint64, uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		log.Fatalf("unable to read /proc/stat: %s", err)
	}

	var idle, total uint64
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

func GetRAMUsage() (float32, float32) {
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		log.Fatalf("unable to read /proc/stat: %s", err)
	}

	var memUsed, memTotal float32
	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)

		if fields[0] == "MemTotal:" {
			raw, err := strconv.Atoi(fields[1])
			if err != nil {
				log.Fatal(err)
			}

			memTotal = ConvertKBtoGB(raw)
		}

		if fields[0] == "MemFree:" {
			raw, err := strconv.Atoi(fields[1])
			if err != nil {
				log.Fatal(err)
			}

			memFree := ConvertKBtoGB(raw)
			memUsed = memTotal - memFree

			return memUsed, memTotal
		}
	}

	return memUsed, memTotal
}

func GetTemperature() float32 {
	stdout, stderr, err := execute("sensors")
	if err != nil {
		log.Fatalf("unable to get cpu temperature: %s", stderr.String())
	}

	var tempMeasurements []float32
	lines := strings.Split(stdout.String(), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)

		if len(fields) > 0 {
			if strings.Contains(fields[0], "Core") {
				tempMeasurement, err := strconv.ParseFloat(fields[2][:len(fields[2]) - 3], 32)
				if err != nil {
					log.Fatal(err)
				}

				tempMeasurements = append(tempMeasurements, float32(tempMeasurement))
			}
		}
	}

	return FindAverage(tempMeasurements)
}

func GetStorageUsage() (int, int) {
	stdout, stderr, err := execute("df", "/")
	if err != nil {
		log.Fatalf("unable to get storage info: %s", stderr.String())
	}

	var memUsed, memTotal int
	lines := strings.Split(stdout.String(), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)

		if strings.Contains(fields[0], "/dev") {
			raw, err := strconv.Atoi(fields[1])
			if err != nil {
				log.Fatal(err)
			}
			memTotal = int(ConvertKBtoGB(raw))

			raw, err = strconv.Atoi(fields[3])
			if err != nil {
				log.Fatal(err)
			}
			memAvailable := int(ConvertKBtoGB(raw))
			memUsed = memTotal - memAvailable

			return memUsed, memTotal
		}
	}

	return memUsed, memTotal
}

func execute(name string, arg ...string) (bytes.Buffer, bytes.Buffer, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(name, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("stdErr: %s\n", stderr.String())
		return stdout, stderr, err
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("stdErr: %s\n", stderr.String())
		return stdout, stderr, err
	}

	return stdout, stderr, nil
}