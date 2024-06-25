package main

import (
	"fmt"
	"os"

	"tmelot.jsonparser/internal/profiler"
	"tmelot.jsonparser/internal/repetitionTester"
)

func readViaOSStat(rt *repetitionTester.RepetitionTester, fileName string) {
	// Get file size for bandwidth calculation purposes.
	// fileInfo, err := os.Stat(fileName)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// }

	for rt.IsTesting() {
		// Read file
		rt.BeginTime()
		_, err := os.ReadFile(fileName)
		rt.EndTime()
		if err != nil {
			msg := fmt.Sprintln("Error:", err)
			panic(msg)
		}
	}
}

func main() {
	rt := repetitionTester.NewRepetitionTester()
	cpuFreq := profiler.EstimateCPUTimerFreq(false)

	for i := 0; i < 3; i++ {
		bytes := uint64(0)
		secondsToTry := uint32(1)
		rt.NewTestWave(bytes, cpuFreq, secondsToTry)
		readViaOSStat(rt, "../../pairs.json")
		fmt.Println("=========================================")
	}

	fmt.Println("\nDone!")
}
