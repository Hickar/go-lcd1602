package main

import (
	"fmt"
	"log"
	"time"
)

// Temperature
// Storage Usage (Gb)

func main() {
	lcd, err := NewLCD(0x27, 1)
	if err != nil {
		log.Fatal(err)
	}

	lcd.Init()

	for {
		cpuPercentage := fmt.Sprintf("CPU: %.2f%", GetCPUPercentage())

		memUsed, memTotal := GetRAMUsage()
		ramUsage := fmt.Sprintf("RAM: %.2f/%.2fGb", memUsed, memTotal)

		lcd.WriteString(cpuPercentage, LINE_1)
		lcd.WriteString(ramUsage, LINE_2)
		time.Sleep(time.Second * 3)
	}
}

func convertKBtoGB(num int) float32 {
	return float32(num) / 1024 / 1024
}
