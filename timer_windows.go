// +build windows

package main

import (
	"fmt"
)

func ReadCpuTimer() int64 {
	return 0
}

func GetOSTimerFreq() (int64, error) {
	return 0, nil
}

func ReadOSTimer() (int64, error) {
	return 0, nil
}

func PrintTimerStats() {
	fmt.Println("[CPU timer stats]")
	fmt.Println("Windows not yet supported! All functions return 0.")
}