/*
	Generates Haversine points JSON data file. Also computes Haversine & prints reference answer.
*/

package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
)


func createJsonItem(x0, y0, x1, y1 float64) string {
	return fmt.Sprintf("{ \"x0\":%.16f, \"y0\":%.16f, \"x1\":%.16f, \"y1\":%.16f }", x0, y0, x1, y1)
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

func getRandomPoint(centerX, centerY float64) (float64, float64) {
	x := centerX + X_MIN + (X_MAX - X_MIN) * rand.Float64()
	y := centerY + Y_MIN + (Y_MAX - Y_MIN) * rand.Float64()

	// Cap x
	if x < X_MIN {
		x = X_MIN
	} else if x > X_MAX {
		x = X_MAX
	}
	// Cap y
	if y < Y_MIN {
		y = Y_MIN
	} else if y > Y_MAX {
		y = Y_MAX
	}

	return x,y
}


const EARTH_RADIUS = 6372.8
// x: -180 to 180
const X_MIN = -180
const X_MAX = 180
// y: -90 to 90
const Y_MIN = -90
const Y_MAX = 90
const CLUSTER_SIZE = 64 // TODO: Make const or input arg

func main() {
	// Parse input args
	pairsArg := flag.Int("pairs", 10, "Number of pairs of points to generate")
	methodArg := flag.String("method", "uniform", "Point distribution method: uniform or cluster")
	flag.Parse()
	pairs := *pairsArg

	// Setup
	clusterCountLeft := 0
	haversineSum := 0.0
	caseyHaversineSum := 0.0
	centerX := 0.0
	centerY := 0.0

	// Validate method
	if *methodArg == "cluster" {
		clusterCountLeft = CLUSTER_SIZE
	} else if *methodArg != "uniform" {
		*methodArg = "uniform"
	}

	// Open file
	file, err := os.Create("pairs.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "{\"pairs\":[\n")

	for i := 0; i < pairs; i++ {
		clusterCountLeft -= 1
		if clusterCountLeft == 0 {
			clusterCountLeft = CLUSTER_SIZE
			centerX, centerY = getRandomPoint(0,0)
		}

		x0,y0 := getRandomPoint(centerX, centerY)
		x1,y1 := getRandomPoint(centerX, centerY)
		newItem := createJsonItem(x0, y0, x1, y1)
		haversineDistance := referenceHaversine(x0, y0, x1, y1, EARTH_RADIUS)
		haversineSum += haversineDistance
		caseyHaversineSum += (1.0/float64(pairs)) * haversineDistance
		// Comma
		if i < pairs-1 {
			newItem = fmt.Sprintf("%s,", newItem)
		}
		fmt.Fprintf(file, "\t%s\n", newItem)
	}

	avg := haversineSum/float64(pairs)
	fmt.Fprintf(file, "]}")
	fmt.Printf("Count: %d\n  Sum: %.16f\n  Avg: %.16f\n CSum: %.16f\n", pairs, haversineSum, avg, caseyHaversineSum)
	fmt.Printf(" Diff: %.16f\n", math.Abs(avg-caseyHaversineSum))
}

// Ref:
// {"x0":102.1633205722960440, "y0":-24.9977499718717624, "x1":-14.3322557404258362, "y1":62.6708294856625940},
