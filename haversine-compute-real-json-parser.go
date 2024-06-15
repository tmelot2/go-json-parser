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


// Main
//
func main() {
	prof := newProfiler()
	prof.Start("Startup")

	// Setup
	const EARTH_RADIUS = 6372.8
	p := GetPrinter()

	// Get input args
	inputFileArg := flag.String("input", "pairs.json", "Name of input file containing point pairs")
	flag.Parse()
	prof.End("Startup")


	// Read JSON file, convert to string
	prof.Start("Read")
	data, err := os.ReadFile(*inputFileArg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	prof.End("Read")
	prof.Start("ReadToStr")
	strData := string(data)
	prof.End("ReadToStr")
	DebugPrintln(strData)

	// Parse
	prof.Start("Parse")
	jsonResult, err := ParseJson(strData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	prof.End("Parse")

	// Loop over JSON to do stuff
	prof.Start("Sum")
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
	prof.End("Sum")

	prof.Start("MiscOutput")
	p.Printf("Count: %*d\nHaversine sum: %.16f\nHaversine avg: %.16f\n", 14, len(pairs), haversineSum, avg)
	prof.End("MiscOutput")

	prof.Print()
}
