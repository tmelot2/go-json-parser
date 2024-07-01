// Tests "stars align" best-possible speed for reading files different ways with
// the standard library.

package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"unsafe"

	"tmelot.jsonparser/internal/profiler"
	"tmelot.jsonparser/internal/repetitionTester"
)

type AllocType int

const (
	AllocType_None AllocType = iota
	AllocType_Malloc

	AllocCount_Count
)

func handleAllocation(params *ReadParams, buffer *[]byte) {
	switch params.allocType {
	case AllocType_None:
		// fmt.Println("    allocType = NONE")
		break
	case AllocType_Malloc:
		// fmt.Println("    allocType = Malloc")
		*buffer = make([]byte, len(params.dest))
	default:
		HandleError(errors.New("Unrecognized allocation type"))
	}
}

type ReadParams struct {
	allocType AllocType
	dest      []byte
	fileName  string
}

func writeToAllBytes(rt *repetitionTester.RepetitionTester, params *ReadParams) {
	for rt.IsTesting() {
		destBuffer := params.dest
		fmt.Println("before: data starts at ", unsafe.Pointer(&(destBuffer)[0]))
		handleAllocation(params, &destBuffer)
		fmt.Println("after: data starts at ", unsafe.Pointer(&(destBuffer)[0]))

		rt.BeginTime()
		for i := 0; i < len(destBuffer); i++ {
			destBuffer[i] = uint8(i)
		}
		rt.EndTime()

		rt.CountBytes(uint64(len(destBuffer)))
	}
}

func readViaOSReadFile(rt *repetitionTester.RepetitionTester, params *ReadParams) {
	for rt.IsTesting() {
		// Read file
		rt.BeginTime()
		_, err := os.ReadFile(params.fileName)
		rt.EndTime()
		if err != nil {
			HandleError(err)
		}
		rt.CountBytes(uint64(len(params.dest)))
	}
}

func readViaOSReadFull(rt *repetitionTester.RepetitionTester, params *ReadParams) {
	for rt.IsTesting() {
		file, err := os.Open(params.fileName)
		if err != nil {
			HandleError(err)
		}
		defer file.Close()

		// fmt.Printf("%p\n", params.dest)
		buffer := params.dest
		handleAllocation(params, &buffer)

		rt.BeginTime()
		_, err = io.ReadFull(file, buffer)
		rt.EndTime()
		if err != nil {
			HandleError(err)
		}
		rt.CountBytes(uint64(len(buffer)))
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

type ReadOverheadTestFunc func(*repetitionTester.RepetitionTester, *ReadParams)
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
	testFunctions := [2]TestFunction{
		// {name: "OS.ReadFile", fun: readViaOSReadFile},
		{name: "WriteToAllBytes", fun: writeToAllBytes},
		{name: "OS.ReadFull", fun: readViaOSReadFull},
		// {name: "ioutil.ReadFile", fun: readViaIOUtilReadFile},
		// {name: "bufio.Reader", fun: readViaBufIOReader},
		// {name: "bytes.Buffer", fun: readViaBytesBuffer},
	}

	// Create multiple testers, one for each test function.
	var testers [len(testFunctions)]*repetitionTester.RepetitionTester
	for i, _ := range testers {
		testers[i] = repetitionTester.NewRepetitionTester()
	}

	cpuFreq := profiler.EstimateCPUTimerFreq(false)

	// Get file size for bandwidth calculation purposes.
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		fmt.Println("Error:", err)
	}
	byteCount := uint64(fileInfo.Size())

	params := ReadParams{
		// allocType: AllocType_None,
		allocType: AllocType_Malloc,
		dest:      make([]byte, byteCount),
		fileName:  fileName,
	}

	// Run tests!
	for i := 0; i < 1; i++ {
		// for true {
		for i, testFunc := range testFunctions {
			fmt.Println("---", testFunc.name, "---")
			secondsToTry := uint32(1)
			testers[i].NewTestWave(byteCount, cpuFreq, secondsToTry)

			testFunc.fun(testers[i], &params)
			fmt.Println("")
		}
		fmt.Println("=========================================")
	}

	fmt.Println("\nDone!")
}
