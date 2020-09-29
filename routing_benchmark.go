package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime/pprof"
	"strconv"
	"time"
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

	var desPos = getArrayPositionFromCords(len(result), len(result[0]), laDes, lnDes)

	var startPos = getArrayPositionFromCords(len(result), len(result[0]), la, ln)

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
	/*
		var pre, distances, pre2, _, u, _ = a_stern_bi(startPos[0], startPos[1], desPos[0], desPos[1], mapPointSquares)

		var wayCords [][2]float64

		println("despos")
		println(desPos[0], desPos[1])

		var way = getShortestPath2(u[0], u[1], pre)

		for i := 0; i <= len(way)-1; i++ {

			wayCords = append(wayCords, getCordsFromArrayPosition(result, way[i][0], way[i][1]))
		}

		var way2 = getShortestPath2(u[0], u[1], pre2)

		var wayCords2 [][2]float64
		for i := 0; i <= len(way2)-1; i++ {

			wayCords2 = append(wayCords2, getCordsFromArrayPosition(result, way2[i][0], way2[i][1]))
		}

		//https://stackoverflow.com/questions/19239449/how-do-i-reverse-an-array-in-go
		for i, j := 0, len(wayCords2)-1; i < j; i, j = i+1, j-1 {
			wayCords2[i], wayCords2[j] = wayCords2[j], wayCords2[i]
		}

		//wayCords2 = wayCords2[1:]

		wayCords = append(wayCords, wayCords2...)
	*/
	var way = getShortestPath2(desPos[0], desPos[1], pre)
	var wayCords [][2]float64
	for i := 0; i <= len(way)-1; i++ {

		wayCords = append(wayCords, getCordsFromArrayPosition(len(result), len(result[0]), way[i][0], way[i][1]))
	}

	var payload, err = json.Marshal(wayCords)

	if err != nil {
		log.Println(err)
		log.Println(distances)

		log.Println(pre)

		//log.Println(u2)
		//log.Println(distances2)

		//log.Println(pre2)
	}

	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)

}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	jsonFile, err := os.Open("data/bitarray")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &result)

	var rnd [][4]int

	rand.Seed(time.Now().UnixNano())

	for i := 0; i <= 1000; i++ {

		var y1 = rand.Intn(500)

		var y2 = rand.Intn(500)

		var v1 = rand.Intn(1000)

		var v2 = rand.Intn(1000)
		for result[y1][v1] == false || result[y2][v2] == false {
			y1 = rand.Intn(500)
			y2 = rand.Intn(500)
			v1 = rand.Intn(1000)
			v2 = rand.Intn(1000)
		}

		rnd = append(rnd, [4]int{y1, v1, y2, v2})

	}

	println("single")
	t := time.Now()
	println(t.String())
	var counterPOPS = 0

	for i := 0; i <= 100; i++ {

		//	var _, _, _, counter = dijkstra_single(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3])
		//	counterPOPS = counterPOPS + counter

	}
	println("dijkstra single pops")
	println(counterPOPS)
	t = time.Now()
	println(t.String())

	jsonFile, err = os.Open("data/optimization_squares")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ = ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &optEdges)

	for i := 0; i <= len(optEdges)-1; i++ {

		for j := 0; j <= optEdges[i][1][0]-optEdges[i][0][0]; j++ {

			for k := 0; k <= optEdges[i][1][1]-optEdges[i][0][1]; k++ {

				var list = mapPointSquares[[2]int{optEdges[i][0][0] + j, optEdges[i][0][1] + k}]
				list = append(list, i)
				mapPointSquares[[2]int{optEdges[i][0][0] + j, optEdges[i][0][1] + k}] = list

			}
		}

	}

	jsonFile, err = os.Open("data/optimization_squares_distances")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ = ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &distsOpt)

	println("dijksta opt rnd")
	t = time.Now()
	println(t.String())
	counterPOPS = 0

	for i := 0; i <= -1; i++ {
		//println("sdfsdfs")
		var _, _, _, counter = dijkstra_single_optimized(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3], mapPointSquares, optEdges, distsOpt)

		counterPOPS = counterPOPS + counter
	}
	println("dijkstra opt pops")
	println(counterPOPS)
	t = time.Now()
	println(t.String())

	println("dijksta opt rnd old")
	t = time.Now()
	println(t.String())
	counterPOPS = 0

	for i := 0; i <= -100; i++ {
		//println("sdfsdfs")
		var _, _, _, counter = dijkstra_single_optimized_old(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3], mapPointSquares, optEdges)

		counterPOPS = counterPOPS + counter
	}
	println("dijkstra opt pops old")
	println(counterPOPS)
	t = time.Now()
	println(t.String())

	println("a stern rnd")
	t = time.Now()
	println(t.String())
	counterPOPS = 0
	for i := 0; i <= -1; i++ {

		var _, _, _, counter = a_stern_single(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3])
		counterPOPS = counterPOPS + counter

	}
	println("pops a stern rnd")
	println(counterPOPS)

	t = time.Now()
	println(t.String())

	println("a stern single opt")
	t = time.Now()
	println(t.String())
	counterPOPS = 0

	for i := 0; i <= -1; i++ {

		var _, _, _, counter = a_stern_single_optimized(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3], mapPointSquares, optEdges)
		counterPOPS = counterPOPS + counter

	}
	println("pops a stern single opt ")
	println(counterPOPS)

	println("a stern single opt with pre")
	t = time.Now()
	println(t.String())
	counterPOPS = 0

	for i := 0; i <= -1; i++ {

		var _, _, _, counter = a_stern_single_optimized_with_pre(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3], mapPointSquares, optEdges, distsOpt)
		counterPOPS = counterPOPS + counter

	}
	println("pops a stern single opt with pre")
	println(counterPOPS)

	println("dijkstra single")
	t = time.Now()
	println(t.String())
	counterPOPS = 0

	for i := 0; i <= -1; i++ {
		//println("sdfsdfs")
		var _, _, _, counter = dijkstra_single(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3])

		counterPOPS = counterPOPS + counter
	}
	println("dijkstra single")

	t = time.Now()
	println(t.String())
	println("dijkstra single pops")
	println(counterPOPS)
	counterPOPS = 0

	for i := 0; i <= 10; i++ {

		//	var _, _, _, _, _, counter = a_stern_bi(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3], mapPointSquares)
		//	counterPOPS = counterPOPS + counter

	}

	println("a stern bi pops")
	println(counterPOPS)
	t = time.Now()
	println(t.String())

	println("a stern bi vergleich")
	t = time.Now()
	println(t.String())
	counterPOPS = 0

	for i := 0; i <= -1; i++ {

		var _, dists, _, counter = a_stern_single(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3])
		counterPOPS = counterPOPS + counter

		var _, dists2, _, _ = a_stern_single_optimized(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3], mapPointSquares, optEdges)
		counterPOPS = counterPOPS + counter

		if dists2[rnd[i][2]*len(result[0])+rnd[i][3]] > dists[rnd[i][2]*len(result[0])+rnd[i][3]] {
			println("____________________________________________________________________________________________")
		}

	}
	t = time.Now()
	println(t.String())

	println("test")
	for i := 0; i <= 100; i++ {

		var _, dis, _, _ = dijkstra_single(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3])

		var _, dis2, _, _ = a_stern_single_optimized_with_pre(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3], mapPointSquares, optEdges, distsOpt)

		var _, dis3, _, _ = a_stern_single_optimized(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3], mapPointSquares, optEdges)

		if dis2[rnd[i][2]*len(result[0])+rnd[i][3]] > (dis[rnd[i][2]*len(result[0])+rnd[i][3]] + 0.1) {
			println("error")
			println(getCordsFromArrayPosition(len(result), len(result[0]), rnd[i][0], rnd[i][1])[0])
			println(getCordsFromArrayPosition(len(result), len(result[0]), rnd[i][0], rnd[i][1])[1])
			println(getCordsFromArrayPosition(len(result), len(result[0]), rnd[i][2], rnd[i][3])[0])
			println(getCordsFromArrayPosition(len(result), len(result[0]), rnd[i][2], rnd[i][3])[1])
		}
		if dis3[rnd[i][2]*len(result[0])+rnd[i][3]] > dis[rnd[i][2]*len(result[0])+rnd[i][3]]+0.1 {
			println("error2")
		}
	}

	for i := 0; i <= len(result)-1; i = i + 10 {
		for j := 0; j <= len(result[i])-1; j = j + 10 {

			if len(mapPointSquares[[2]int{i, j}]) == 1 {
				print("X")
			} else {
				print(" ")
			}
		}
		println("")
	}

}
