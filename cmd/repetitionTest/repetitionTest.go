package main

import (
	"flag"
	"fmt"
	"os"

	"tmelot.jsonparser/internal/profiler"
	"tmelot.jsonparser/internal/repetitionTester"
)

func readViaOSStat(rt *repetitionTester.RepetitionTester, fileName string, byteCount uint64) {
	for rt.IsTesting() {
		// Read file
		rt.BeginTime()
		_, err := os.ReadFile(fileName)
		rt.EndTime()
		if err != nil {
			msg := fmt.Sprintln("Error:", err)
			panic(msg)
		}
		rt.CountBytes(byteCount)
	}
}

func main() {
	fileNameArg := flag.String("fileName", "../../pairs.json", "Path to pairs JSON file")
	flag.Parse()
	fileName := *fileNameArg

	rt := repetitionTester.NewRepetitionTester()
	cpuFreq := profiler.EstimateCPUTimerFreq(false)

	// Get file size for bandwidth calculation purposes.
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		fmt.Println("Error:", err)
	}
	byteCount := uint64(fileInfo.Size())

	for i := 0; i < 3; i++ {
		// bytes := uint64(0)
		secondsToTry := uint32(2)
		rt.NewTestWave(byteCount, cpuFreq, secondsToTry)
		readViaOSStat(rt, fileName, byteCount)
		fmt.Println("=========================================")
	}

	fmt.Println("\nDone!")
}
