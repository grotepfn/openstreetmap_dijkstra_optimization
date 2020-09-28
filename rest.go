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
)

var result [][]bool
var optEdges [][2][2]int
var distsOpt [][][]float64

var mapPointSquares = make(map[[2]int][]int)

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

	var desPos = getArrayPositionFromCords(result, laDes, lnDes)

	var startPos = getArrayPositionFromCords(result, la, ln)

	var pre, distances, _, pops = dijkstra_single_optimized(startPos[0], startPos[1], desPos[0], desPos[1], mapPointSquares, optEdges, distsOpt)
	var pre2, distances3, _, pops3 = a_stern_single_optimized(startPos[0], startPos[1], desPos[0], desPos[1], mapPointSquares, optEdges)
	var _, distances4, _, pops4 = dijkstra_single_optimized_old(startPos[0], startPos[1], desPos[0], desPos[1], mapPointSquares, optEdges)
	var _, distances5, _, pops5 = a_stern_single_optimized_with_pre(startPos[0], startPos[1], desPos[0], desPos[1], mapPointSquares, optEdges, distsOpt)

	var _, distances2, _, pops2 = dijkstra_single(startPos[0], startPos[1], desPos[0], desPos[1])

	println("opti")
	println(distances[desPos[0]*len(result[0])+desPos[1]])
	println(pops)

	println("opti")
	println(distances3[desPos[0]*len(result[0])+desPos[1]])
	println(pops3)

	println("opti")
	println(distances4[desPos[0]*len(result[0])+desPos[1]])
	println(pops4)

	println("opti")
	println(distances5[desPos[0]*len(result[0])+desPos[1]])
	println(pops5)

	println("slow")
	println(distances2[desPos[0]*len(result[0])+desPos[1]])
	println(pops2)

	var way [][2]int
	if algorithm == "dijkstra" {
		println("dijkstra")
		way = getShortestPath2(desPos[0], desPos[1], pre)
	} else if algorithm == "astar" {
		println("astern")
		way = getShortestPath2(desPos[0], desPos[1], pre2)
	}

	var divided = true
	for divided {
		divided = false
		for i := 0; i <= len(way)-2; i = i + 2 {

			if math.Abs(float64(way[i][0]-way[i+1][0])) >= 2 || math.Abs(float64(way[i][1]-way[i+1][1])) >= 2 {
				divided = true

				var midPoint = getMidPoint(way[i][0], way[i][1], way[i+1][0], way[i+1][1])

				way = insert(way, i+1, midPoint)

			}

		}
	}
	//println(len(way))

	var wayCords [][2]float64
	for i := 0; i <= len(way)-1; i++ {

		wayCords = append(wayCords, getCordsFromArrayPosition(result, way[i][0], way[i][1]))
	}

	var payload, err = json.Marshal(wayCords)

	if err != nil {
		log.Println(err)
		log.Println(distances)

		log.Println(pre)

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

				var list = mapPointSquares[[2]int{optEdges[i][0][0] + j, optEdges[i][0][1] + k}]
				list = append(list, i)
				mapPointSquares[[2]int{optEdges[i][0][0] + j, optEdges[i][0][1] + k}] = list

			}
		}

	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	log.Fatal(http.ListenAndServe(":8080", router))

}

//https://de.mathworks.com/matlabcentral/answers/229312-how-to-calculate-the-middle-point-between-two-points-on-the-earth-in-matlab
func getMidPoint(lat1, lng1, lat2, lng2 int) [2]int {

	var l = getCordsFromArrayPosition(result, lat1, lng1)
	var la = l[0]
	var ln = l[1]
	l = getCordsFromArrayPosition(result, lat2, lng2)
	var la2 = l[0]
	var ln2 = l[1]

	var Bx = math.Cos(la2*(math.Pi/180.0)) * math.Cos((ln2-ln)*(math.Pi/180.0))
	var By = math.Cos(la2*(math.Pi/180.0)) * math.Sin((ln2-ln)*(math.Pi/180.0))

	var latMid = (180 / math.Pi) * math.Atan2(math.Sin(la*(math.Pi/180.0))+math.Sin(la2*(math.Pi/180.0)), math.Sqrt((math.Cos(la*(math.Pi/180.0))+Bx)*(math.Cos(la*(math.Pi/180.0))+Bx)+By*By))

	var lonMid = ln + (180/math.Pi)*math.Atan2(By, math.Cos(la*(math.Pi/180.0))+Bx)

	return getArrayPositionFromCords(result, latMid, lonMid)

}

//https://stackoverflow.com/questions/46128016/insert-a-value-in-a-slice-at-a-given-index
// 0 <= index <= len(a)
func insert(a [][2]int, index int, value [2]int) [][2]int {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}
