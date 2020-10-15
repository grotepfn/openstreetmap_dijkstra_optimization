package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/grotepfn/openstreetmap_dijkstra_optimization/bitArray"
	geojson "github.com/paulmach/go.geojson"
)

var result [][]bool

func main() {

	jsonFile, err := os.Open("data/bitArray")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &result)
	println("done")

	var sq = findSquares()
	println(len(sq))

	jsonString1, _ := json.Marshal(sq)
	ioutil.WriteFile("data/optimization_squares", jsonString1, 0644)
	println("done writing")

	var sq2 = make([][][]float64, len(sq))

	for i := 0; i <= len(sq)-1; i++ {

		//for k := 0; k <= len(sq[i][j])-1; k++ {
		sq2[i] = append(sq2[i], []float64{bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[i][0][0], sq[i][0][1])[1], bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[i][0][0], sq[i][0][1])[0]})

		sq2[i] = append(sq2[i], []float64{bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[i][0][0], sq[i][1][1])[1], bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[i][0][0], sq[i][1][1])[0]})

		sq2[i] = append(sq2[i], []float64{bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[i][1][0], sq[i][1][1])[1], bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[i][1][0], sq[i][1][1])[0]})

		sq2[i] = append(sq2[i], []float64{bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[i][1][0], sq[i][0][1])[1], bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[i][1][0], sq[i][0][1])[0]})

		sq2[i] = append(sq2[i], []float64{bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[i][0][0], sq[i][0][1])[1], bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[i][0][0], sq[i][0][0])[0]})

	}

	var greatCircleDistances [][][]float64

	for k := 0; k <= len(sq)-1; k++ {
		var oneSquare [][]float64

		for x := sq[k][0][1]; x <= sq[k][1][1]; x++ {
			var help []float64
			for xAxis := sq[k][0][1]; xAxis <= sq[k][1][1]; xAxis++ {

				//upper line
				var distance = bitArray.GreatCircleDistance(bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[k][0][0], x), bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[k][0][0], xAxis))

				var l = help
				l = append(l, distance)

				help = l
			}

			for xAxis := sq[k][0][1]; xAxis <= sq[k][1][1]; xAxis++ {

				//lower line

				var distance = bitArray.GreatCircleDistance(bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[k][0][0], x), bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[k][1][0], xAxis))

				var l = help
				l = append(l, distance)

				help = l

			}

			//left right
			for yAxis := sq[k][0][0]; yAxis <= sq[k][1][0]; yAxis++ {

				//left line

				var distance = bitArray.GreatCircleDistance(bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[k][0][0], x), bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), yAxis, sq[k][0][1]))

				var l = help
				l = append(l, distance)

				help = l
			}

			for yAxis := sq[k][0][0]; yAxis <= sq[k][1][0]; yAxis++ {
				//right line

				var distance = bitArray.GreatCircleDistance(bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[k][0][0], x), bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), yAxis, sq[k][1][1]))

				var l = help
				l = append(l, distance)

				help = l
			}

			oneSquare = append(oneSquare, help)

		}

		for x2 := sq[k][0][1]; x2 <= sq[k][1][1]; x2++ {
			var help []float64
			for xAxis := sq[k][0][1]; xAxis <= sq[k][1][1]; xAxis++ {

				//upper line
				var distance = bitArray.GreatCircleDistance(bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[k][1][0], x2), bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[k][0][0], xAxis))

				var l = help
				l = append(l, distance)

				help = l
			}
			for xAxis := sq[k][0][1]; xAxis <= sq[k][1][1]; xAxis++ {
				//lower line

				var distance = bitArray.GreatCircleDistance(bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[k][1][0], x2), bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[k][1][0], xAxis))

				var l = help
				l = append(l, distance)

				help = l
			}

			//left right
			for yAxis := sq[k][0][0]; yAxis <= sq[k][1][0]; yAxis++ {

				//left line

				var distance = bitArray.GreatCircleDistance(bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[k][1][0], x2), bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), yAxis, sq[k][0][1]))

				var l = help
				l = append(l, distance)

				help = l

			}
			for yAxis := sq[k][0][0]; yAxis <= sq[k][1][0]; yAxis++ {
				//right line

				var distance = bitArray.GreatCircleDistance(bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[k][1][0], x2), bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), yAxis, sq[k][1][1]))

				var l = help
				l = append(l, distance)

				help = l
			}

			oneSquare = append(oneSquare, help)
		}

		for y := sq[k][0][0]; y <= sq[k][1][0]; y++ {
			var help []float64
			for xAxis := sq[k][0][1]; xAxis <= sq[k][1][1]; xAxis++ {

				//upper line
				var distance = bitArray.GreatCircleDistance(bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), y, sq[k][0][1]), bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[k][0][0], xAxis))

				var l = help
				l = append(l, distance)
				help = l
			}
			for xAxis := sq[k][0][1]; xAxis <= sq[k][1][1]; xAxis++ {

				//lower line

				var distance = bitArray.GreatCircleDistance(bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), y, sq[k][0][1]), bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[k][1][0], xAxis))
				var l = help
				l = append(l, distance)

				help = l

			}

			//left right
			for yAxis := sq[k][0][0]; yAxis <= sq[k][1][0]; yAxis++ {

				//left line

				var distance = bitArray.GreatCircleDistance(bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), y, sq[k][0][1]), bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), yAxis, sq[k][0][1]))

				var l = help
				l = append(l, distance)
				help = l

			}
			for yAxis := sq[k][0][0]; yAxis <= sq[k][1][0]; yAxis++ {

				//right line

				var distance = bitArray.GreatCircleDistance(bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), y, sq[k][0][1]), bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), yAxis, sq[k][1][1]))
				var l = help
				l = append(l, distance)

				help = l
			}
			oneSquare = append(oneSquare, help)
		}

		for y2 := sq[k][0][0]; y2 <= sq[k][1][0]; y2++ {
			var help []float64
			for xAxis := sq[k][0][1]; xAxis <= sq[k][1][1]; xAxis++ {

				//upper line
				var distance = bitArray.GreatCircleDistance(bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), y2, sq[k][1][1]), bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[k][0][0], xAxis))

				var l = help
				l = append(l, distance)

				help = l

			}
			for xAxis := sq[k][0][1]; xAxis <= sq[k][1][1]; xAxis++ {
				//lower line

				var distance = bitArray.GreatCircleDistance(bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), y2, sq[k][1][1]), bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), sq[k][1][0], xAxis))

				var l = help
				l = append(l, distance)

				help = l

			}

			//left right
			for yAxis := sq[k][0][0]; yAxis <= sq[k][1][0]; yAxis++ {

				//left line

				var distance = bitArray.GreatCircleDistance(bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), y2, sq[k][1][1]), bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), yAxis, sq[k][0][1]))
				var l = help
				l = append(l, distance)

				help = l
			}
			for yAxis := sq[k][0][0]; yAxis <= sq[k][1][0]; yAxis++ {

				//right line

				var distance = bitArray.GreatCircleDistance(bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), y2, sq[k][1][1]), bitArray.GetCordsFromArrayPosition(len(result), len(result[0]), yAxis, sq[k][1][1]))
				var l = help
				l = append(l, distance)

				help = l
			}
			oneSquare = append(oneSquare, help)
		}

		greatCircleDistances = append(greatCircleDistances, oneSquare)
	}

	jsonString1, _ = json.Marshal(greatCircleDistances)
	ioutil.WriteFile("data/optimization_squares_distances", jsonString1, 0644)
	println("done writing")

	g := geojson.NewPolygonFeature(sq2[0:844])

	jsonString1, _ = json.Marshal(g)
	ioutil.WriteFile("data/geoJson", jsonString1, 0644)
	println("done writing")

}

func findSquares() [][2][2]int {

	var s [][2][2]int

	for i := 2; i <= len(result)-2; i++ {
		//inner:
		for j := 2; j <= len(result[i])-2; j++ {

			//var length, length2 = findBiggestSquare(i, j, s)
			var length3 = findBiggestSquare2(i, j, s)
			if length3 >= 5 { //&& length2 >= 7 && math.Max(float64(length), float64(length2))/math.Min(float64(length), float64(length2)) <= 1.5 {

				s = append(s, [2][2]int{{i, j}, {i + length3, j + length3}})
			}

		}

	}
	return s
}

func findBiggestSquare(i, j int, s [][2][2]int) (int, int) {

	var k = 0
	var l = 0
	var l_max = 0
	for i+k <= len(result)-1 && result[i+k][j] == true && k <= int(1.2*float64(l_max)) {
		l = 0
		for j+l <= len(result[i])-1 && l <= l_max && result[i][j+l] == true {

			for m := 0; m <= l; m++ {

				for z := 0; z <= len(s)-1; z++ {

					if i+k >= s[z][0][0] && j+m >= s[z][0][1] && i+k <= s[z][1][0] && j+m <= s[z][1][1] {

						return -1, -1
					}

					//links
					if i+k >= s[z][0][0] && i+k <= s[z][1][0] && j+m >= s[z][0][1]-2 && j+m <= s[z][0][1] {

						return k - 1, l_max - 1
					}
					//rechts
					if i+k >= s[z][0][0] && i+k <= s[z][1][0] && j+m <= s[z][1][1]+2 && j+m >= s[z][1][1] {

						return k - 1, l_max - 1
					}

					//oben!!!!! falsch?
					if j+m >= s[z][0][1] && j+m <= s[z][1][1] && i+k >= s[z][0][0]-2 && i+k <= s[z][0][0] {

						return k - 1, l_max - 1
					}
					//unten
					if j+m >= s[z][0][1] && j+m <= s[z][1][1] && i+k <= s[z][1][0]+2 && i+k >= s[z][1][0] {

						return k - 1, l_max - 1
					}

				}

				if (result[i+k][j+m]) != true {

					return k - 1, l_max - 1
				}
			}

			l++
			if k == 0 {
				l_max = l
			}
		}
		k++
	}

	return k - 3, l_max - 3
}

func findBiggestSquare2(i, j int, s [][2][2]int) int {
	var k = 0
	for result[i+k][j+k] == true && i+k <= len(result)-3 && j+k <= len(result[0])-3 {

		for y := i - 2; y <= i+k+2; y++ {
			for x := j - 2; x <= j+k+2; x++ {

				//for m := 0; m <= k; m++ {

				for z := 0; z <= len(s)-1; z++ {

					if x >= s[z][0][1]-2 && y >= s[z][0][0]-2 && x <= s[z][1][1]+2 && y <= s[z][1][0]+2 {

						return k - 2
					}

					//links
					if x >= s[z][0][1]-2 && x <= s[z][0][1] && y >= s[z][0][0]-2 && y <= s[z][1][0]+2 {

						return k - 2
					}
					//rechts
					if x >= s[z][1][1] && x <= s[z][1][1]+2 && y <= s[z][0][0]-2 && y >= s[z][1][0]+2 {

						return k - 2
					}

					//oben
					if y >= s[z][0][0]-2 && y <= s[z][0][0] && x >= s[z][0][1]-2 && x <= s[z][1][1]+2 {

						return k - 2
					}
					//unten
					if y <= s[z][1][0]+2 && y >= s[z][1][0] && x <= s[z][1][1]+2 && x >= s[z][0][1]-2 {

						return k - 2
					}

				}

				//	}

				if result[y][x] != true {
					return k - 2
				}

			}
		}
		k++
	}

	return k - 2
}
