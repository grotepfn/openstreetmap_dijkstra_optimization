package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

func TestSquares(t *testing.T) {

	jsonFile, err := os.Open("../../data/bitArray")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &result)

	var optEdges [][2][2]int
	var distsOpt [][][]float64
	jsonFile, err = os.Open("../../data/optimization_squares")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ = ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &optEdges)

	if (len(optEdges)) == 0 {
		t.Errorf("error")
	}

	for i := 0; i <= len(optEdges)-1; i++ {

		for j := optEdges[i][0][0]; j <= optEdges[i][1][0]; j++ {

			for k := optEdges[i][0][1]; k <= optEdges[i][1][1]; k++ {

				if result[j][k] != true {
					t.Errorf(strconv.Itoa(j))
					t.Errorf(strconv.Itoa(k))
					t.Errorf("error")
				}

			}

		}
	}
	var mapPointSquares = make(map[[2]int][]int)

	for i := 0; i <= len(optEdges)-1; i++ {

		for j := 0; j <= optEdges[i][1][0]-optEdges[i][0][0]; j++ {

			for k := 0; k <= optEdges[i][1][1]-optEdges[i][0][1]; k++ {

				var list = mapPointSquares[[2]int{optEdges[i][0][0] + j, optEdges[i][0][1] + k}]
				list = append(list, i)
				mapPointSquares[[2]int{optEdges[i][0][0] + j, optEdges[i][0][1] + k}] = list

			}
		}

	}

	for i := 0; i <= len(result)-1; i = i + 1 {
		for j := 0; j <= len(result[i])-1; j = j + 1 {

			if len(mapPointSquares[[2]int{i, j}]) <= 1 {

			} else {
				t.Errorf("error")
				t.Errorf(strconv.Itoa(len(mapPointSquares[[2]int{i, j}])))

			}
		}

	}
	for k := 0; k <= len(optEdges)-1; k++ {

		for xAxis := optEdges[k][0][1] - 1; xAxis <= optEdges[k][1][1]+1; xAxis++ {

			if len(mapPointSquares[[2]int{optEdges[k][0][0] - 1, xAxis}]) > 0 {
				t.Errorf("error 1")
				t.Errorf(strconv.Itoa(optEdges[k][0][0]))
				t.Errorf(strconv.Itoa(optEdges[k][0][1]))
				t.Errorf(strconv.Itoa(optEdges[k][1][0]))
				t.Errorf(strconv.Itoa(optEdges[k][1][1]))

				var x = mapPointSquares[[2]int{optEdges[k][0][0] - 1, xAxis}]

				t.Errorf(strconv.Itoa(optEdges[x[0]][0][0]))
				t.Errorf(strconv.Itoa(optEdges[x[0]][0][1]))
				t.Errorf(strconv.Itoa(optEdges[x[0]][1][0]))
				t.Errorf(strconv.Itoa(optEdges[x[0]][1][1]))

			}

			//lower line

			if len(mapPointSquares[[2]int{optEdges[k][1][0] + 1, xAxis}]) > 0 {
				t.Errorf("error 2")

			}
		}

		//left right
		for yAxis := optEdges[k][0][0] - 1; yAxis <= optEdges[k][1][0]+1; yAxis++ {

			if len(mapPointSquares[[2]int{yAxis, optEdges[k][0][1] - 1}]) > 0 {
				t.Errorf("error3")

			}

			if len(mapPointSquares[[2]int{yAxis, optEdges[k][1][1] + 1}]) > 0 {
				t.Errorf("error4")

			}

		}
	}

	for k := 0; k <= len(optEdges)-1; k++ {

		for xAxis := optEdges[k][0][1] - 2; xAxis <= optEdges[k][1][1]+2; xAxis++ {

			if len(mapPointSquares[[2]int{optEdges[k][0][0] - 1, xAxis}]) > 0 {
				t.Errorf("error 1")
				t.Errorf(strconv.Itoa(optEdges[k][0][0]))
				t.Errorf(strconv.Itoa(optEdges[k][0][1]))
				t.Errorf(strconv.Itoa(optEdges[k][1][0]))
				t.Errorf(strconv.Itoa(optEdges[k][1][1]))

				var x = mapPointSquares[[2]int{optEdges[k][0][0] - 1, xAxis}]

				t.Errorf(strconv.Itoa(optEdges[x[0]][0][0]))
				t.Errorf(strconv.Itoa(optEdges[x[0]][0][1]))
				t.Errorf(strconv.Itoa(optEdges[x[0]][1][0]))
				t.Errorf(strconv.Itoa(optEdges[x[0]][1][1]))

			}

			//lower line

			if len(mapPointSquares[[2]int{optEdges[k][1][0] + 1, xAxis}]) > 0 {
				t.Errorf("error 2")

			}
		}

		//left right
		for yAxis := optEdges[k][0][0] - 1; yAxis <= optEdges[k][1][0]+1; yAxis++ {

			if len(mapPointSquares[[2]int{yAxis, optEdges[k][0][1] - 1}]) > 0 {
				t.Errorf("error3")

			}

			if len(mapPointSquares[[2]int{yAxis, optEdges[k][1][1] + 1}]) > 0 {
				t.Errorf("error4")

			}

		}
	}

	jsonFile, err = os.Open("../../data/optimization_squares_distances")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ = ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &distsOpt)

	if (len(distsOpt)) == 0 {
		t.Errorf("error5")
	}

	for i := 0; i <= len(distsOpt)-1; i++ {

		if len(distsOpt[i][4]) != (optEdges[i][1][0]-optEdges[i][0][0])*4+4 {
			t.Errorf("error5")
			t.Errorf(strconv.Itoa(len(distsOpt[i][0])))
			t.Errorf(strconv.Itoa(optEdges[i][1][0] - optEdges[i][0][0]))
		} else {
			//	t.Errorf("error3242")
		}
		if len(distsOpt) != len(optEdges) {
			t.Errorf("error5")
		}
		if len(distsOpt[i]) != (optEdges[i][1][0]-optEdges[i][0][0])*4+4 {
			t.Errorf("error5")
			t.Errorf(strconv.Itoa(len(distsOpt[i][0])))
			t.Errorf(strconv.Itoa(optEdges[i][1][0] - optEdges[i][0][0]))
		}

	}

}
