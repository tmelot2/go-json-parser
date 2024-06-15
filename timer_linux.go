// +build linux

package main

import (
	"golang.org/x/sys/unix"
)

// Declare asm functions that return CPU timer values.
// TODO: Only ARM64 is currently supported, add X64.
// Also, the ARM impl uses CNTVCT which returns a 24MHz timer on MacOS Sonoma.
// Without a cycle counter from the CPU & OS, the 24Mhz doesn't really help us.
func ReadCPUTimer() int64

// GetOSTimerFreq returns the frequency of the OS timer.
func GetOSTimerFreq() (int64, error) {
    // On Unix-like systems, the frequency can be considered as nanoseconds in a second
    return 1e9, nil // 1 second = 1e9 nanoseconds
}

// ReadOSTimer returns the current time from the OS high-resolution timer.
func ReadOSTimer() (int64, error) {
    var ts unix.Timespec
    err := unix.ClockGettime(unix.CLOCK_MONOTONIC, &ts)
    if err != nil {
        return 0, err
    }
    osTimerFreq, _ := GetOSTimerFreq()
    return osTimerFreq * ts.Sec + ts.Nsec, nil // Convert to nanoseconds
}

func PrintTimerStats() {
	// Setup
	millisecondsToWait := int64(10)
	// TODO: Optionally get ms from input args
	width := 24 // Output width
	p := message.NewPrinter(language.English) // For printing large numbers with commas

	osFreq, err := GetOSTimerFreq()
	if err != nil {
		fmt.Println("Error getting OS timer frequency:", err)
		return
	}
	// In nanoseconds per second
	p.Printf("OS Timer Frequency [reported]:          %*d\n", width, osFreq)

	cpuStart   := ReadCPUTimer()
	osStart, _ := ReadOSTimer()

	var osEnd int64
	var osElapsed int64
	osWaitTime := osFreq * millisecondsToWait / 1000
	for osElapsed < osWaitTime {
		osEnd, _ = ReadOSTimer()
		osElapsed = osEnd - osStart
	}
	cpuEnd := ReadCPUTimer()
	cpuElapsed := cpuEnd - cpuStart

	cpuFreq := int64(0)
	if osElapsed > 0 {
		cpuFreq = osFreq * cpuElapsed / osElapsed
	}

	// p.Printf(  "OS Timer:      %*d -> %*d = %*d elapsed\n", width, osStart, width, osEnd, width, osElapsed)
	p.Printf("OS Timer:                               %*d elapsed\n", width, osElapsed)
	p.Printf("OS Seconds (elapsed/freq):                   %*.4f\n", width, float64(osElapsed) / float64(osFreq))

	// p.Printf(  "CPU timer:     %*d -> %*d = %*d elapsed\n", width, cpuStart, width, cpuEnd, width, cpuElapsed)
	p.Printf("CPU timer:                              %*d elapsed\n", width, cpuElapsed)
	p.Printf("CPU freq (guessed):                     %*d\n", width, cpuFreq)
 }
