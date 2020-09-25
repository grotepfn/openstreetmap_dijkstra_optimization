package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

	var _, distances, _, pops = dijkstra_single_optimized(startPos[0], startPos[1], desPos[0], desPos[1], mapPointSquares, optEdges, distsOpt)
	var _, distances3, _, pops3 = a_stern_single_optimized(startPos[0], startPos[1], desPos[0], desPos[1], mapPointSquares, optEdges)
	var _, distances4, _, pops4 = dijkstra_single_optimized_old(startPos[0], startPos[1], desPos[0], desPos[1], mapPointSquares, optEdges)
	var pre, distances5, _, pops5 = a_stern_single_optimized_with_pre(startPos[0], startPos[1], desPos[0], desPos[1], mapPointSquares, optEdges, distsOpt)

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

	var way = getShortestPath2(desPos[0], desPos[1], pre)
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

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	log.Fatal(http.ListenAndServe(":8080", router))

}
