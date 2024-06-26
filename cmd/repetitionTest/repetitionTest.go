package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"tmelot.jsonparser/internal/profiler"
	"tmelot.jsonparser/internal/repetitionTester"
)

func readViaOSReadFile(rt *repetitionTester.RepetitionTester, fileName string, byteCount uint64) {
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

func readViaIOUtilReadFile(rt *repetitionTester.RepetitionTester, fileName string, byteCount uint64) {
	for rt.IsTesting() {
		// Read file
		rt.BeginTime()
		_, err := ioutil.ReadFile(fileName)
		rt.EndTime()
		if err != nil {
			msg := fmt.Sprintln("Error:", err)
			panic(msg)
		}
		rt.CountBytes(byteCount)
	}
}

func readViaBufIOReader(rt *repetitionTester.RepetitionTester, fileName string, byteCount uint64) {
	for rt.IsTesting() {
		// Read file
		file, err := os.Open(fileName)
		if err != nil {
			msg := fmt.Sprintln("Error:", err)
			panic(msg)
		}
		reader := bufio.NewReader(file)

		rt.BeginTime()
		_, err = io.ReadAll(reader)
		rt.EndTime()
		if err != nil {
			msg := fmt.Sprintln("Error:", err)
			panic(msg)
		}
		rt.CountBytes(byteCount)
	}
}

type ReadOverheadTestFunc func(*repetitionTester.RepetitionTester, string, uint64)
type TestFunction struct {
	name string
	fun  ReadOverheadTestFunc
}


func main() {
	fileNameArg := flag.String("fileName", "../../pairs.json", "Path to pairs JSON file")
	flag.Parse()
	fileName := *fileNameArg

	// Table of test functions to test.
	testFunctions := [3]TestFunction{
		{name: "OS.ReadFile", fun: readViaOSReadFile},
		{name: "ioutil.ReadFile", fun: readViaIOUtilReadFile},
		{name: "bufio.Reader", fun: readViaBufIOReader},
	}

	// Create multiple testers, one for each test function.
	var testers [len(testFunctions)]*repetitionTester.RepetitionTester
	for i,_ := range testers {
		testers[i] = repetitionTester.NewRepetitionTester()
	}

	// rt := repetitionTester.NewRepetitionTester()
	cpuFreq := profiler.EstimateCPUTimerFreq(false)

	// Get file size for bandwidth calculation purposes.
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		fmt.Println("Error:", err)
	}
	byteCount := uint64(fileInfo.Size())

	// Run tests!
	for i := 0; i < 3; i++ {
		for i, testFunc := range testFunctions {
			fmt.Println(testFunc.name, ":")
			secondsToTry := uint32(2)
			testers[i].NewTestWave(byteCount, cpuFreq, secondsToTry)
			testFunc.fun(testers[i], fileName, byteCount)
			fmt.Println("")
		}
		fmt.Println("=========================================")
	}

	fmt.Println("\nDone!")
}
