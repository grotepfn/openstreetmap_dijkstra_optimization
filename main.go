package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/grotepfn/openstreetmap_dijkstra_optimization/bitArray"
)

var result [][]bool
var optEdges [][2][2]int
var distsOpt [][][]float64

var mapPointSquares = make(map[[2]int]int)

func homeLink(w http.ResponseWriter, r *http.Request) {

	var algorithm = r.URL.Query().Get("algorithm")

	var lat = r.URL.Query().Get("lat")
	var lng = r.URL.Query().Get("lng")
	println(lat)
	println(lng)

	var latDes = r.URL.Query().Get("latDes")
	var lngDes = r.URL.Query().Get("lngDes")

	var la, error = strconv.ParseFloat(lat, 64)
	var ln, error2 = strconv.ParseFloat(lng, 64)
	var laDes, error3 = strconv.ParseFloat(latDes, 64)
	var lnDes, error4 = strconv.ParseFloat(lngDes, 64)

	if error != nil {
		log.Println(error)
	}
	if error2 != nil {
		log.Println(error2)
	}
	if error3 != nil {
		log.Println(error3)
	}
	if error4 != nil {
		log.Println(error4)
	}

	var desPos = bitArray.GetArrayPositionFromCords(len(result), len(result[0]), laDes, lnDes)

	var startPos = bitArray.GetArrayPositionFromCords(len(result), len(result[0]), la, ln)

	var preDij, distancesDij, _, popsDij = bitArray.Dijkstra_single(startPos[0], startPos[1], desPos[0], desPos[1], result)

	var preDijBi, distancesDijBi, preDijBi2, distancesDijBi2, bestMeetingPoint, popsDijBi = bitArray.Dijkstra_bi(result, startPos[0], startPos[1], desPos[0], desPos[1])
	var preDijOpt, distancesDijOpt, _, popsDijOpt = bitArray.Dijkstra_single_optimized(result, startPos[0], startPos[1], desPos[0], desPos[1], mapPointSquares, optEdges)
	var preDijOptPre, distancesDijOptPre, _, popsDijOptPre = bitArray.Dijkstra_single_optimized_pre(result, startPos[0], startPos[1], desPos[0], desPos[1], mapPointSquares, optEdges, distsOpt)
	var preAStar, distancesAStar, _, popsAStar = bitArray.A_stern_single(result, startPos[0], startPos[1], desPos[0], desPos[1])
	var preAStarOpt, distancesAStarOpt, _, popsAStarOpt = bitArray.A_stern_single_optimized(result, startPos[0], startPos[1], desPos[0], desPos[1], mapPointSquares, optEdges)
	var preAStarOptPre, distancesAStarOptPre, _, popsAStarOptPre = bitArray.A_stern_single_optimized_with_pre(result, startPos[0], startPos[1], desPos[0], desPos[1], mapPointSquares, optEdges, distsOpt)

	println("ergebnisse")
	println("dijkstra single")
	println(distancesDij[desPos[0]*len(result[0])+desPos[1]])
	println(popsDij)

	println("dijkstra bi")
	println(distancesDijBi[bestMeetingPoint[0]*len(result[0])+bestMeetingPoint[1]] + distancesDijBi2[bestMeetingPoint[0]*len(result[0])+bestMeetingPoint[1]])
	println(popsDijBi)

	println("dijkstra single optimized")
	println(distancesDijOpt[desPos[0]*len(result[0])+desPos[1]])
	println(popsDijOpt)

	println("dijkstra single optimized pre")
	println(distancesDijOptPre[desPos[0]*len(result[0])+desPos[1]])
	println(popsDijOptPre)

	println("a star single")
	println(distancesAStar[desPos[0]*len(result[0])+desPos[1]])
	println(popsAStar)

	println("a star single opt without pre")
	println(distancesAStarOpt[desPos[0]*len(result[0])+desPos[1]])
	println(popsAStarOpt)

	println("a star single opt with pre")
	println(distancesAStarOptPre[desPos[0]*len(result[0])+desPos[1]])
	println(popsAStarOptPre)

	var way [][2]int
	if algorithm == "dijkstra" {
		println("dijkstra")
		way = bitArray.GetShortestPath(result, desPos[0], desPos[1], preDij)
		println(distancesDij[desPos[0]*len(result[0])+desPos[1]])
	} else if algorithm == "astar" {
		println("astern")
		way = bitArray.GetShortestPath(result, desPos[0], desPos[1], preAStar)
		println(distancesAStar[desPos[0]*len(result[0])+desPos[1]])
	} else if algorithm == "dijkstraOpt" {
		println("dijkstra opt")
		way = bitArray.GetShortestPath(result, desPos[0], desPos[1], preDijOpt)
		println(distancesDijOpt[desPos[0]*len(result[0])+desPos[1]])
	} else if algorithm == "dijkstraOptWithPre" {
		println("dijkstra opt with pre")
		way = bitArray.GetShortestPath(result, desPos[0], desPos[1], preDijOptPre)
		println(distancesDijOptPre[desPos[0]*len(result[0])+desPos[1]])
	} else if algorithm == "astarOpt" {
		println("astar opt")
		way = bitArray.GetShortestPath(result, desPos[0], desPos[1], preAStarOpt)
		println(distancesAStarOpt[desPos[0]*len(result[0])+desPos[1]])
	} else if algorithm == "astarOptWithPre" {
		println("astar opt with pre")
		way = bitArray.GetShortestPath(result, desPos[0], desPos[1], preAStarOptPre)
		println(distancesAStarOptPre[desPos[0]*len(result[0])+desPos[1]])
	} else if algorithm == "biDijkstra" {
		println("dijkstra bi")
		var wayCords [][2]float64

		way = bitArray.GetShortestPath(result, bestMeetingPoint[0], bestMeetingPoint[1], preDijBi)

		for i := 0; i <= len(way)-1; i++ {

			wayCords = append(wayCords, bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), way[i][0], way[i][1]))
		}

		var way2 = bitArray.GetShortestPath(result, bestMeetingPoint[0], bestMeetingPoint[1], preDijBi2)

		//https://stackoverflow.com/questions/19239449/how-do-i-reverse-an-array-in-go
		for i, j := 0, len(way2)-1; i < j; i, j = i+1, j-1 {
			way2[i], way2[j] = way2[j], way2[i]
		}

		way = append(way, way2...)

	}

	println("found a way")

	var wayCords [][2]float64
	for i := 0; i <= len(way)-1; i++ {

		wayCords = append(wayCords, bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), way[i][0], way[i][1]))
	}

	var payload, err = json.Marshal(wayCords)

	if err != nil {
		log.Println(err)

	}

	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)

}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {

	//import files for routing
	jsonFile, err := os.Open("data/bitArray")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &result)

	jsonFile, err = os.Open("data/optimization_squares")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ = ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &optEdges)

	jsonFile, err = os.Open("data/optimization_squares_distances")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ = ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &distsOpt)

	for i := 0; i <= len(optEdges)-1; i++ {

		for j := 0; j <= optEdges[i][1][0]-optEdges[i][0][0]; j++ {

			for k := 0; k <= optEdges[i][1][1]-optEdges[i][0][1]; k++ {

				mapPointSquares[[2]int{optEdges[i][0][0] + j, optEdges[i][0][1] + k}] = i

			}
		}

	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	log.Fatal(http.ListenAndServe(":8080", router))

}

//https://de.mathworks.com/matlabcentral/answers/229312-how-to-calculate-the-middle-point-between-two-points-on-the-earth-in-matlab
func getMidPoint(lat1, lng1, lat2, lng2 float64) [2]float64 {

	var Bx = math.Cos(lat2*(math.Pi/180.0)) * math.Cos((lng2-lng1)*(math.Pi/180.0))
	var By = math.Cos(lat2*(math.Pi/180.0)) * math.Sin((lng2-lng1)*(math.Pi/180.0))

	var latMid = (180 / math.Pi) * math.Atan2(math.Sin(lat1*(math.Pi/180.0))+math.Sin(lat2*(math.Pi/180.0)), math.Sqrt((math.Cos(lat1*(math.Pi/180.0))+Bx)*(math.Cos(lat1*(math.Pi/180.0))+Bx)+By*By))

	var lonMid = lng1 + (180/math.Pi)*math.Atan2(By, math.Cos(lat1*(math.Pi/180.0))+Bx)

	return [2]float64{latMid, lonMid}

}

//https://stackoverflow.com/questions/46128016/insert-a-value-in-a-slice-at-a-given-index
// 0 <= index <= len(a)
func insert(a [][2]float64, index int, value [2]float64) [][2]float64 {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}
