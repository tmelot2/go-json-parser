// +build darwin

package profiler

import (
	"fmt"

	"golang.org/x/sys/unix"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Declare asm functions that return CPU timer values.
// NOTE-1: The ARM impl uses CNTVCT which returns a 24MHz counter timer on MacOS.
// Without a cycle counter from the CPU & OS, the 24Mhz doesn't really help us
// count cycles. BUT, the relative percentages of cycle time is still useful!
func ReadCPUTimer() uint64

// GetOSTimerFreq returns the frequency of the OS timer.
func GetOSTimerFreq() (uint64, error) {
    // On Unix-like systems, the frequency can be considered as nanoseconds in a second
    return 1e9, nil // 1 second = 1e9 nanoseconds
}

// ReadOSTimer returns the current time from the OS high-resolution timer.
func ReadOSTimer() (uint64, error) {
    var ts unix.Timespec
    err := unix.ClockGettime(unix.CLOCK_MONOTONIC, &ts)
    if err != nil {
        return 0, err
    }
    osTimerFreq, _ := GetOSTimerFreq()
    return osTimerFreq * uint64(ts.Sec) + uint64(ts.Nsec), nil // Convert to nanoseconds
}

// Prints read, measurement, & guess of CPU frequency & related data.
func EstimateCPUTimerFreq(printDebug bool) uint64 {
	// Setup
	millisecondsToWait := uint64(100)
	width := 20 // Output width
	p := message.NewPrinter(language.English) // For printing large numbers with commas

	// Get OS timer frequency
	osFreq, err := GetOSTimerFreq()
	if err != nil {
		fmt.Println("Error getting OS timer frequency:", err)
		return 0
	}
	// In nanoseconds per second
	if printDebug {
		p.Printf("OS Timer Frequency [reported]: %*d\n", width, osFreq)
	}

	cpuStart   := ReadCPUTimer()
	osStart, _ := ReadOSTimer()

	var osEnd uint64
	var osElapsed uint64
	osWaitTime := osFreq * millisecondsToWait / 1000
	for osElapsed < osWaitTime {
		osEnd, _ = ReadOSTimer()
		osElapsed = osEnd - osStart
	}
	cpuEnd := ReadCPUTimer()
	cpuElapsed := cpuEnd - cpuStart

	cpuFreq := uint64(0)
	if osElapsed > 0 {
		cpuFreq = osFreq * cpuElapsed / osElapsed
	}

	if printDebug {
		p.Printf("OS Timer:                      %*d elapsed\n", width, osElapsed)
		p.Printf("OS Seconds (elapsed/freq):          %*.4f\n", width, float64(osElapsed) / float64(osFreq))

		p.Printf("CPU timer:                     %*d elapsed\n", width, cpuElapsed)
		p.Printf("CPU freq (guessed):            %*d\n", width, cpuFreq)
	}

	return uint64(cpuFreq)
}