package main

import (
	"flag"
	"fmt"
	"os"
	// "strings"
	// "strconv"
)


type Point struct {
	X0 float64
	Y0 float64
	X1 float64
	Y1 float64
}

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

	// Parse input args
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

	// Lex JSON tokens
	lexer := newLexer(strData)
	lexer.Debug = DEBUG
	lexedTokens, err := lexer.lex()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Println("=========================================\nLEXED TOKENS (partial, see above error)\n=========================================")
	} else {
		fmt.Println("\n=====================\nLEXED TOKENS\n=====================")
	}
	for _,t := range lexedTokens {
		fmt.Printf("(\"%s\", %s), \n", t.Value, t.Type)
	}
	fmt.Printf("\nToken count: %d\n", len(lexedTokens))

	fmt.Println("")

	fmt.Println("Parser stuff")
	parser := newParser(lexedTokens)
	result, err := parser.Parse()
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Parser result:", result)

	// Haversine average for fake data
	var points []Point
	points = append(points, Point{X0:1.1, Y0:2.2, X1:3.3, Y1:4.4})
	points = append(points, Point{X0:11.11, Y0:22.22, X1:33.33, Y1:44.44})
	// Compute Haversines & average
	haversineSum := 0.0
	for _,p := range points {
		haversineSum += referenceHaversine(p.X0, p.Y0, p.X1, p.Y1, EARTH_RADIUS)
	}

	avg := haversineSum / float64(len(points))
	fmt.Printf("Count: %d\nHaversine sum: %.16f\nHaversine avg: %.16f\n", len(points), haversineSum, avg)
}
