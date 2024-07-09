// Tests page fault patterns for increasingly larger memory allocations & touches.

package main

import (
	// "bufio"
	// "bytes"
	// "errors"
	// "flag"
	"fmt"
	// "io"
	// "io/ioutil"
	// "os"
	"runtime/debug"
	// // "unsafe"

	// "tmelot.jsonparser/internal/profiler"
	"tmelot.jsonparser/internal/repetitionTester"
)

func handleAllocation(size int) *[]byte {
	buffer := make([]byte, size)
	return &buffer
}

func main() {
	// Turn off the garbage collector. This is a short-running app, & the testing needs to be done
	// without the GC doing sensible things like reusing memory with make().
	debug.SetGCPercent(-1)

	// TODO: Change to input arg
	pageCount := 1024
	// NOTE: May not be OS page size, this is just our testing page size.
	pageSize := 4096
	totalSize := pageCount * pageSize

	// fmt.Println("Page Count, Touch Count, Fault Count, Extra Faults")
	fmt.Println("Test #, Touch Count, Fault Count, Extra Faults")

	for touchCount := 0; touchCount <= pageCount; touchCount++ {
		touchSize := pageSize * touchCount
		data := handleAllocation(totalSize)

		startFaultCount := repetitionTester.GetPageFaultCount()
		for i := 0; i < touchSize; i++ {
			(*data)[i] = byte(i)
		}
		endFaultCount := repetitionTester.GetPageFaultCount()
		faultCount := int(endFaultCount - startFaultCount)

		// fmt.Printf("%d, %d, %d, %d\n", pageCount, touchCount, faultCount, (faultCount - touchCount))
		fmt.Printf("%d, %d, %d, %d\n", touchCount, touchCount, faultCount, (faultCount - touchCount))
	}
}
