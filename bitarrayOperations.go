package main

import "math"

//lat lng
func getCordsFromArrayPosition(leny, lenx, pos1, pos2 int) [2]float64 {

	return [2]float64{90 - (float64(pos1) / float64(leny) * 180), -180 + 360*(float64(pos2)/float64(lenx))}

}

//lat lng
func getArrayPositionFromCords(leny, lenx int, lat, lng float64) [2]int {

	return [2]int{int(math.Round((lat - 90) / 180 * float64(leny) * -1)), int(math.Round((lng + 180) / 360 * float64(lenx-1)))}

}

//https://stackoverflow.com/questions/43018206/modulo-of-negative-integers-in-go
func modLikePython(d, m int) int {
	var res int = d % m
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return res + m
	}
	return res
}

//https://stackoverflow.com/questions/43018206/modulo-of-negative-integers-in-go
func modLikePythonFloat(d, m float64) float64 {
	var res float64 = math.Mod(d, m)
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return res + m
	}
	return res
}

type GeoPoint struct {
	lat float64
	lng float64
}

//https://github.com/kellydunn/golang-geo/blob/master/point.go
// GreatCircleDistance: Calculates the Haversine distance between two points in kilometers.
// Original Implementation from: http://www.movable-type.co.uk/scripts/latlong.html
func GreatCircleDistance(l1 [2]float64, l2 [2]float64) float64 {
	var EARTH_RADIUS = 6371.0
	dLat := (l2[0] - l1[0]) * (math.Pi / 180.0)
	dLon := (l2[1] - l1[1]) * (math.Pi / 180.0)

	lat1 := l1[0] * (math.Pi / 180.0)
	lat2 := l2[0] * (math.Pi / 180.0)

	a1 := math.Sin(dLat/2) * math.Sin(dLat/2)
	a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)

	a := a1 + a2

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EARTH_RADIUS * c
}
