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
	// "time"
	// "unsafe"

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

var DEBUG = false
var OUTPUT_WIDTH = 10
var globalProfiler = newProfiler()
var globalProfilerParent string

func readEntireFile(fileName string) ([]byte, error) {
	// Get file size for bandwidth calculation purposes.
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Start profile.
	globalProfiler.StartBandwidth("Read", uint64(fileInfo.Size()))

	// Read file
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	globalProfiler.EndBandwidth("Read")
	return data, nil
}

var EARTH_RADIUS = 6372.8

func haversineSum(json *JsonValue) {
	p := GetPrinter()

	fmt.Println("===============================")
	haversineSum := 0.0
	// Profile to get time for GetArray() call.
	globalProfiler.StartBlock("SumHaversine")
	pairs, _ := json.GetArray("pairs")
	globalProfiler.EndBlock("SumHaversine")

	// Profile rest of haversine sum. 32 is bytes per haversine set. 4 points,
	// each a float64, so 8 bytes. 8*4 = 32.
	globalProfiler.StartBandwidth("SumHaversine", uint64(len(pairs)*32))
	for _, p := range pairs {
		x0, _ := p.GetFloat("x0")
		y0, _ := p.GetFloat("y0")
		x1, _ := p.GetFloat("x1")
		y1, _ := p.GetFloat("y1")
		haversineSum += referenceHaversine(x0, y0, x1, y1, EARTH_RADIUS)
	}
	avg := haversineSum / float64(len(pairs))
	globalProfiler.EndBandwidth("SumHaversine")

	globalProfiler.StartBlock("MiscOutput")
	p.Printf("Count: %*d\nHaversine sum: %.16f\nHaversine avg: %.16f\n", 14, len(pairs), haversineSum, avg)
	globalProfiler.EndBlock("MiscOutput")
}

// Main
//
func main() {
	globalProfiler.BeginProfile()

	// Get input args
	globalProfiler.StartBlock("Startup")
	inputFileArg := flag.String("input", "pairs.json", "Name of input file containing point pairs")
	flag.Parse()
	globalProfiler.EndBlock("Startup")

	// Read JSON file
	data, err := readEntireFile(*inputFileArg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Convert to string
	globalProfiler.StartBlock("ReadToStr")
	strData := string(data)
	globalProfiler.EndBlock("ReadToStr")
	DebugPrintln(strData)

	// Parse
	jsonResult, err := ParseJson(strData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Compute Haversine & print results
	haversineSum(jsonResult)

	globalProfiler.EndAndPrintProfile()
}
