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
		rt.Print()
	}
}

func main() {
	rt := repetitionTester.NewRepetitionTester()
	fmt.Println(rt)

	for i := 0; i < 40; i++ {
		rt.NewTestWave(0, profiler.EstimateCPUTimerFreq(false), 3)
		readViaOSStat(rt, "../../pairs.json")
	}

	fmt.Println("\nDone!")
}
