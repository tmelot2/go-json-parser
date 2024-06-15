// +build realjsonparser

/*
	Uses real JSON parsing library that I am writing.
*/

package main

import (
	"flag"
	"fmt"
	"os"
	// "strings"
	// "strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)


func DebugPrintf(format string, a ...interface{}) {
	if DEBUG {
		fmt.Printf(format, a...)
	}
}

func DebugPrintln(a ...interface{}) {
	if DEBUG {
		fmt.Println(a...)
	}
}

func GetPrinter() *message.Printer {
	return message.NewPrinter(language.English) // For printing large numbers with commas
}

func PrintTimeElapsed(label string, totalTSCElapsed, start, end uint64) {
	elapsed := end - start
	percent := float64(100.0 * (float64(elapsed) / float64(totalTSCElapsed)))
	p := GetPrinter()
	s := p.Sprintf("%*d", OUTPUT_WIDTH, elapsed)
	p.Printf("  %s: %14s (%.2f%%)\n", label, s, percent)
}

var DEBUG = false
var OUTPUT_WIDTH = 10


// Main
//
func main() {
	// Profiling vars
	var (
		prof_begin         = uint64(0)
		prof_read          = uint64(0)
		prof_read_tostr    = uint64(0)
		prof_miscSetup     = uint64(0)
		prof_parse         = uint64(0)
		prof_sum           = uint64(0)
		prof_miscOutput    = uint64(0)
		prof_end          = uint64(0)
		// TODO: Add steps for JsonValue stuff, even though it's mixed around
	)
	// Setup
	const EARTH_RADIUS = 6372.8
	p := GetPrinter()

	// cpuFreq := EstimateCPUTimerFreq(true)

	// Get input args
	inputFileArg := flag.String("input", "pairs.json", "Name of input file containing point pairs")
	flag.Parse()

	prof_begin = ReadCPUTimer()

	// Read JSON file, convert to string
	prof_read = ReadCPUTimer()
	data, err := os.ReadFile(*inputFileArg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	prof_read_tostr = ReadCPUTimer()
	strData := string(data)
	prof_miscSetup = ReadCPUTimer()
	DebugPrintln(strData)

	// Parse
	prof_parse = ReadCPUTimer()
	jsonResult, err := ParseJson(strData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Loop over JSON to do stuff
	prof_sum = ReadCPUTimer()
	fmt.Println("===============================")
	haversineSum := 0.0
	pairs, _ := jsonResult.GetArray("pairs")
	for _, p := range pairs {
		x0, _ := p.GetFloat("x0")
		y0, _ := p.GetFloat("y0")
		x1, _ := p.GetFloat("x1")
		y1, _ := p.GetFloat("y1")
		haversineSum += referenceHaversine(x0, y0, x1, y1, EARTH_RADIUS)
	}
	avg := haversineSum / float64(len(pairs))

	prof_miscOutput = ReadCPUTimer()
	p.Printf("Count: %*d\nHaversine sum: %.16f\nHaversine avg: %.16f\n", 14, len(pairs), haversineSum, avg)

	prof_end = ReadCPUTimer()

	// Calculate & print profile timings
	totalCPUElapsed := prof_end - prof_begin

	cpuFreq := EstimateCPUTimerFreq(false)
	ms := 1000.0 * float64(totalCPUElapsed) / float64(cpuFreq)
	fmt.Println("\n[CPU profiling stats]")
	p.Printf("Total time: %0.4fms (CPU freq %*d)\n", ms, 14, cpuFreq)
	PrintTimeElapsed("Startup     ", totalCPUElapsed, prof_begin, prof_read)
	PrintTimeElapsed("Read        ", totalCPUElapsed, prof_read, prof_read_tostr)
	PrintTimeElapsed("Read ToStr  ", totalCPUElapsed, prof_read_tostr, prof_miscSetup)
	PrintTimeElapsed("Misc Setup  ", totalCPUElapsed, prof_miscSetup, prof_parse)
	PrintTimeElapsed("Parse       ", totalCPUElapsed, prof_parse, prof_sum)
	PrintTimeElapsed("Sum         ", totalCPUElapsed, prof_sum, prof_miscOutput)
	PrintTimeElapsed("Misc Output ", totalCPUElapsed, prof_miscOutput, prof_end)
}
