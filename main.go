package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	lcd, err := NewLCD(0x27, 1)
	if err != nil {
		log.Fatal(err)
	}
	defer lcd.Close()

	lcd.Init()

	for {
		cpuPercentage := fmt.Sprintf("CPU: %.2f%%", GetCPUPercentage())
		memUsed, memTotal := GetRAMUsage()
		ramUsage := fmt.Sprintf("RAM: %.1f/%.1fGb", memUsed, memTotal)

		lcd.WriteString(cpuPercentage, LCD_LINE_1)
		lcd.WriteString(ramUsage, LCD_LINE_2)
		time.Sleep(time.Second * 5)

		temperature := fmt.Sprintf("Temp: +%.2fC", GetTemperature())
		storageUsed, storageTotal := GetStorageUsage()
		storageUsage := fmt.Sprintf("Mem: %d/%dGb", storageUsed, storageTotal)

		lcd.WriteString(temperature, LCD_LINE_1)
		lcd.WriteString(storageUsage, LCD_LINE_2)
		time.Sleep(time.Second * 4)
	}
}