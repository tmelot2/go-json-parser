/*
Tests sets of memory reads (from the same address) to see how many read ports are on my CPU's backend.

Output on my Zen3 Ryzen 5900x:

--- Read_x1 ---
Min: 837532 (0.226355ms) 4.553680gb/s
Max: 1989268 (0.537627ms) 1.917214gb/s
Avg: 893534 (0.241490ms) 4.268277gb/s

--- Read_x2 ---
Min: 837532 (0.226355ms) 4.553680gb/s
Max: 2123319 (0.573857ms) 1.796175gb/s
Avg: 886816 (0.239674ms) 4.300616gb/s

--- Read_x3 ---
Min: 837532 (0.226355ms) 4.553680gb/s
Max: 1868056 (0.504868ms) 2.041616gb/s
Avg: 881333 (0.238193ms) 4.327368gb/s

--- Read_x4 ---
Min: 1673510 (0.452289ms) 2.278954gb/s
Max: 3000811 (0.811011ms) 1.270941gb/s
Avg: 1769710 (0.478289ms) 2.155073gb/s

As we can plainly see, there's a wall after 3x read. This is confirmed by the AMD Zen3 architecture
manual which states that the Load-Store unit can do 3 memory uops per cycle.
*/

package main

/*
// Compiler flags: -I. look for .h files in cur dir. Not needed here because I put code below.
// #cgo CFLAGS: -I.

// Linker flags: -L. look for libraries in cur dir. -ltheName link against file "theName".
#cgo LDFLAGS: -L. -lmemoryReadMaxPorts

// Used as a wrapper so that we can all asm routines without making a Go wrapper for each.
#include <stdint.h>
typedef void (*ASMFuncPtr)(uint64_t count, char* data);
void callASMFunction(ASMFuncPtr func, uint64_t count, char* data) {
    func(count, data);
}

typedef char u8;
typedef long long unsigned u64;

// Prototypes
void Read_x1(u64 count, u8 *data);
void Read_x2(u64 count, u8 *data);
void Read_x3(u64 count, u8 *data);
void Read_x4(u64 count, u8 *data);
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
	testFunctions := [4]TestFunction{
		{name: "Read_x1", fun: wrapASMFunction(C.Read_x1)},
		{name: "Read_x2", fun: wrapASMFunction(C.Read_x2)},
		{name: "Read_x3", fun: wrapASMFunction(C.Read_x3)},
		{name: "Read_x4", fun: wrapASMFunction(C.Read_x4)},
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
