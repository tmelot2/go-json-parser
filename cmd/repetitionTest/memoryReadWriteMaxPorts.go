/*
Tests sets of memory reads & writes (from the same address) to see how many read & write ports are on my CPU's backend.

Output on my Zen3 Ryzen 5900x:

--- Read_x1 ---
Min: 829961208 (224.307593ms) 4.555077gb/s
Max: 855320711 (231.161322ms) 4.420023gb/s
Avg: 836111407 (225.969763ms) 4.521571gb/s PF: 5 (214274.0416k/fault)

--- Read_x2 ---
Min: 417322371 (112.786689ms) 9.059032gb/s
Max: 427022106 (115.408166ms) 8.853258gb/s
Avg: 418978102 (113.234171ms) 9.023232gb/s

--- Read_x3 ---
Min: 277102916 (74.890594ms) 13.643078gb/s
Max: 286161630 (77.338827ms) 13.211194gb/s
Avg: 279378168 (75.505510ms) 13.531969gb/s

--- Read_x4 ---
Min: 277795482 (75.077769ms) 13.609065gb/s
Max: 293696787 (79.375300ms) 12.872244gb/s
Avg: 285612182 (77.190332ms) 13.236609gb/s

--- Write_x1 ---
Min: 830239263 (224.382741ms) 4.553551gb/s
Max: 873916098 (236.186963ms) 4.325972gb/s
Avg: 843678161 (228.014775ms) 4.481018gb/s

--- Write_x2 ---
Min: 417625253 (112.868547ms) 9.052462gb/s
Max: 424076388 (114.612048ms) 8.914754gb/s
Avg: 419151657 (113.281077ms) 9.019496gb/s

--- Write_x3 ---
Min: 416096487 (112.455378ms) 9.085722gb/s
Max: 429491634 (116.075587ms) 8.802353gb/s
Avg: 417987042 (112.966325ms) 9.044627gb/s

--- Write_x4 ---
Min: 416248002 (112.496327ms) 9.082414gb/s
Max: 423301793 (114.402704ms) 8.931067gb/s
Avg: 418604758 (113.133271ms) 9.031280gb/s


For reads, as we can plainly see, there's a wall after 3x. This is confirmed by the AMD Zen3 architecture
manual which states that the Load-Store unit can do 3 load memory uops per cycle.

For writes, there's a wall after 2x, which is also confirmed by the manual which states that 2 of the 3
can be writes.
*/

package main

/*
// Compiler flags: -I. look for .h files in cur dir. Not needed here because I put code below.
// #cgo CFLAGS: -I.

// Linker flags: -L. look for libraries in cur dir. -ltheName link against file "theName".
#cgo LDFLAGS: -L. -lmemoryReadWriteMaxPorts

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
void Write_x1(u64 count, u8 *data);
void Write_x2(u64 count, u8 *data);
void Write_x3(u64 count, u8 *data);
void Write_x4(u64 count, u8 *data);
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
	testFunctions := [8]TestFunction{
		{name: "Read_x1", fun: wrapASMFunction(C.Read_x1)},
		{name: "Read_x2", fun: wrapASMFunction(C.Read_x2)},
		{name: "Read_x3", fun: wrapASMFunction(C.Read_x3)},
		{name: "Read_x4", fun: wrapASMFunction(C.Read_x4)},
		{name: "Write_x1", fun: wrapASMFunction(C.Write_x1)},
		{name: "Write_x2", fun: wrapASMFunction(C.Write_x2)},
		{name: "Write_x3", fun: wrapASMFunction(C.Write_x3)},
		{name: "Write_x4", fun: wrapASMFunction(C.Write_x4)},
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
