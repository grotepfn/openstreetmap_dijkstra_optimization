package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/qedus/osmpbf"
)

var nodeIdsAllWays [][]int64
var nodesIdsLocations = make(map[int64][2]float64)
var mapNodesBeginsWith = make(map[int64]int)
var polygon [][][2]float64
var nc, wc, rc uint64
var pointInPolygon [][2]float64
var pointInWater GeoPoint
var pointInWater2 GeoPoint

var bitArray [500][1000]bool

var boundingBox [][4]float64

var preRotateBitArray [len(bitArray[0]) * len(bitArray)]float64
var lock = sync.RWMutex{}
var mapPreCalcPoly = make(map[[2]float64]float64)

type GeoPoint struct {
	lat float64
	lng float64
}

func main() {
	println("numcpus " + strconv.Itoa(runtime.NumCPU()))
	runtime.GOMAXPROCS(runtime.NumCPU())

	pointInWater = GeoPoint{90, 0}

	pointInWater2 = GeoPoint{-20.3034175184893, -10.546875}

	t := time.Now()

	println("starting at " + t.String())

	f, err := os.Open("data/planet-coastlines.pbf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)

	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)

	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		log.Fatal(err)
	}

	t = time.Now()

	setup(d)

	println("done with reading" + t.String())

	fillMap()

	t = time.Now()
	println("done with ids" + t.String())
	println("number of ways before merge: " + strconv.Itoa(len(nodeIdsAllWays)))

	mergeSingleWays()

	println("number of ways after merge: " + strconv.Itoa(len(nodeIdsAllWays)))
	t = time.Now()
	println(t.String())

	createPolygon()
	println("created polygon")
	t = time.Now()
	println(t.String())

	fillRotMap()
	fillRotMap2()

	fillBitArray()

	t = time.Now()
	println(t.String())

	jsonString1, _ := json.Marshal(bitArray)
	ioutil.WriteFile("data/bitArray", jsonString1, 0644)

}

func setup(d *osmpbf.Decoder) {
	var nodeIds []int64
	//get all ways ids and fill map with location data
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				// Process Node v.
				nodesIdsLocations[v.ID] = [2]float64{v.Lon, v.Lat}
				nc++

			case *osmpbf.Way:
				// Process Way v.
				wc++

				value, ok := v.Tags["natural"]

				if ok {
					if value == "coastline" {
						nodeIds = []int64{}
						//coordinates = append(coordinates, []float64{v., v.Lat})

						for i := 0; i < len(v.NodeIDs); i++ {
							//println(v.NodeIDs[i])
							nodeIds = append(nodeIds, v.NodeIDs[i])

						}
						//print(value)
						nodeIdsAllWays = append(nodeIdsAllWays, nodeIds)
					}
				}
			case *osmpbf.Relation:
				// Process Relation v.
				rc++
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}

	}

}

func mergeSingleWays() {

inner:
	for i := 0; i <= len(nodeIdsAllWays)-1; i++ {

		if nodeIdsAllWays[i] == nil {
			continue inner
		}

		for nodeIdsAllWays[i][0] != nodeIdsAllWays[i][len(nodeIdsAllWays[i])-1] {
			var help = nodeIdsAllWays[i]

			var helpx = nodeIdsAllWays[i][len(nodeIdsAllWays[i])-1]

			var indexxx = mapNodesBeginsWith[helpx]

			var app = nodeIdsAllWays[indexxx]

			for m := 1; m <= len(app)-1; m++ {
				help = append(help, app[m])

			}

			nodeIdsAllWays[i] = help

			// Remove p[indexxx]
			nodeIdsAllWays[indexxx] = nil // Truncate slice.
		}
	}

	var changed = true
	for changed {
		changed = false
		for i := 0; i <= len(nodeIdsAllWays)-1; i++ {

			if nodeIdsAllWays[i] == nil {

				// Remove p[indexxx]
				nodeIdsAllWays[i] = nodeIdsAllWays[len(nodeIdsAllWays)-1] // Copy last element to index i.
				nodeIdsAllWays[len(nodeIdsAllWays)-1] = nil               // Erase last element (write zero value).
				nodeIdsAllWays = nodeIdsAllWays[:len(nodeIdsAllWays)-1]
				changed = true
			}

		}
	}
}

func createPolygon() {

	for i := 0; i <= len(nodeIdsAllWays)-1; i++ {

		//bounding box
		var smallestLat float64 = 10000
		var biggestLat float64 = -10000
		var smallestLng float64 = 10000
		var biggestLng float64 = -10000
		//

		var helparray [][2]float64
		for j := 0; j <= len(nodeIdsAllWays[i])-1; j++ {
			var helparray2 [2]float64

			helparray2[0] = nodesIdsLocations[nodeIdsAllWays[i][j]][0]
			helparray2[1] = nodesIdsLocations[nodeIdsAllWays[i][j]][1]
			helparray = append(helparray, helparray2)

			if helparray2[0] > biggestLng {

				biggestLng = helparray2[0]
			}

			if helparray2[0] < smallestLng {
				smallestLng = helparray2[0]
			}

			if helparray2[1] > biggestLat {
				biggestLat = helparray2[1]
			}

			if helparray2[1] < smallestLat {
				smallestLat = helparray2[1]
			}
		}

		boundingBox = append(boundingBox, [4]float64{smallestLng, biggestLng, smallestLat, biggestLat})

		polygon = append(polygon, helparray)
	}

}

func isCrossing(X GeoPoint, P GeoPoint, A GeoPoint, B GeoPoint) bool {

	//https://gis.stackexchange.com/questions/10808/manually-transforming-rotated-lat-lon-to-regular-lat-lon

	//A, B and X on the same longitude or P antipodal to X

	var t = false
	if A.lng == B.lng || A.lng == P.lng || B.lng == P.lng || X.lng == A.lng || X.lng == B.lng || X.lng == P.lng || P.lng == -180 || P.lng == -90 || P.lng == 90 {
		//	println("hier")
		X = pointInWater2

		t = true
	}
	X = pointInWater2
	t = true
	//t = true
	if !azimuthMiddle(X, P, A, B, t) {

		return false
	}

	var X2 float64
	X2 = rotateLng(X, A)

	var B2 float64
	B2 = rotateLng(B, A)

	var P2 float64
	P2 = rotateLng(P, A)

	var A2 = 0.0
	var ln = 0.0
	if (P2 - B2) < -180 {
		ln = modLikePython((P2 - B2), 180.0)
	} else if (P2 - B2) > 180 {
		ln = -180 + (P2 - B2 - 180)
	} else {
		ln = P2 - B2
	}

	P2 = ln
	if (A2 - B2) < -180 {
		ln = modLikePython((A2 - B2), 180.0)
	} else if (A2 - B2) > 180 {
		ln = -180 + (A2 - B2 - 180)
	} else {
		ln = A2 - B2
	}

	if (X2 - B2) < -180 {
		ln = modLikePython((X2 - B2), 180.0)
	} else if (X2 - B2) > 180 {
		ln = -180 + (X2 - B2 - 180)
	} else {
		ln = X2 - B2
	}

	X2 = ln

	B2 = 0

	if (P2 >= 0 && X2 <= 0) || (P2 <= 0 && X2 >= 0) {

		return true
	}

	return false
}

func azimuthMiddle(X GeoPoint, P GeoPoint, A GeoPoint, B GeoPoint, c bool) bool {

	var P2 float64
	P2 = P.lng
	if c {
		var da = getArrayPositionFromCords(P.lat, P.lng)
		P2 = preRotateBitArray[da[0]*len(bitArray[0])+da[1]]
		//P2 = rotateLng(P, X)
	}

	var A2 float64
	A2 = A.lng
	if c {
		//A2 = rotateLng(A, X)
		A2 = mapPreCalcPoly[[2]float64{A.lat, A.lng}]
	}

	var B2 float64
	B2 = B.lng
	if c {
		//B2 = rotateLng(B, X)
		B2 = mapPreCalcPoly[[2]float64{B.lat, B.lng}]
	}

	var ln = 0.0
	if (P2 - B2) < -180 {
		ln = modLikePython((P2 - B2), 180.0)
	} else if (P2 - B2) > 180 {
		ln = -180 + (P2 - B2 - 180)
	} else {
		ln = P2 - B2
	}

	P2 = ln

	if (A2 - B2) < -180 {
		ln = modLikePython((A2 - B2), 180.0)
	} else if (A2 - B2) > 180 {
		ln = -180 + (A2 - B2 - 180)
	} else {
		ln = A2 - B2
	}
	A2 = ln

	B2 = 0

	if (P2 >= B2 && P2 <= A2) || (P2 <= B2 && P2 >= A2) {

		return true

	}

	return false

}

func checkPolygon(X GeoPoint, P GeoPoint, B [][2]float64) bool {
	var counter = 0

	for i := 0; i <= len(B)-2; i++ {

		//lat lon geojson format
		if isCrossing(X, P, GeoPoint{B[i][1], B[i][0]}, GeoPoint{B[i+1][1], B[i+1][0]}) {

			counter++
		}
	}
	if counter%2 == 1 {
		return false
	}
	return true
}

func fillMap() {
	for i := 0; i <= len(nodeIdsAllWays)-1; i++ {

		mapNodesBeginsWith[nodeIdsAllWays[i][0]] = i

	}

}

func isPointInWater(point GeoPoint) bool {

	for i := 0; i <= len(polygon)-1; i++ {

		if point.lng >= boundingBox[i][0] && point.lng <= boundingBox[i][1] && point.lat >= boundingBox[i][2] && point.lat <= boundingBox[i][3] { //|| (point.lng <= -160 && boundingBox[i][2] >= 179) {

			if checkPolygon(pointInWater, point, polygon[i]) == false {

				return false
			}
		}

	}

	return true

}

func rotateLng(P GeoPoint, X GeoPoint) float64 {

	var latX = X.lat
	var lonX = X.lng

	latX = -latX

	//south north pole
	if lonX > 0 {
		lonX = -(180 - lonX)
	} else if lonX < 0 {
		lonX = 180 + lonX
	} else if lonX == 0 {
		lonX = 180
	}

	var lon = (P.lng * math.Pi) / 180
	var lat = (P.lat * math.Pi) / 180

	var x = math.Cos(lon) * math.Cos(lat)
	var y = math.Sin(lon) * math.Cos(lat)
	var z = math.Sin(lat)

	var theta = (90 + latX)
	theta = (theta * math.Pi) / 180
	var phi = lonX
	phi = (phi * math.Pi) / 180

	var newX = math.Cos(theta)*math.Cos(phi)*x + math.Cos(theta)*math.Sin(phi)*y + math.Sin(theta)*z

	var newY = -math.Sin(phi)*x + math.Cos(phi)*y

	var newLng = math.Atan2(newY, newX)
	newLng = (newLng * 180) / math.Pi
	return newLng
}

func fillBitArray() {

	var wg sync.WaitGroup
	wg.Add(len(bitArray))

	for i := 0; i <= len(bitArray)-1; i++ {

		go func(i int) {
			defer wg.Done()
			for j := 0; j <= len(bitArray[i])-1; j++ {

				var bla = getCordsFromArrayPosition(i, j)
				var x = GeoPoint{bla[0], bla[1]}
				bitArray[i][j] = isPointInWater(x)

			}

		}(i)
	}
	wg.Wait()

	for i := 0; i <= len(bitArray)-1; i = i + 10 {
		for j := 0; j <= len(bitArray[i])-1; j = j + 10 {

			if bitArray[i][j] {
				print(" ")

			} else {
				print("X")
			}
		}
		println("")
	}
}

//https://stackoverflow.com/questions/43018206/modulo-of-negative-integers-in-go
func modLikePython(d, m float64) float64 {
	var res float64 = math.Mod(d, m)
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return res + m
	}
	return res
}

func fillRotMap() {

	var wg sync.WaitGroup
	wg.Add(len(bitArray))
	for i := 0; i <= len(bitArray)-1; i++ {

		go func(i int) {
			defer wg.Done()
			for j := 0; j <= len(bitArray[0])-1; j++ {

				var x = getCordsFromArrayPosition(i, j)
				var z = rotateLng(GeoPoint{x[0], x[1]}, pointInWater2)
				preRotateBitArray[i*len(bitArray[0])+j] = z
			}
		}(i)

	}
	wg.Wait()
}

func fillRotMap2() {
	var wg sync.WaitGroup
	wg.Add(len(polygon))
	for i := 0; i <= len(polygon)-1; i++ {
		go func(i int) {
			defer wg.Done()
			for j := 0; j <= len(polygon[i])-1; j++ {
				var x = polygon[i][j]
				//CAUTION 1 0
				var z = rotateLng(GeoPoint{x[1], x[0]}, pointInWater2)

				lock.Lock()
				mapPreCalcPoly[[2]float64{x[1], x[0]}] = z
				lock.Unlock()
			}
		}(i)

	}
	wg.Wait()
}

//lat lng
func getCordsFromArrayPosition(pos1 int, pos2 int) [2]float64 {

	return [2]float64{90 - (float64(pos1) / float64(len(bitArray)) * 180), -180 + 360*(float64(pos2)/float64(len(bitArray[0])))}

}

//lat lng
func getArrayPositionFromCords(lat float64, lng float64) [2]int {

	return [2]int{int(math.Round((lat - 90) / 180 * float64(len(bitArray)) * -1)), int(math.Round((lng + 180) / 360 * float64(len(bitArray[0])-1)))}

}
