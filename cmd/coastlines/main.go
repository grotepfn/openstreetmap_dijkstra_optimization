package main

import (
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/grotepfn/openstreetmap_dijkstra_optimization/bitArray"
	"github.com/qedus/osmpbf"
)

var nodeIdsAllWays [][]int64
var nodesIdsLocations = make(map[int64][2]float64)
var mapNodesBeginsWith = make(map[int64]int)
var polygon [][][2]float64
var nc, wc, rc uint64
var pointInPolygon [][2]float64
var pointInWater bitArray.GeoPoint
var pointInWater2 bitArray.GeoPoint

var result [20][20]bool

var boundingBox [][4]float64

var preRotateresult [len(result[0]) * len(result)]float64
var lock = sync.RWMutex{}
var mapPreCalcPoly = make(map[[2]float64]float64)

func main() {
	println("numcpus " + strconv.Itoa(runtime.NumCPU()))
	runtime.GOMAXPROCS(runtime.NumCPU())

	pointInWater = bitArray.GeoPoint{90, 0}

	pointInWater2 = bitArray.GeoPoint{-20.3034175184893, -10.546875}

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

	fillresult()

	t = time.Now()
	println(t.String())

	//jsonString1, _ := json.Marshal(result)
	//ioutil.WriteFile("data/result", jsonString1, 0644)

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

func isCrossing(X bitArray.GeoPoint, P bitArray.GeoPoint, A bitArray.GeoPoint, B bitArray.GeoPoint) bool {

	//https://gis.stackexchange.com/questions/10808/manually-transforming-rotated-lat-lon-to-regular-lat-lon

	//A, B and X on the same longitude or P antipodal to X

	var t = false
	if A.Lng == B.Lng || A.Lng == P.Lng || B.Lng == P.Lng || X.Lng == A.Lng || X.Lng == B.Lng || X.Lng == P.Lng || P.Lng == -180 || P.Lng == -90 || P.Lng == 90 {
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
		ln = bitArray.ModLikePythonFloat((P2 - B2), 180.0)
	} else if (P2 - B2) > 180 {
		ln = -180 + (P2 - B2 - 180)
	} else {
		ln = P2 - B2
	}

	P2 = ln
	if (A2 - B2) < -180 {
		ln = bitArray.ModLikePythonFloat((A2 - B2), 180.0)
	} else if (A2 - B2) > 180 {
		ln = -180 + (A2 - B2 - 180)
	} else {
		ln = A2 - B2
	}

	if (X2 - B2) < -180 {
		ln = bitArray.ModLikePythonFloat((X2 - B2), 180.0)
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

func azimuthMiddle(X bitArray.GeoPoint, P bitArray.GeoPoint, A bitArray.GeoPoint, B bitArray.GeoPoint, c bool) bool {

	var P2 float64
	P2 = P.Lng
	if c {

		var da = bitArray.GetArrayPositionFromCords(len(result), len(result[0]), P.Lat, P.Lng)
		P2 = preRotateresult[da[0]*len(result[0])+da[1]]
		//P2 = rotateLng(P, X)
	}

	var A2 float64
	A2 = A.Lng
	if c {
		//A2 = rotateLng(A, X)
		A2 = mapPreCalcPoly[[2]float64{A.Lat, A.Lng}]
	}

	var B2 float64
	B2 = B.Lng
	if c {
		//B2 = rotateLng(B, X)
		B2 = mapPreCalcPoly[[2]float64{B.Lat, B.Lng}]
	}

	var ln = 0.0
	if (P2 - B2) < -180 {
		ln = bitArray.ModLikePythonFloat((P2 - B2), 180.0)
	} else if (P2 - B2) > 180 {
		ln = -180 + (P2 - B2 - 180)
	} else {
		ln = P2 - B2
	}

	P2 = ln

	if (A2 - B2) < -180 {
		ln = bitArray.ModLikePythonFloat((A2 - B2), 180.0)
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

func checkPolygon(X bitArray.GeoPoint, P bitArray.GeoPoint, B [][2]float64) bool {
	var counter = 0

	for i := 0; i <= len(B)-2; i++ {

		//lat lon geojson format
		if isCrossing(X, P, bitArray.GeoPoint{B[i][1], B[i][0]}, bitArray.GeoPoint{B[i+1][1], B[i+1][0]}) {

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

func isPointInWater(point bitArray.GeoPoint) bool {

	for i := 0; i <= len(polygon)-1; i++ {

		if point.Lng >= boundingBox[i][0] && point.Lng <= boundingBox[i][1] && point.Lat >= boundingBox[i][2] && point.Lat <= boundingBox[i][3] { //|| (point.lng <= -160 && boundingBox[i][2] >= 179) {

			if checkPolygon(pointInWater, point, polygon[i]) == false {

				return false
			}
		}

	}

	return true

}

func rotateLng(P bitArray.GeoPoint, X bitArray.GeoPoint) float64 {

	var latX = X.Lat
	var lonX = X.Lng

	latX = -latX

	//south north pole
	if lonX > 0 {
		lonX = -(180 - lonX)
	} else if lonX < 0 {
		lonX = 180 + lonX
	} else if lonX == 0 {
		lonX = 180
	}

	var lon = (P.Lng * math.Pi) / 180
	var lat = (P.Lat * math.Pi) / 180

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

func fillresult() {

	var wg sync.WaitGroup
	wg.Add(len(result))

	for i := 0; i <= len(result)-1; i++ {

		go func(i int) {
			defer wg.Done()
			for j := 0; j <= len(result[i])-1; j++ {

				var bla = bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), i, j)
				var x = bitArray.GeoPoint{bla[0], bla[1]}
				result[i][j] = isPointInWater(x)

			}

		}(i)
	}
	wg.Wait()

	for i := 0; i <= len(result)-1; i = i + 1 {
		for j := 0; j <= len(result[i])-1; j = j + 1 {

			if result[i][j] {
				print(" ")

			} else {
				print("X")
			}
		}
		println("")
	}
}

func fillRotMap() {

	var wg sync.WaitGroup
	wg.Add(len(result))
	for i := 0; i <= len(result)-1; i++ {

		go func(i int) {
			defer wg.Done()
			for j := 0; j <= len(result[0])-1; j++ {

				var x = bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), i, j)
				var z = rotateLng(bitArray.GeoPoint{x[0], x[1]}, pointInWater2)
				preRotateresult[i*len(result[0])+j] = z
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
				var z = rotateLng(bitArray.GeoPoint{x[1], x[0]}, pointInWater2)

				lock.Lock()
				mapPreCalcPoly[[2]float64{x[1], x[0]}] = z
				lock.Unlock()
			}
		}(i)

	}
	wg.Wait()
}
