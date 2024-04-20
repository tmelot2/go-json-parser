package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"strconv"
)


type Point struct {
	X0 float64
	Y0 float64
	X1 float64
	Y1 float64
}


func radiansFromDegrees(degrees float64) float64 {
	result := 0.01745329251994329577 * degrees
	return result
}


func square(a float64) float64 {
	result := a*a
	return result
}


func referenceHaversine(x0, y0, x1, y1, radius float64) float64 {
	lat1 := y0
	lat2 := y1
	lon1 := x0
	lon2 := x1

	dLat := radiansFromDegrees(lat2 - lat1)
	dLon := radiansFromDegrees(lon2 - lon1)
	lat1 = radiansFromDegrees(lat1)
	lat2 = radiansFromDegrees(lat2)

	a := square(math.Sin(dLat/2.0)) + math.Cos(lat1)*math.Cos(lat2)*square(math.Sin(dLon/2))
	c := 2.0 * math.Asin(math.Sqrt(a))

	result := radius * c
	return result
}


func parseJsonFile(path string) []Point {
	// Read JSON from file
	jsonFile, loadFileErr := os.Open(path)
	if loadFileErr != nil {
		fmt.Println("Error:", loadFileErr)
		return nil
	}
	defer jsonFile.Close()

	// Read line by line
	var points []Point
	scanner := bufio.NewScanner(jsonFile)
	lineNum := 0
	for scanner.Scan() {
		// Skip 1st line
		if lineNum == 0 {
			lineNum += 1
			continue
		}

		line := scanner.Text()

		// Skip last line
		if line == "]}" {
			continue
		}

		// Parse line
		x0, y0, x1, y1 := parseLine(line)
		points = append(points, Point{x0,y0,x1,y1})
	}

	return points
}


func parseLine(line string) (float64, float64, float64, float64) {
	substrings := []string{"x0", "y0", "x1", "y1"}
	x0, y0, x1, y1 := 0.0, 0.0, 0.0, 0.0

	for _,substr := range substrings {
		// Find start of point part
		startIndex := strings.Index(line, substr) + 4
		if startIndex < 0 {
			continue
		}

		// Find end of point part (relative to the start)
		endIndex := startIndex + strings.IndexAny(line[startIndex:], " ,}\n")
		if endIndex < startIndex {
			continue
		}

		// Parse as float64
		parsed, err := strconv.ParseFloat(line[startIndex:endIndex], 64)
		if err != nil {
			fmt.Println("Error:", err)
			return 0,0,0,0
		}

		// Assign parts
		switch substr {
		case "x0":
			x0 = parsed
		case "y0":
			y0 = parsed
		case "x1":
			x1 = parsed
		case "y1":
			y1 = parsed
		}
	}

	return x0,y0,x1,y1
}



// Main
//
func main() {
	const EARTH_RADIUS = 6372.8

	// Parse input args
	inputFileArg := flag.String("input", "pairs.json", "Name of input file containing point pairs")
	flag.Parse()

	points := parseJsonFile(*inputFileArg)

	haversineSum := 0.0
	for _,p := range points {
		haversineSum += referenceHaversine(p.X0, p.Y0, p.X1, p.Y1, EARTH_RADIUS)
	}

	avg := haversineSum / float64(len(points))
	fmt.Printf("Count: %d\nHaversine sum: %.16f\nHaversine avg: %.16f\n", len(points), haversineSum, avg)
}
