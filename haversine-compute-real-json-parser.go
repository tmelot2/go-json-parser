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

var DEBUG = false


// Main
//
func main() {
	const EARTH_RADIUS = 6372.8

	// Get input args
	inputFileArg := flag.String("input", "pairs.json", "Name of input file containing point pairs")
	flag.Parse()

	// Read JSON file, convert to string
	data, err := os.ReadFile(*inputFileArg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	strData := string(data)
	DebugPrintln(strData)

	// TODO: Combine lexer & parser into one function call

	// Lex JSON
	lexer := newLexer(strData)
	lexer.Debug = DEBUG
	lexedTokens, err := lexer.lex()
	if err != nil {
		fmt.Printf("Lexer error: %s\n", err)
	}

	// Parse JSON
	parser := newParser(lexedTokens)
	jsonResult, err := parser.Parse()
	if err != nil {
		fmt.Println("Parser error:", err)
	}
	// fmt.Println("Parser result:", jsonResult)

	// Loop over JSON to do stuff
	// TODO: Figure out how to abstract casting stuff into separate client logic
	fmt.Println("===============================")
	points, ok := jsonResult["pairs"].([]any)
	if !ok {
		fmt.Println("Error casting pairs array")
		return
	}

	haversineSum := 0.0
	for _, p := range points {
		point, ok2 := p.(map[string]any)
		if !ok2 {
			fmt.Println("Error casting point to map")
			continue
		}
		x0 := point["x0"].(float64)
		y0 := point["y0"].(float64)
		x1 := point["x1"].(float64)
		y1 := point["y1"].(float64)
		haversineSum += referenceHaversine(x0, y0, x1, y1, EARTH_RADIUS)
	}
	avg := haversineSum / float64(len(points))
	fmt.Printf("Count: %d\nHaversine sum: %.16f\nHaversine avg: %.16f\n", len(points), haversineSum, avg)
}
