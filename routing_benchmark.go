package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var result [][]bool
var optEdges [][2][2]int
var distsOpt [][][]float64

var mapPointSquares = make(map[[2]int][]int)

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

	println("dij single")
	t := time.Now()
	println(t.String())
	var counterPOPS = 0

	for i := 0; i <= 100; i++ {

		var _, _, _, counter = dijkstra_single(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3])
		counterPOPS = counterPOPS + counter

	}
	println("dijkstra single pops")
	println(counterPOPS)
	t = time.Now()
	println(t.String())

	println("dij bi")
	t = time.Now()
	println(t.String())
	counterPOPS = 0

	for i := 0; i <= -100; i++ {

		var _, _, _, _, _, counter = dijkstra_bi(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3])
		counterPOPS = counterPOPS + counter

	}
	println("dijkstra bi pops")
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

	for i := 0; i <= -100; i++ {
		//println("sdfsdfs")
		var _, _, _, counter = dijkstra_single_optimized(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3], mapPointSquares, optEdges)

		counterPOPS = counterPOPS + counter
	}
	println("dijkstra opt pops")
	println(counterPOPS)
	t = time.Now()
	println(t.String())

	println("dijksta opt rnd pre")
	t = time.Now()
	println(t.String())
	counterPOPS = 0

	for i := 0; i <= 100; i++ {
		//println("sdfsdfs")
		var _, _, _, counter = dijkstra_single_optimized_pre(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3], mapPointSquares, optEdges, distsOpt)

		counterPOPS = counterPOPS + counter
	}
	println("dijkstra opt pops pre")
	println(counterPOPS)
	t = time.Now()
	println(t.String())

	println("a stern rnd")
	t = time.Now()
	println(t.String())
	counterPOPS = 0
	for i := 0; i <= 100; i++ {

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

	for i := 0; i <= 100; i++ {

		var _, _, _, counter = a_stern_single_optimized(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3], mapPointSquares, optEdges)
		counterPOPS = counterPOPS + counter

	}
	println("pops a stern single opt ")
	println(counterPOPS)

	println("a stern single opt with pre")
	t = time.Now()
	println(t.String())
	counterPOPS = 0

	for i := 0; i <= 100; i++ {

		var _, _, _, counter = a_stern_single_optimized_with_pre(rnd[i][0], rnd[i][1], rnd[i][2], rnd[i][3], mapPointSquares, optEdges, distsOpt)
		counterPOPS = counterPOPS + counter

	}
	println("pops a stern single opt with pre")
	println(counterPOPS)
	t = time.Now()
	println(t.String())
	counterPOPS = 0

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
