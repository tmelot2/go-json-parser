// Tests "stars align" best-possible speed for reading files different ways with
// the standard library.

package main

import (
	"bufio"
	"bytes"
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
			HandleError(err)
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
			HandleError(err)
		}
		rt.CountBytes(byteCount)
	}
}

func readViaBufIOReader(rt *repetitionTester.RepetitionTester, fileName string, byteCount uint64) {
	for rt.IsTesting() {
		// Read file
		file, err := os.Open(fileName)
		if err != nil {
			HandleError(err)
		}
		defer file.Close()
		reader := bufio.NewReader(file)

		rt.BeginTime()
		_, err = io.ReadAll(reader)
		rt.EndTime()
		if err != nil {
			HandleError(err)
		}
		rt.CountBytes(byteCount)
	}
}

func readViaBytesBuffer(rt *repetitionTester.RepetitionTester, fileName string, byteCount uint64) {
	for rt.IsTesting() {
		// Read file
		file, err := os.Open(fileName)
		if err != nil {
			HandleError(err)
		}

		var buf bytes.Buffer

		rt.BeginTime()
		_, err = io.Copy(&buf, file)
		rt.EndTime()
		if err != nil {
			HandleError(err)
		}
		rt.CountBytes(byteCount)
	}
}

type ReadOverheadTestFunc func(*repetitionTester.RepetitionTester, string, uint64)
type TestFunction struct {
	name string
	fun  ReadOverheadTestFunc
}

func HandleError(err error) {
	msg := fmt.Sprintln("Error:", err)
	panic(msg)
}


func main() {
	fileNameArg := flag.String("fileName", "../../pairs.json", "Path to pairs JSON file")
	flag.Parse()
	fileName := *fileNameArg

	// Table of test functions to test.
	testFunctions := [4]TestFunction{
		{name: "OS.ReadFile", fun: readViaOSReadFile},
		{name: "ioutil.ReadFile", fun: readViaIOUtilReadFile},
		{name: "bufio.Reader", fun: readViaBufIOReader},
		{name: "bytes.Buffer", fun: readViaBytesBuffer},
	}

	// Create multiple testers, one for each test function.
	var testers [len(testFunctions)]*repetitionTester.RepetitionTester
	for i,_ := range testers {
		testers[i] = repetitionTester.NewRepetitionTester()
	}

	cpuFreq := profiler.EstimateCPUTimerFreq(false)

	// Get file size for bandwidth calculation purposes.
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		fmt.Println("Error:", err)
	}
	byteCount := uint64(fileInfo.Size())

	// Run tests!
	for i := 0; i < 1; i++ {
	// for true {
		for i, testFunc := range testFunctions {
			fmt.Println("---", testFunc.name, "---")
			secondsToTry := uint32(3)
			testers[i].NewTestWave(byteCount, cpuFreq, secondsToTry)
			testFunc.fun(testers[i], fileName, byteCount)
			fmt.Println("")
		}
		fmt.Println("=========================================")
	}

	fmt.Println("\nDone!")
}
