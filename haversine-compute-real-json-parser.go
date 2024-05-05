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

	// Parse
	jsonResult, err := ParseJson(strData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Loop over JSON to do stuff
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
	fmt.Printf("Count: %d\nHaversine sum: %.16f\nHaversine avg: %.16f\n", len(pairs), haversineSum, avg)

	// JsonValue tests
	//
	// JSON i was testing stuff with
	// {
	// 	"a": "b",
	// 	"b": 1,
	// 	"c": 2.222,
	// 	"d": {
	// 		"d1": "d",
	// 		"d2": 1,
	// 		"d3": 2.222
	// 	},
	// 	"e": [
	// 		[1,2]
	// 	]
	// }
	// j := NewJsonValue(jsonResult)
	// s, sErr := j.GetString("a")
	// fmt.Println(s, sErr)

	// i, iErr := j.GetInt("b")
	// fmt.Println(i, iErr)

	// f, fErr := j.GetFloat("c")
	// fmt.Println(f, fErr)

	// o, oErr := j.GetObject("d")
	// fmt.Println(o, oErr)
	// fmt.Println(o.GetString("d1"))
	// fmt.Println(o.GetInt("d2"))
	// fmt.Println(o.GetFloat("d3"))

	// a, _ := j.GetArray("e")
	// for _, aVal := range a {
	// 	fmt.Println(aVal)
	// 	fmt.Println(aVal.GetObject(""))
	// 	z, _ := aVal.GetArray("")
	// 	for _, zVal := range z {
	// 		fmt.Println(zVal.GetInt(""))
	// 	}
	// }
}
