package bitArray

import (
	"container/heap"
	"math"
)

//From https://golang.org/pkg/container/heap/
// An Item is something we manage in a priority queue.
type Item struct {
	pos                    [2]int  // The value of the item; arbitrary.
	distancePriority       float64 // The priority of the item in the queue.
	actualDistanceForAStar float64 //dijkstra
	// The index is needed by update and is maintained by the heap.Interface methods.
	index     int // The index of the item in the heap.
	optimized bool
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest, not highest, priority so we use smaller than here.
	return pq[i].distancePriority < pq[j].distancePriority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, pos [2]int, distance float64) {
	item.pos = pos
	item.distancePriority = distance
	heap.Fix(pq, item.index)
}

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.

//func dijktstra_minheap(pos1, pos2, pos3, pos4 int) ([][][]int, [][]float64) {
//}

//Implementation fo Dijkstras algorithm without optimization https://dev.to/douglasmakey/implementation-of-dijkstra-using-heap-in-go-6e3
func Dijkstra_single(pos1, pos2, pos3, pos4 int, result [][]bool) ([][]int, []float64, [2]int, int) {

	var counterPOPS = 0

	distances := make([]float64, len(result)*len(result[0]))

	pq := make(PriorityQueue, 0)
	pre := make([][]int, len(result)*len(result[0]))

	for i := 0; i < len(result); i++ {

		pre[i] = make([]int, len(result[0]))

	}
	for i := 0; i <= len(result)-1; i++ {
		for j := 0; j <= len(result[i])-1; j++ {

			distances[i*(len(result[0]))+j] = math.MaxFloat64

			pre[i*(len(result[0]))+j] = []int{-1, -1}

			if i == pos1 && j == pos2 {
				// Insert a new item and then modify its priority.
				item := &Item{
					pos:              [2]int{pos1, pos2},
					distancePriority: 0,
				}
				heap.Push(&pq, item) // Money Boy

			}
		}
	}

	distances[pos1*(len(result[0]))+pos2] = 0

	var node [2]int

	for pq.Len() > 0 {

		// Find the nearest yet to visit node
		p := heap.Pop(&pq).(*Item)

		node = p.pos

		if pre[node[0]*len(result[0])+node[1]][0] != -1 && p.distancePriority > distances[node[0]*len(result[0])+node[1]] {

			continue
		}

		counterPOPS++

		var neighbours = GetEdges(result, node[0], node[1])
		for k := 0; k <= len(neighbours)-1; k++ {

			var distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node[0], node[1]), GetCordsFromArrayPosition(len(result), len(result[0]), neighbours[k][0], neighbours[k][1]))

			var alt = distances[node[0]*(len(result[0]))+node[1]] + distance

			if alt < distances[neighbours[k][0]*(len(result[0]))+neighbours[k][1]] {

				newItem := &Item{pos: neighbours[k], distancePriority: alt}
				pre[neighbours[k][0]*(len(result[0]))+neighbours[k][1]] = []int{node[0], node[1]}
				distances[neighbours[k][0]*(len(result[0]))+neighbours[k][1]] = alt
				heap.Push(&pq, newItem)

			}

		}

		if node[0] == pos3 && node[1] == pos4 {

			return pre, distances, node, counterPOPS

		}

	}

	return pre, distances, node, counterPOPS

}

//A star implementation optimized without precalculation of the distances between the points of the optimization squares
func A_stern_single_optimized(result [][]bool, pos1, pos2, pos3, pos4 int, mapPointSquares map[[2]int]int, optEdges [][2][2]int) ([][]int, []float64, [2]int, int) {

	var visitedSquares = map[int]bool{}
	suchraum := make([]bool, len(result)*len(result[0]))
	var counterPOPS = 0
	var optimizedAmount = 0
	distances := make([]float64, len(result)*len(result[0]))

	visitedPoints := make([]bool, len(result)*len(result[0]))

	pq := make(PriorityQueue, 0)
	pre := make([][]int, len(result)*len(result[0]))

	for i := 0; i < len(result); i++ {

		pre[i] = make([]int, len(result[0]))

	}

	for i := 0; i <= len(result)-1; i++ {

		for j := 0; j <= len(result[i])-1; j++ {

			distances[i*len(result[0])+j] = math.MaxFloat64
			suchraum[i*len(result[0])+j] = false

			visitedPoints[i*len(result[0])+j] = false

			pre[i*(len(result[0]))+j] = []int{-1, -1}

			if i == pos1 && j == pos2 {
				// Insert a new item and then modify its priority.
				item := &Item{
					pos:                    [2]int{pos1, pos2},
					distancePriority:       GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), pos1, pos2), GetCordsFromArrayPosition(len(result), len(result[0]), pos3, pos4)),
					actualDistanceForAStar: 0,
				}
				heap.Push(&pq, item)

			}
		}
	}

	distances[pos1*(len(result[0]))+pos2] = 0

	var node [2]int

	for pq.Len() > 0 {

		// Find the nearest yet to visit node
		p := heap.Pop(&pq).(*Item)

		node = p.pos
		var isOptimized = p.optimized

		if pre[node[0]*len(result[0])+node[1]][0] != -1 && p.actualDistanceForAStar > distances[node[0]*len(result[0])+node[1]] {

			continue
		}

		counterPOPS++

		var k = -1
		var ok = false
		k, ok = mapPointSquares[[2]int{node[0], node[1]}]

		if !isOptimized && ok && node[0] >= optEdges[k][0][0] && node[1] >= optEdges[k][0][1] && node[0] <= optEdges[k][1][0] && node[1] <= optEdges[k][1][1] && ((pos3 < optEdges[k][0][0] || pos4 < optEdges[k][0][1]) || (pos3 > optEdges[k][1][0] || pos4 > optEdges[k][1][1])) {
			if !visitedSquares[k] {
				optimizedAmount += (optEdges[k][1][1] - optEdges[k][0][1]) * (optEdges[k][1][0] - optEdges[k][0][0])
			}
			visitedSquares[k] = true

			for xAxis := optEdges[k][0][1]; xAxis <= optEdges[k][1][1]; xAxis++ {

				if visitedPoints[optEdges[k][0][0]*len(result[0])+xAxis] == true {

				}

				if visitedPoints[optEdges[k][1][0]*len(result[0])+xAxis] == true {

					continue
				}

				//upper line
				var distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node[0], node[1]), GetCordsFromArrayPosition(len(result), len(result[0]), optEdges[k][0][0], xAxis))

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				var altHelp = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), optEdges[k][0][0], xAxis), GetCordsFromArrayPosition(len(result), len(result[0]), pos3, pos4))

				if alt < distances[optEdges[k][0][0]*len(result[0])+xAxis] {
					newItem := &Item{pos: [2]int{optEdges[k][0][0], xAxis}, distancePriority: alt + altHelp, optimized: true, actualDistanceForAStar: alt}
					pre[optEdges[k][0][0]*len(result[0])+xAxis] = []int{node[0], node[1]}
					distances[optEdges[k][0][0]*len(result[0])+xAxis] = alt

					heap.Push(&pq, newItem)
					if visitedPoints[optEdges[k][0][0]*len(result[0])+xAxis] == true {

					}
				}

				//lower line

				distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node[0], node[1]), GetCordsFromArrayPosition(len(result), len(result[0]), optEdges[k][1][0], xAxis))

				alt = distances[node[0]*len(result[0])+node[1]] + distance
				altHelp = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), optEdges[k][1][0], xAxis), GetCordsFromArrayPosition(len(result), len(result[0]), pos3, pos4))

				if alt < distances[optEdges[k][1][0]*len(result[0])+xAxis] {
					newItem2 := &Item{pos: [2]int{optEdges[k][1][0], xAxis}, distancePriority: alt + altHelp, optimized: true, actualDistanceForAStar: alt}
					pre[optEdges[k][1][0]*len(result[0])+xAxis] = []int{node[0], node[1]}
					distances[optEdges[k][1][0]*len(result[0])+xAxis] = alt

					heap.Push(&pq, newItem2)

					if visitedPoints[optEdges[k][1][0]*len(result[0])+xAxis] == true {

					}
				}

			}

			//left right
			for yAxis := optEdges[k][0][0]; yAxis <= optEdges[k][1][0]; yAxis++ {

				//left line

				var distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node[0], node[1]), GetCordsFromArrayPosition(len(result), len(result[0]), yAxis, optEdges[k][0][1]))

				var alt = distances[node[0]*len(result[0])+node[1]] + distance
				var altHelp = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), yAxis, optEdges[k][0][1]), GetCordsFromArrayPosition(len(result), len(result[0]), pos3, pos4))

				if alt < distances[(yAxis)*len(result[0])+optEdges[k][0][1]] {
					newItem := &Item{pos: [2]int{yAxis, optEdges[k][0][1]}, distancePriority: alt + altHelp, optimized: true, actualDistanceForAStar: alt}
					pre[(yAxis)*len(result[0])+optEdges[k][0][1]] = []int{node[0], node[1]}
					distances[(yAxis)*len(result[0])+optEdges[k][0][1]] = alt

					heap.Push(&pq, newItem)

				}

				//right line

				distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node[0], node[1]), GetCordsFromArrayPosition(len(result), len(result[0]), yAxis, optEdges[k][1][1]))

				alt = distances[node[0]*len(result[0])+node[1]] + distance
				altHelp = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), yAxis, optEdges[k][1][1]), GetCordsFromArrayPosition(len(result), len(result[0]), pos3, pos4))

				if alt < distances[(yAxis)*len(result[0])+optEdges[k][1][1]] {
					newItem2 := &Item{pos: [2]int{yAxis, optEdges[k][1][1]}, distancePriority: alt + altHelp, optimized: true, actualDistanceForAStar: alt}
					pre[(yAxis)*len(result[0])+optEdges[k][1][1]] = []int{node[0], node[1]}
					distances[(yAxis)*len(result[0])+optEdges[k][1][1]] = alt

					heap.Push(&pq, newItem2)

				}
			}

			//set nodes within negative
			for yAxis := optEdges[k][0][0] + 1; yAxis < optEdges[k][1][0]; yAxis++ {
				for xAxis := optEdges[k][0][1] + 1; xAxis < optEdges[k][1][1]; xAxis++ {
					visitedPoints[yAxis*len(result[0])+xAxis] = true
				}
			}

		}

		//destination point within square
		if ok && node[0] >= optEdges[k][0][0] && node[1] >= optEdges[k][0][1] && node[0] <= optEdges[k][1][0] && node[1] <= optEdges[k][1][1] && (pos3 >= optEdges[k][0][0] && pos4 >= optEdges[k][0][1] && pos3 <= optEdges[k][1][0] && pos4 <= optEdges[k][1][1]) {

			var distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node[0], node[1]), GetCordsFromArrayPosition(len(result), len(result[0]), pos3, pos4))

			var alt = distances[node[0]*len(result[0])+node[1]] + distance

			if alt <= distances[pos3*len(result[0])+pos4] {

				pre[pos3*len(result[0])+pos4] = []int{node[0], node[1]}
				distances[pos3*len(result[0])+pos4] = alt

			}

			return pre, distances, node, counterPOPS
		}

		var neighbours = GetEdges(result, node[0], node[1])

	neig:
		for k := 0; k <= len(neighbours)-1; k++ {

			if visitedPoints[neighbours[k][0]*len(result[0])+neighbours[k][1]] == true {

				continue neig
			}

			var distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node[0], node[1]), GetCordsFromArrayPosition(len(result), len(result[0]), neighbours[k][0], neighbours[k][1]))

			var alt = distances[node[0]*len(result[0])+node[1]] + distance

			if alt < distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] {
				newItem := &Item{pos: neighbours[k], actualDistanceForAStar: alt, distancePriority: alt + GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), neighbours[k][0], neighbours[k][1]), GetCordsFromArrayPosition(len(result), len(result[0]), pos3, pos4))}
				pre[neighbours[k][0]*len(result[0])+neighbours[k][1]] = []int{node[0], node[1]}
				distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] = alt

				heap.Push(&pq, newItem)

			}

		}

		if node[0] == pos3 && node[1] == pos4 {

			return pre, distances, node, counterPOPS

		}

	}

	return pre, distances, node, counterPOPS

}

//Dijkstra implementation optimized with precalculation of the distances between the points of the optimization squares
func Dijkstra_single_optimized_pre(result [][]bool, pos1, pos2, pos3, pos4 int, mapPointSquares map[[2]int]int, optEdges [][2][2]int, dists [][][]float64) ([][]int, []float64, [2]int, int) {

	var counterPOPS = 0

	distances := make([]float64, len(result)*len(result[0]))

	visitedPointsInSquare := make([]bool, len(result)*len(result[0]))

	pq := make(PriorityQueue, 0)
	pre := make([][]int, len(result)*len(result[0]))

	for i := 0; i < len(result); i++ {

		pre[i] = make([]int, len(result[0]))

	}

	for i := 0; i <= len(result)-1; i++ {

		for j := 0; j <= len(result[i])-1; j++ {

			distances[i*len(result[0])+j] = math.MaxFloat64

			visitedPointsInSquare[i*len(result[0])+j] = false

			pre[i*(len(result[0]))+j] = []int{-1, -1}

			if i == pos1 && j == pos2 {
				// Insert a new item and then modify its priority.
				item := &Item{
					pos:              [2]int{pos1, pos2},
					distancePriority: 0,
				}
				heap.Push(&pq, item)

			}
		}
	}

	distances[pos1*(len(result[0]))+pos2] = 0

	var node [2]int
	//var counter = 0

	for pq.Len() > 0 {

		// Find the nearest yet to visit node
		p := heap.Pop(&pq).(*Item)

		node = p.pos
		var isOptimized = p.optimized

		if pre[node[0]*len(result[0])+node[1]][0] != -1 && p.distancePriority > distances[node[0]*len(result[0])+node[1]] {

			continue
		}

		if visitedPointsInSquare[node[0]*len(result[0])+node[1]] == true {

			continue
		}

		counterPOPS++
		var k = -1
		var ok = false
		k, ok = mapPointSquares[[2]int{node[0], node[1]}]

		//node within square and dest pos not within square
		if !isOptimized && ok && node[0] >= optEdges[k][0][0] && node[1] >= optEdges[k][0][1] && node[0] <= optEdges[k][1][0] && node[1] <= optEdges[k][1][1] && ((pos3 < optEdges[k][0][0] || pos4 < optEdges[k][0][1]) || (pos3 > optEdges[k][1][0] || pos4 > optEdges[k][1][1])) && (node[0] == optEdges[k][0][0] || node[1] == optEdges[k][0][1] || node[0] == optEdges[k][1][0] || node[1] == optEdges[k][1][1]) && ((pos1 < optEdges[k][0][0] || pos2 < optEdges[k][0][1]) || (pos1 > optEdges[k][1][0] || pos2 > optEdges[k][1][1])) {

			var zahl = 0

			var foundSquarePoint = false
			//upper
			if !foundSquarePoint && node[1] >= optEdges[k][0][1] && node[1] <= optEdges[k][1][1] && node[0] == optEdges[k][0][0] {

				zahl = zahl + node[1] - optEdges[k][0][1]
				foundSquarePoint = true
			} else if !foundSquarePoint {
				zahl = zahl + optEdges[k][1][1] - optEdges[k][0][1] + 1
			}

			//lower
			if !foundSquarePoint && node[1] >= optEdges[k][0][1] && node[1] <= optEdges[k][1][1] && node[0] == optEdges[k][1][0] {

				zahl = zahl + node[1] - optEdges[k][0][1]
				foundSquarePoint = true
			} else if !foundSquarePoint {
				zahl = zahl + optEdges[k][1][1] - optEdges[k][0][1] + 1
			}

			//left side
			if !foundSquarePoint && node[0] >= optEdges[k][0][0] && node[0] <= optEdges[k][1][0] && node[1] == optEdges[k][0][1] {

				zahl = zahl + node[0] - optEdges[k][0][0]
				foundSquarePoint = true
			} else if !foundSquarePoint {

				zahl = zahl + optEdges[k][1][0] - optEdges[k][0][0] + 1
			}

			//right side
			if !foundSquarePoint && node[0] >= optEdges[k][0][0] && node[0] <= optEdges[k][1][0] && node[1] == optEdges[k][1][1] {

				zahl = zahl + node[0] - optEdges[k][0][0]
				foundSquarePoint = true
			} else if !foundSquarePoint {

				zahl = zahl + optEdges[k][1][0] - optEdges[k][0][0]

			}

			var help = dists[k][zahl]

			for i := 0; i <= optEdges[k][1][1]-optEdges[k][0][1]; i++ {

				var distance = help[i]

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[optEdges[k][0][0]*len(result[0])+(optEdges[k][0][1]+i)] {

					newItem := &Item{pos: [2]int{optEdges[k][0][0], (optEdges[k][0][1] + i)}, optimized: true, distancePriority: alt}
					pre[optEdges[k][0][0]*len(result[0])+(optEdges[k][0][1]+i)] = []int{node[0], node[1]}
					distances[optEdges[k][0][0]*len(result[0])+(optEdges[k][0][1]+i)] = alt

					heap.Push(&pq, newItem)

				}

			}

			for i := 0; i <= optEdges[k][1][1]-optEdges[k][0][1]; i++ {

				var distance = help[1+optEdges[k][1][1]-optEdges[k][0][1]+i]

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[optEdges[k][1][0]*len(result[0])+i+optEdges[k][0][1]] {

					newItem := &Item{pos: [2]int{optEdges[k][1][0], optEdges[k][0][1] + i}, optimized: true, distancePriority: alt}
					pre[optEdges[k][1][0]*len(result[0])+optEdges[k][0][1]+i] = []int{node[0], node[1]}
					distances[optEdges[k][1][0]*len(result[0])+optEdges[k][0][1]+i] = alt

					heap.Push(&pq, newItem)

				}

			}

			//left
			for i := 0; i <= optEdges[k][1][1]-optEdges[k][0][1]; i++ {

				var distance = help[2+2*(optEdges[k][1][1]-optEdges[k][0][1])+i]

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][0][1]] {

					newItem := &Item{pos: [2]int{(optEdges[k][0][0] + i), optEdges[k][0][1]}, optimized: true, distancePriority: alt}
					pre[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][0][1]] = []int{node[0], node[1]}
					distances[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][0][1]] = alt

					heap.Push(&pq, newItem)

				}

			}

			//right
			for i := 0; i <= optEdges[k][1][1]-optEdges[k][0][1]; i++ {

				var distance = help[3+3*(optEdges[k][1][1]-optEdges[k][0][1])+i]

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][1][1]] {

					newItem := &Item{pos: [2]int{(optEdges[k][0][0] + i), optEdges[k][1][1]}, optimized: true, distancePriority: alt}
					pre[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][1][1]] = []int{node[0], node[1]}
					distances[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][1][1]] = alt

					heap.Push(&pq, newItem)

				}

			}

			//set nodes within negative
			for yAxis := optEdges[k][0][0] + 1; yAxis < optEdges[k][1][0]; yAxis++ {
				for xAxis := optEdges[k][0][1] + 1; xAxis < optEdges[k][1][1]; xAxis++ {

					visitedPointsInSquare[yAxis*len(result[0])+xAxis] = true
				}
			}

		}

		var neighbours = GetEdges(result, node[0], node[1])

	neig:
		for k := 0; k <= len(neighbours)-1; k++ {

			if visitedPointsInSquare[neighbours[k][0]*len(result[0])+neighbours[k][1]] == true {

				continue neig
			}

			var distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node[0], node[1]), GetCordsFromArrayPosition(len(result), len(result[0]), neighbours[k][0], neighbours[k][1]))

			var alt = distances[node[0]*len(result[0])+node[1]] + distance

			if alt < distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] {
				newItem := &Item{pos: neighbours[k], distancePriority: alt}
				pre[neighbours[k][0]*len(result[0])+neighbours[k][1]] = []int{node[0], node[1]}
				distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] = alt

				heap.Push(&pq, newItem)

			}

		}

		if node[0] == pos3 && node[1] == pos4 {

			return pre, distances, node, counterPOPS

		}

	}

	return pre, distances, node, counterPOPS

}

//A star implementation without opimization
func A_stern_single(result [][]bool, pos1, pos2, pos3, pos4 int) ([][]int, []float64, [2]int, int) {

	var counterPOPS = 0

	distances := make([]float64, len(result)*len(result[0]))

	pq := make(PriorityQueue, 0)
	pre := make([][]int, len(result)*len(result[0]))

	for i := 0; i < len(result); i++ {

		pre[i] = make([]int, len(result[0]))

	}

	for i := 0; i <= len(result)-1; i++ {

		for j := 0; j <= len(result[i])-1; j++ {

			distances[i*len(result[0])+j] = math.MaxFloat64

			pre[i*(len(result[0]))+j] = []int{-1, -1}

			if i == pos1 && j == pos2 {
				// Insert a new item and then modify its priority.
				item := &Item{
					pos:              [2]int{pos1, pos2},
					distancePriority: GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), pos1, pos2), GetCordsFromArrayPosition(len(result), len(result[0]), pos3, pos4)),
				}
				heap.Push(&pq, item)

			}
		}
	}

	distances[pos1*(len(result[0]))+pos2] = 0

	var node [2]int

	for pq.Len() > 0 {

		p := heap.Pop(&pq).(*Item)

		node = p.pos

		if pre[node[0]*len(result[0])+node[1]][0] != -1 && (p.actualDistanceForAStar) > distances[node[0]*len(result[0])+node[1]] {

			continue
		}

		counterPOPS++

		var neighbours = GetEdges(result, node[0], node[1])

		for k := 0; k <= len(neighbours)-1; k++ {

			var distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node[0], node[1]), GetCordsFromArrayPosition(len(result), len(result[0]), neighbours[k][0], neighbours[k][1]))

			var alt = distances[node[0]*len(result[0])+node[1]] + distance

			if alt < distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] {

				newItem := &Item{pos: neighbours[k], distancePriority: alt + GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), neighbours[k][0], neighbours[k][1]), GetCordsFromArrayPosition(len(result), len(result[0]), pos3, pos4))}
				pre[neighbours[k][0]*len(result[0])+neighbours[k][1]] = []int{node[0], node[1]}
				distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] = alt

				heap.Push(&pq, newItem)

			}

		}

		if node[0] == pos3 && node[1] == pos4 {

			return pre, distances, node, counterPOPS

		}

	}

	return pre, distances, node, counterPOPS

}

//Dijkstra bidirectional implementation
func Dijkstra_bi(result [][]bool, pos1, pos2, pos3, pos4 int) ([][]int, []float64, [][]int, []float64, [2]int, int) {

	var counterPOPS = 0

	distances := make([]float64, len(result)*len(result[0]))

	distances2 := make([]float64, len(result)*len(result[0]))

	pq := make(PriorityQueue, 0)
	pq2 := make(PriorityQueue, 0)
	pre := make([][]int, len(result)*len(result[0]))
	pre2 := make([][]int, len(result)*len(result[0]))

	var alreadyVisited = make(map[[2]int]bool)
	var alreadyVisited2 = make(map[[2]int]bool)
	var bothVisited = false
	var bothVisited2 = false

	var found = false

	var bestMeetingPointDistance = math.MaxFloat64
	var bestMeetingPoint [2]int

	for i := 0; i < len(result); i++ {

		pre[i] = make([]int, len(result[0]))
		pre2[i] = make([]int, len(result[0]))

	}

	for i := 0; i <= len(result)-1; i++ {

		for j := 0; j <= len(result[i])-1; j++ {

			distances[i*len(result[0])+j] = math.MaxFloat64
			distances2[i*len(result[0])+j] = math.MaxFloat64

			pre[i*(len(result[0]))+j] = []int{-1, -1}
			pre2[i*(len(result[0]))+j] = []int{-1, -1}

			if i == pos1 && j == pos2 {
				// Insert a new item and then modify its priority.
				item := &Item{
					pos:              [2]int{pos1, pos2},
					distancePriority: 0,
				}
				heap.Push(&pq, item)

			}

			if i == pos3 && j == pos4 {
				// Insert a new item and then modify its priority.
				item := &Item{
					pos:              [2]int{pos3, pos4},
					distancePriority: 0,
				}
				heap.Push(&pq2, item)

			}
		}
	}

	distances[pos1*(len(result[0]))+pos2] = 0
	distances2[pos3*(len(result[0]))+pos4] = 0

	var node [2]int
	var node2 [2]int

	for pq.Len() > 0 || pq2.Len() > 0 {

		if pq.Len() > 0 {
			p := heap.Pop(&pq).(*Item)

			node = p.pos

			if pre[node[0]*len(result[0])+node[1]][0] != -1 && p.distancePriority > distances[node[0]*len(result[0])+node[1]] {

				continue
			}
			counterPOPS++

			bothVisited = alreadyVisited2[[2]int{node[0], node[1]}]
			alreadyVisited[[2]int{node[0], node[1]}] = true

			var neighbours = GetEdges(result, node[0], node[1])

			for k := 0; k <= len(neighbours)-1; k++ {

				var distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node[0], node[1]), GetCordsFromArrayPosition(len(result), len(result[0]), neighbours[k][0], neighbours[k][1]))

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] && !found {

					newItem := &Item{pos: neighbours[k], distancePriority: alt}
					pre[neighbours[k][0]*len(result[0])+neighbours[k][1]] = []int{node[0], node[1]}
					distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] = alt

					heap.Push(&pq, newItem)

				}

			}
		}
		if pq2.Len() > 0 {
			p2 := heap.Pop(&pq2).(*Item)

			node2 = p2.pos

			if pre2[node2[0]*len(result[0])+node2[1]][0] != -1 && p2.distancePriority > distances2[node2[0]*len(result[0])+node2[1]] {

				continue
			}

			counterPOPS++
			bothVisited2 = alreadyVisited[[2]int{node2[0], node2[1]}]
			alreadyVisited2[[2]int{node2[0], node2[1]}] = true

			var neighbours = GetEdges(result, node2[0], node2[1])

			for k := 0; k <= len(neighbours)-1; k++ {

				var distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node2[0], node2[1]), GetCordsFromArrayPosition(len(result), len(result[0]), neighbours[k][0], neighbours[k][1]))

				var alt = distances2[node2[0]*len(result[0])+node2[1]] + distance

				if alt < distances2[neighbours[k][0]*len(result[0])+neighbours[k][1]] && !found {

					newItem := &Item{pos: neighbours[k], distancePriority: alt}
					pre2[neighbours[k][0]*len(result[0])+neighbours[k][1]] = []int{node2[0], node2[1]}
					distances2[neighbours[k][0]*len(result[0])+neighbours[k][1]] = alt

					heap.Push(&pq2, newItem)

				}

			}
		}

		if bothVisited || bothVisited2 || node == node2 {

			found = true
			if distances[node[0]*len(result[0])+node[1]]+distances2[node[0]*len(result[0])+node[1]] < bestMeetingPointDistance {

				bestMeetingPointDistance = distances[node[0]*len(result[0])+node[1]] + distances2[node[0]*len(result[0])+node[1]]
				bestMeetingPoint = node
			}
			if distances[node2[0]*len(result[0])+node2[1]]+distances2[node2[0]*len(result[0])+node2[1]] < bestMeetingPointDistance {

				bestMeetingPointDistance = distances[node2[0]*len(result[0])+node2[1]] + distances2[node2[0]*len(result[0])+node2[1]]
				bestMeetingPoint = node2
			}

			if (bothVisited || bothVisited2) && pq.Len() == 0 && pq2.Len() == 0 {

				return pre, distances, pre2, distances2, bestMeetingPoint, counterPOPS
			}
		}

	}

	return pre, distances, pre2, distances2, bestMeetingPoint, counterPOPS

}

//get edges to neighbours for a specific point in the bitarray
func GetEdges(result [][]bool, laInt, lnInt int) [][2]int {

	var edges [][2]int

	var m = -1
	//oben links
	if laInt-1 >= 0 {

		m = ModLikePython(lnInt-1, len(result[laInt-1])-1)
	}
	if laInt-1 >= 0 && result[laInt-1][m] == true {

		edges = append(edges, [2]int{laInt - 1, m})
	}
	//oben
	if laInt-1 >= 0 && result[laInt-1][lnInt] == true {

		edges = append(edges, [2]int{laInt - 1, lnInt})
	}

	//oben rechts
	if laInt-1 >= 0 {

		m = ModLikePython(lnInt+1, len(result[laInt-1])-1)
	}
	if laInt-1 >= 0 && result[laInt-1][m] == true {

		edges = append(edges, [2]int{laInt - 1, m})
	}

	//links

	m = ModLikePython(lnInt-1, len(result[laInt])-1)

	if result[laInt][m] == true {

		edges = append(edges, [2]int{laInt, m})
	}

	//rechts

	m = ModLikePython(lnInt+1, len(result[laInt])-1)

	if result[laInt][m] == true {

		edges = append(edges, [2]int{laInt, m})
	}
	//unten links
	if laInt+1 <= len(result)-1 {

		m = ModLikePython(lnInt-1, len(result[laInt+1])-1)
	}
	if laInt+1 <= len(result)-1 && result[laInt+1][m] == true {

		edges = append(edges, [2]int{laInt + 1, m})
	}
	//unten
	if laInt+1 <= len(result)-1 && result[laInt+1][lnInt] == true {

		edges = append(edges, [2]int{laInt + 1, lnInt})
	}

	//unten rechts
	if laInt+1 <= len(result)-1 {

		m = ModLikePython(lnInt+1, len(result[laInt+1])-1)
	}
	if laInt+1 <= len(result)-1 && result[laInt+1][m] == true {

		edges = append(edges, [2]int{laInt + 1, m})
	}

	return edges
}

//Dijkstra implementation optimized without precalculation of the distances between the points of the optimization squares
func Dijkstra_single_optimized(result [][]bool, pos1, pos2, pos3, pos4 int, mapPointSquares map[[2]int]int, optEdges [][2][2]int) ([][]int, []float64, [2]int, int) {

	var counterPOPS = 0

	distances := make([]float64, len(result)*len(result[0]))

	visitedPointsInSquare := make([]bool, len(result)*len(result[0]))

	pq := make(PriorityQueue, 0)
	pre := make([][]int, len(result)*len(result[0]))

	for i := 0; i < len(result); i++ {

		pre[i] = make([]int, len(result[0]))

	}

	for i := 0; i <= len(result)-1; i++ {

		for j := 0; j <= len(result[i])-1; j++ {

			distances[i*len(result[0])+j] = math.MaxFloat64

			visitedPointsInSquare[i*len(result[0])+j] = false

			pre[i*(len(result[0]))+j] = []int{-1, -1}

			if i == pos1 && j == pos2 {
				// Insert a new item and then modify its priority.
				item := &Item{
					pos:              [2]int{pos1, pos2},
					distancePriority: 0,
				}
				heap.Push(&pq, item)

			}
		}
	}

	distances[pos1*(len(result[0]))+pos2] = 0

	var node [2]int

	for pq.Len() > 0 {

		// Find the nearest yet to visit node
		p := heap.Pop(&pq).(*Item)

		node = p.pos

		var isOptimized = p.optimized
		if pre[node[0]*len(result[0])+node[1]][0] != -1 && (p.distancePriority) > distances[node[0]*len(result[0])+node[1]] {

			continue
		}

		if visitedPointsInSquare[node[0]*len(result[0])+node[1]] != false {

			continue
		}

		counterPOPS++

		var k = -1
		var ok = false
		k, ok = mapPointSquares[[2]int{node[0], node[1]}]

		//node within square and dest pos not within square
		if !isOptimized && ok && node[0] >= optEdges[k][0][0] && node[1] >= optEdges[k][0][1] && node[0] <= optEdges[k][1][0] && node[1] <= optEdges[k][1][1] && ((pos3 < optEdges[k][0][0] || pos4 < optEdges[k][0][1]) || (pos3 > optEdges[k][1][0] || pos4 > optEdges[k][1][1])) {

			for xAxis := optEdges[k][0][1]; xAxis <= optEdges[k][1][1]; xAxis++ {

				//upper line
				var distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node[0], node[1]), GetCordsFromArrayPosition(len(result), len(result[0]), optEdges[k][0][0], xAxis))

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[optEdges[k][0][0]*len(result[0])+xAxis] {
					newItem := &Item{pos: [2]int{optEdges[k][0][0], xAxis}, distancePriority: alt, optimized: true}
					pre[optEdges[k][0][0]*len(result[0])+xAxis] = []int{node[0], node[1]}
					distances[optEdges[k][0][0]*len(result[0])+xAxis] = alt

					heap.Push(&pq, newItem)

				}

				//lower line

				distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node[0], node[1]), GetCordsFromArrayPosition(len(result), len(result[0]), optEdges[k][1][0], xAxis))

				alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[optEdges[k][1][0]*len(result[0])+xAxis] {
					newItem2 := &Item{pos: [2]int{optEdges[k][1][0], xAxis}, distancePriority: alt, optimized: true}
					pre[optEdges[k][1][0]*len(result[0])+xAxis] = []int{node[0], node[1]}
					distances[optEdges[k][1][0]*len(result[0])+xAxis] = alt

					heap.Push(&pq, newItem2)

				}

			}

			//left right
			for yAxis := optEdges[k][0][0]; yAxis <= optEdges[k][1][0]; yAxis++ {

				var distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node[0], node[1]), GetCordsFromArrayPosition(len(result), len(result[0]), yAxis, optEdges[k][0][1]))

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[(yAxis)*len(result[0])+optEdges[k][0][1]] {
					newItem := &Item{pos: [2]int{yAxis, optEdges[k][0][1]}, distancePriority: alt, optimized: true}
					pre[(yAxis)*len(result[0])+optEdges[k][0][1]] = []int{node[0], node[1]}
					distances[(yAxis)*len(result[0])+optEdges[k][0][1]] = alt

					heap.Push(&pq, newItem)

				}

				//right line
				distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node[0], node[1]), GetCordsFromArrayPosition(len(result), len(result[0]), yAxis, optEdges[k][1][1]))

				alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[(yAxis)*len(result[0])+optEdges[k][1][1]] {
					newItem2 := &Item{pos: [2]int{yAxis, optEdges[k][1][1]}, distancePriority: alt, optimized: true}
					pre[(yAxis)*len(result[0])+optEdges[k][1][1]] = []int{node[0], node[1]}
					distances[(yAxis)*len(result[0])+optEdges[k][1][1]] = alt

					heap.Push(&pq, newItem2)

				}
			}

			//set nodes within negative
			for yAxis := optEdges[k][0][0] + 1; yAxis < optEdges[k][1][0]; yAxis++ {
				for xAxis := optEdges[k][0][1] + 1; xAxis < optEdges[k][1][1]; xAxis++ {

					visitedPointsInSquare[yAxis*len(result[0])+xAxis] = true
				}
			}

		}

		var neighbours = GetEdges(result, node[0], node[1])

	neig:
		for k := 0; k <= len(neighbours)-1; k++ {

			if visitedPointsInSquare[neighbours[k][0]*len(result[0])+neighbours[k][1]] == true {

				continue neig
			}

			var distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node[0], node[1]), GetCordsFromArrayPosition(len(result), len(result[0]), neighbours[k][0], neighbours[k][1]))

			var alt = distances[node[0]*len(result[0])+node[1]] + distance

			if alt < distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] {
				newItem := &Item{pos: neighbours[k], distancePriority: alt}
				pre[neighbours[k][0]*len(result[0])+neighbours[k][1]] = []int{node[0], node[1]}
				distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] = alt

				heap.Push(&pq, newItem)

			}

		}

		if node[0] == pos3 && node[1] == pos4 {

			return pre, distances, node, counterPOPS

		}

	}

	return pre, distances, node, counterPOPS

}

//A star implementation optimized with precalculation of the distances between the points of the optimization squares
func A_stern_single_optimized_with_pre(result [][]bool, pos1, pos2, pos3, pos4 int, mapPointSquares map[[2]int]int, optEdges [][2][2]int, dists [][][]float64) ([][]int, []float64, [2]int, int) {

	var visitedSquares = map[int]bool{}

	var counterPOPS = 0
	var optimizedAmount = 0
	distances := make([]float64, len(result)*len(result[0]))

	visitedPointsInSquare := make([]bool, len(result)*len(result[0]))

	pq := make(PriorityQueue, 0)
	pre := make([][]int, len(result)*len(result[0]))

	for i := 0; i < len(result); i++ {

		pre[i] = make([]int, len(result[0]))

	}

	for i := 0; i <= len(result)-1; i++ {

		for j := 0; j <= len(result[i])-1; j++ {

			distances[i*len(result[0])+j] = math.MaxFloat64

			visitedPointsInSquare[i*len(result[0])+j] = false

			pre[i*(len(result[0]))+j] = []int{-1, -1}

			if i == pos1 && j == pos2 {
				// Insert a new item and then modify its priority.
				item := &Item{
					pos:                    [2]int{pos1, pos2},
					distancePriority:       GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), pos1, pos2), GetCordsFromArrayPosition(len(result), len(result[0]), pos3, pos4)),
					actualDistanceForAStar: 0,
				}
				heap.Push(&pq, item)

			}
		}
	}

	distances[pos1*(len(result[0]))+pos2] = 0

	var node [2]int

	for pq.Len() > 0 {

		// Find the nearest yet to visit node
		p := heap.Pop(&pq).(*Item)

		node = p.pos
		var isOptimized = p.optimized

		if pre[node[0]*len(result[0])+node[1]][0] != -1 && p.actualDistanceForAStar > distances[node[0]*len(result[0])+node[1]] {

			continue
		}

		counterPOPS++

		var k = -1
		var ok = false
		k, ok = mapPointSquares[[2]int{node[0], node[1]}]

		//node within square and dest pos not within square
		if !isOptimized && ok && node[0] >= optEdges[k][0][0] && node[1] >= optEdges[k][0][1] && node[0] <= optEdges[k][1][0] && node[1] <= optEdges[k][1][1] && ((pos3 < optEdges[k][0][0] || pos4 < optEdges[k][0][1]) || (pos3 > optEdges[k][1][0] || pos4 > optEdges[k][1][1])) && (node[0] == optEdges[k][0][0] || node[1] == optEdges[k][0][1] || node[0] == optEdges[k][1][0] || node[1] == optEdges[k][1][1]) && ((pos1 < optEdges[k][0][0] || pos2 < optEdges[k][0][1]) || (pos1 > optEdges[k][1][0] || pos2 > optEdges[k][1][1])) {

			if !visitedSquares[k] {
				optimizedAmount += (optEdges[k][1][1] - optEdges[k][0][1]) * (optEdges[k][1][0] - optEdges[k][0][0])
			}
			visitedSquares[k] = true

			var zahl = 0

			var foundSquarePoint = false
			//upper
			if !foundSquarePoint && node[1] >= optEdges[k][0][1] && node[1] <= optEdges[k][1][1] && node[0] == optEdges[k][0][0] {

				zahl = zahl + node[1] - optEdges[k][0][1]
				foundSquarePoint = true
			} else if !foundSquarePoint {
				zahl = zahl + optEdges[k][1][1] - optEdges[k][0][1] + 1
			}

			//lower
			if !foundSquarePoint && node[1] >= optEdges[k][0][1] && node[1] <= optEdges[k][1][1] && node[0] == optEdges[k][1][0] {

				zahl = zahl + node[1] - optEdges[k][0][1]
				foundSquarePoint = true
			} else if !foundSquarePoint {
				zahl = zahl + optEdges[k][1][1] - optEdges[k][0][1] + 1
			}

			//left side
			if !foundSquarePoint && node[0] >= optEdges[k][0][0] && node[0] <= optEdges[k][1][0] && node[1] == optEdges[k][0][1] {

				zahl = zahl + node[0] - optEdges[k][0][0]
				foundSquarePoint = true
			} else if !foundSquarePoint {

				zahl = zahl + optEdges[k][1][0] - optEdges[k][0][0] + 1
			}

			//right side
			if !foundSquarePoint && node[0] >= optEdges[k][0][0] && node[0] <= optEdges[k][1][0] && node[1] == optEdges[k][1][1] {

				zahl = zahl + node[0] - optEdges[k][0][0]
				foundSquarePoint = true
			} else if !foundSquarePoint {

				zahl = zahl + optEdges[k][1][0] - optEdges[k][0][0]

			}

			var help = dists[k][zahl]

			for i := 0; i <= optEdges[k][1][1]-optEdges[k][0][1]; i++ {

				var distance = help[i]

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[optEdges[k][0][0]*len(result[0])+(optEdges[k][0][1]+i)] {

					var altHelp = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), optEdges[k][0][0], optEdges[k][0][1]+i), GetCordsFromArrayPosition(len(result), len(result[0]), pos3, pos4))

					newItem := &Item{pos: [2]int{optEdges[k][0][0], (optEdges[k][0][1] + i)}, actualDistanceForAStar: alt, optimized: true, distancePriority: alt + altHelp}
					pre[optEdges[k][0][0]*len(result[0])+(optEdges[k][0][1]+i)] = []int{node[0], node[1]}
					distances[optEdges[k][0][0]*len(result[0])+(optEdges[k][0][1]+i)] = alt

					heap.Push(&pq, newItem)

				}

			}

			for i := 0; i <= optEdges[k][1][1]-optEdges[k][0][1]; i++ {

				var distance = help[1+optEdges[k][1][1]-optEdges[k][0][1]+i]

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[optEdges[k][1][0]*len(result[0])+i+optEdges[k][0][1]] {

					var altHelp = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), optEdges[k][1][0], optEdges[k][0][1]+i), GetCordsFromArrayPosition(len(result), len(result[0]), pos3, pos4))

					newItem := &Item{pos: [2]int{optEdges[k][1][0], optEdges[k][0][1] + i}, actualDistanceForAStar: alt, optimized: true, distancePriority: alt + altHelp}
					pre[optEdges[k][1][0]*len(result[0])+optEdges[k][0][1]+i] = []int{node[0], node[1]}
					distances[optEdges[k][1][0]*len(result[0])+optEdges[k][0][1]+i] = alt

					heap.Push(&pq, newItem)

				}

			}

			//left
			for i := 0; i <= optEdges[k][1][1]-optEdges[k][0][1]; i++ {

				var distance = help[2+2*(optEdges[k][1][1]-optEdges[k][0][1])+i]

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][0][1]] {

					var altHelp = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), optEdges[k][0][0]+i, optEdges[k][0][1]), GetCordsFromArrayPosition(len(result), len(result[0]), pos3, pos4))

					newItem := &Item{pos: [2]int{(optEdges[k][0][0] + i), optEdges[k][0][1]}, actualDistanceForAStar: alt, optimized: true, distancePriority: alt + altHelp}
					pre[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][0][1]] = []int{node[0], node[1]}
					distances[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][0][1]] = alt

					heap.Push(&pq, newItem)

				}

			}

			//right
			for i := 0; i <= optEdges[k][1][1]-optEdges[k][0][1]; i++ {

				var distance = help[3+3*(optEdges[k][1][1]-optEdges[k][0][1])+i]

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][1][1]] {

					var altHelp = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), optEdges[k][0][0]+i, optEdges[k][1][1]), GetCordsFromArrayPosition(len(result), len(result[0]), pos3, pos4))

					newItem := &Item{pos: [2]int{(optEdges[k][0][0] + i), optEdges[k][1][1]}, actualDistanceForAStar: alt, optimized: true, distancePriority: alt + altHelp}
					pre[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][1][1]] = []int{node[0], node[1]}
					distances[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][1][1]] = alt

					heap.Push(&pq, newItem)

				}

			}

			//set nodes within negative
			for yAxis := optEdges[k][0][0] + 1; yAxis < optEdges[k][1][0]; yAxis++ {
				for xAxis := optEdges[k][0][1] + 1; xAxis < optEdges[k][1][1]; xAxis++ {

					visitedPointsInSquare[yAxis*len(result[0])+xAxis] = true
				}
			}

		}

		//destination point within square
		if ok && node[0] >= optEdges[k][0][0] && node[1] >= optEdges[k][0][1] && node[0] <= optEdges[k][1][0] && node[1] <= optEdges[k][1][1] && (pos3 >= optEdges[k][0][0] && pos4 >= optEdges[k][0][1] && pos3 <= optEdges[k][1][0] && pos4 <= optEdges[k][1][1]) {

			var distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node[0], node[1]), GetCordsFromArrayPosition(len(result), len(result[0]), pos3, pos4))

			var alt = distances[node[0]*len(result[0])+node[1]] + distance

			if alt <= distances[pos3*len(result[0])+pos4] {

				pre[pos3*len(result[0])+pos4] = []int{node[0], node[1]}
				distances[pos3*len(result[0])+pos4] = alt

			}

			return pre, distances, node, counterPOPS
		}

		var neighbours = GetEdges(result, node[0], node[1])

	neig:
		for k := 0; k <= len(neighbours)-1; k++ {

			if visitedPointsInSquare[neighbours[k][0]*len(result[0])+neighbours[k][1]] == true {

				continue neig
			}

			var distance = GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), node[0], node[1]), GetCordsFromArrayPosition(len(result), len(result[0]), neighbours[k][0], neighbours[k][1]))

			var alt = distances[node[0]*len(result[0])+node[1]] + distance

			if alt < distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] {
				newItem := &Item{pos: neighbours[k], actualDistanceForAStar: alt, distancePriority: alt + GreatCircleDistance(GetCordsFromArrayPosition(len(result), len(result[0]), neighbours[k][0], neighbours[k][1]), GetCordsFromArrayPosition(len(result), len(result[0]), pos3, pos4))}
				pre[neighbours[k][0]*len(result[0])+neighbours[k][1]] = []int{node[0], node[1]}
				distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] = alt

				heap.Push(&pq, newItem)

			}

		}

		if node[0] == pos3 && node[1] == pos4 {

			return pre, distances, node, counterPOPS

		}

	}

	return pre, distances, node, counterPOPS

}
