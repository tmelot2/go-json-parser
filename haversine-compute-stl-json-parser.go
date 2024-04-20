package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
)

type Point struct {
	X0 float64 `json:"x0"`
	Y0 float64 `json:"y0"`
	X1 float64 `json:"x1"`
	Y1 float64 `json:"y1"`
}

type Data struct {
	Points []Point `json:"pairs"`
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


// Main
//
func main() {
	const EARTH_RADIUS = 6372.8

	// Parse input args
	inputFileArg := flag.String("input", "pairs.json", "Name of input file containing point pairs")
	flag.Parse()

	// Read JSON from file
	jsonData, loadFileErr := ioutil.ReadFile(*inputFileArg)
	if loadFileErr != nil {
		fmt.Println("Error:", loadFileErr)
		return
	}

	// Parse JSON
	var data Data
	parseJsonErr := json.Unmarshal(jsonData, &data)
	if parseJsonErr != nil {
		fmt.Println("Error:", parseJsonErr)
		return
	}

	// Loop over points to compute Haversine sum
	haversineSum := 0.0
	for _,p := range data.Points {
		haversineSum += referenceHaversine(p.X0, p.Y0, p.X1, p.Y1, EARTH_RADIUS)
		// fmt.Printf("%d (%4.16f, %3.16f) (%4.16f, %3.16f)\n", i+1, p.X0, p.Y0, p.X1, p.Y1)
	}

	// Compute final average & print answer
	avg := haversineSum / float64(len(data.Points))
	fmt.Printf("Count: %d\nHaversine sum: %.16f\nHaversine avg: %.16f\n", len(data.Points), haversineSum, avg)
}
