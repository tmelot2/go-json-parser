// Tests different implementations of NOP loops to illustrate CPU front end bottlenecks.
// TODO: Simplify by removing wrapper funcs
// TODO: Do asm

package main

/*
// Compiler flags: -I. look for .h files in cur dir. Not needed here because I put code below.
// #cgo CFLAGS: -I.

// Linker flags: -L. look for libraries in cur dir. -ltheName link against file "theName".
#cgo LDFLAGS: -L. -lnopLoopPatterns

// Used as a generic wrapper so that we can all any asm routine without making a Go wrapper.
#include <stdint.h>
typedef void (*ASMFuncPtr)(uint64_t count, char* data);
void callASMFunction(ASMFuncPtr func, uint64_t count, char* data) {
    func(count, data);
}

typedef char u8;
typedef long long unsigned u64;

// Prototypes
void NOP3x1AllBytes(u64 count, u8 *data);
void NOP1x3AllBytes(u64 count, u8 *data);
void NOP1x9AllBytes(u64 count, u8 *data);
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

type ASMFunction func(count C.ulonglong, data *C.char)
type TestFunction struct {
	name string
	fun  ASMFunction
}

func wrapASMFunction(f unsafe.Pointer) ASMFunction {
    return func(count C.ulonglong, data *C.char) {
        C.callASMFunction(C.ASMFuncPtr(f), C.ulonglong(count), data)
    }
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
	testFunctions := [3]TestFunction{
		{name: "NOP3x1AllBytes", fun: wrapASMFunction(C.NOP3x1AllBytes)},
		{name: "NOP1x3AllBytes", fun: wrapASMFunction(C.NOP1x3AllBytes)},
		{name: "NOP1x9AllBytes", fun: wrapASMFunction(C.NOP1x9AllBytes)},
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
			tester := testers[i]

			fmt.Println("---", testFunc.name, "---")
			secondsToTry := uint32(3)
			tester.NewTestWave(byteCount, cpuFreq, secondsToTry)

			for tester.IsTesting() {
				destBuffer := params.dest
				count := C.ulonglong(len(destBuffer))
				cBytes := (*C.char)(unsafe.Pointer(&destBuffer[0]))

				tester.BeginTime()
				testFunc.fun(count, cBytes)
				tester.EndTime()
				tester.CountBytes(uint64(len(destBuffer)))
			}
		}
		fmt.Println("=========================================")
	}

	fmt.Println("\nDone!")
}
