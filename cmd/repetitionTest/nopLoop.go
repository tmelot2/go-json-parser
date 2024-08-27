// Tests memory touch speed vs different forms of reduced instructions.

package main

/*
// Compiler flags: -I. look for .h files in cur dir. Not needed here because I put code below.
// #cgo CFLAGS: -I.

// Linker flags: -L. look for libraries in cur dir. -ltheName link against file "theName".
#cgo LDFLAGS: -L. -lnopLoop

typedef char u8;
typedef long long unsigned u64;

// Prototypes
void MOVAllBytesASM(u64 count, u8 *data);
void NOPAllBytesASM(u64 count);
void CMPAllBytesASM(u64 count);
void DECAllBytesASM(u64 count);
*/
import "C"

import (
	"flag"
	"fmt"
	"os"
	// "runtime/debug"
	"unsafe"

	"tmelot.jsonparser/internal/profiler"
	"tmelot.jsonparser/internal/repetitionTester"
)

type ReadParams struct {
	dest      []byte
	fileName  string
}

func writeToAllBytes(rt *repetitionTester.RepetitionTester, params *ReadParams) {
	for rt.IsTesting() {
		destBuffer := params.dest

		rt.BeginTime()
		for i := 0; i < len(destBuffer); i++ {
			destBuffer[i] = uint8(i)
		}
		rt.EndTime()

		rt.CountBytes(uint64(len(destBuffer)))
	}
}

func MOVAllBytes(rt *repetitionTester.RepetitionTester, params *ReadParams) {
	for rt.IsTesting() {
		destBuffer := params.dest
		count := uint64(len(destBuffer))

		cBytes := (*C.char)(unsafe.Pointer(&destBuffer[0]))
		rt.BeginTime()
		C.MOVAllBytesASM(C.ulonglong(count), cBytes)
		rt.EndTime()
		rt.CountBytes(uint64(len(destBuffer)))
	}
}

func NOPAllBytes(rt *repetitionTester.RepetitionTester, params *ReadParams) {
	for rt.IsTesting() {
		destBuffer := params.dest
		count := uint64(len(destBuffer))

		rt.BeginTime()
		C.NOPAllBytesASM(C.ulonglong(count))
		rt.EndTime()
		rt.CountBytes(uint64(len(destBuffer)))
	}
}

func CMPAllBytes(rt *repetitionTester.RepetitionTester, params *ReadParams) {
	for rt.IsTesting() {
		destBuffer := params.dest
		count := uint64(len(destBuffer))

		rt.BeginTime()
		C.CMPAllBytesASM(C.ulonglong(count))
		rt.EndTime()
		rt.CountBytes(uint64(len(destBuffer)))
	}
}

func DECAllBytes(rt *repetitionTester.RepetitionTester, params *ReadParams) {
	for rt.IsTesting() {
		destBuffer := params.dest
		count := uint64(len(destBuffer))

		rt.BeginTime()
		C.DECAllBytesASM(C.ulonglong(count))
		rt.EndTime()
		rt.CountBytes(uint64(len(destBuffer)))
	}
}


type TestFunc func(*repetitionTester.RepetitionTester, *ReadParams)
type TestFunction struct {
	name string
	fun  TestFunc
}

func HandleError(err error) {
	msg := fmt.Sprintln("Error:", err)
	panic(msg)
}

func main() {
	// Turn off the garbage collector. This is a short-running app, & the testing needs to be done
	// without the GC doing sensible things like reusing memory with make().
	// debug.SetGCPercent(-1)

	// Input args
	fileNameArg  := flag.String("fileName", "../../pairs.json", "Path to pairs JSON file")
	flag.Parse()
	fileName  := *fileNameArg

	// Table of test functions to test.
	testFunctions := [5]TestFunction{
		{name: "WriteToAllBytes", fun: writeToAllBytes},
		{name: "MOVAllBytes", fun: MOVAllBytes},
		{name: "NOPAllBytes", fun: NOPAllBytes},
		{name: "CMPAllBytes", fun: CMPAllBytes},
		{name: "DECAllBytes", fun: DECAllBytes},
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
		dest:      make([]byte, byteCount),
		fileName:  fileName,
	}

	// Run tests!
	for i := 0; i < 1; i++ {
		// for true {
		for i, testFunc := range testFunctions {
			fmt.Println("---", testFunc.name, "---")
			secondsToTry := uint32(3)
			testers[i].NewTestWave(byteCount, cpuFreq, secondsToTry)

			testFunc.fun(testers[i], &params)
			// fmt.Println("")
		}
		fmt.Println("=========================================")
	}

	fmt.Println("\nDone!")
}
