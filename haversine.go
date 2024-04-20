package main

import (
	"math"
)


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
