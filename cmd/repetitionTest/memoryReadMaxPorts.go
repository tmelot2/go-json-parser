/*
Tests sets of memory reads (from the same address) to see how many read ports are on my CPU's backend.

Output on my Zen3 Ryzen 5900x:

--- Read_x1 ---
Min: 860831307 (232.650340ms) 4.391733gb/s
Max: 894005366 (241.616041ms) 4.228768gb/s
Avg: 870280652 (235.204143ms) 4.344049gb/s PF: 10 (107137.0208k/fault)

--- Read_x2 ---
Min: 432655763 (116.930587ms) 8.737990gb/s
Max: 441937767 (119.439164ms) 8.554466gb/s
Avg: 433870482 (117.258880ms) 8.713526gb/s

--- Read_x3 ---
Min: 288269701 (77.908463ms) 13.114599gb/s
Max: 304143404 (82.198528ms) 12.430128gb/s
Avg: 289645586 (78.280313ms) 13.052302gb/s

--- Read_x4 ---
Min: 433797361 (117.239118ms) 8.714994gb/s
Max: 446228049 (120.598666ms) 8.472219gb/s
Avg: 434966322 (117.555044ms) 8.691573gb/s

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
	fileNameArg  := flag.String("fileName", "../../pairs10m.json", "Path to pairs JSON file")
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
