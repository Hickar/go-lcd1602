package main

import (
	"fmt"
	"strconv"
)

func FindAverage(arr []float32) float32 {
	return FindSum(arr) / float32(len(arr))
}

func FindSum(arr []float32) float32 {
	var sum float32
	for _, value := range arr {
		sum += value
	}

	return sum
}

func ConvertKBtoGB(num int) float32 {
	return float32(num) / 1024 / 1024
}

func LeftAlign(str string) string {
	return fmt.Sprintf("%-" + strconv.Itoa(WIDTH) + "v", str)
}

func RightAlign(str string) string {
	return fmt.Sprintf("%" + strconv.Itoa(WIDTH) + "v", str)
}