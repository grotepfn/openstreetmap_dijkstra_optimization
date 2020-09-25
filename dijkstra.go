package main

import (
	"container/heap"
	"math"
)

// An Item is something we manage in a priority queue.
type Item struct {
	pos       [2]int  // The value of the item; arbitrary.
	distance  float64 // The priority of the item in the queue.
	distance2 float64 //dijkstra
	// The index is needed by update and is maintained by the heap.Interface methods.
	index     int // The index of the item in the heap.
	optimized bool
}

// An Item is something we manage in a priority queue.
type ItemOpt struct {
	pos1  [2]int // The value of the item; arbitrary.
	pos2  [2]int // The value of the item; arbitrary.
	index int    // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

type PriorityQueue2 []*ItemOpt

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue2) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].distance > pq[j].distance
}

func (pq PriorityQueue2) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].pos1[0] > pq[j].pos1[0] && pq[i].pos1[1] > pq[j].pos1[1]
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq PriorityQueue2) Swap(i, j int) {
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

func (pq *PriorityQueue2) Push(x interface{}) {
	n := len(*pq)
	ItemOpt := x.(*ItemOpt)
	ItemOpt.index = n
	*pq = append(*pq, ItemOpt)
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

func (pq *PriorityQueue2) Pop() interface{} {
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
	item.distance = distance
	heap.Fix(pq, item.index)
}

func (pq *PriorityQueue2) update(item *ItemOpt, pos1 [2]int, pos2 [2]int) {
	item.pos1 = pos1
	item.pos2 = pos2
	heap.Fix(pq, item.index)
}

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.

//func dijktstra_minheap(pos1, pos2, pos3, pos4 int) ([][][]int, [][]float64) {
//}

//https://dev.to/douglasmakey/implementation-of-dijkstra-using-heap-in-go-6e3
func dijkstra_minheap(pos1, pos2, pos3, pos4 int) ([][][]int, [][]float64, [][][]int, [][]float64, [2]int) {

	//dijkstra_single(pos1, pos2, pos3, pos4)

	var alreadyVisited = make(map[[2]int]bool)
	var alreadyVisited2 = make(map[[2]int]bool)
	var bothVisited = false
	var bothVisited2 = false
	var bestU = math.MaxFloat64
	var bestUp [2]int

	distances := make([][]float64, len(result))

	pq := make(PriorityQueue, 0)
	pre := make([][][]int, len(result))

	distances2 := make([][]float64, len(result))

	pq2 := make(PriorityQueue, 0)
	pre2 := make([][][]int, len(result))

	//var dist float64 = math.MaxFloat64
	var found = false

	//destiny := [2]int{pos3, pos4}

	for i := 0; i < len(result); i++ {
		distances[i] = make([]float64, len(result[0]))

		pre[i] = make([][]int, len(result[0]))

		distances2[i] = make([]float64, len(result[0]))

		pre2[i] = make([][]int, len(result[0]))
	}

	for i := 0; i <= len(result)-1; i++ {
		for j := 0; j <= len(result[i])-1; j++ {

			distances[i][j] = math.MaxFloat64

			distances2[i][j] = math.MaxFloat64

			pre[i][j] = []int{-1, -1}
			pre2[i][j] = []int{-1, -1}

			if i == pos1 && j == pos2 {
				// Insert a new item and then modify its priority.
				item := &Item{
					pos:      [2]int{pos1, pos2},
					distance: 0,
				}
				heap.Push(&pq, item)

			}
			if i == pos3 && j == pos4 {
				// Insert a new item and then modify its priority.
				item := &Item{
					pos:      [2]int{pos3, pos4},
					distance: 0,
				}
				heap.Push(&pq2, item)
			}

		}
	}

	distances[pos1][pos2] = 0
	distances2[pos3][pos4] = 0

	var node [2]int
	var node2 [2]int
	for pq.Len() > 0 || pq2.Len() > 0 {

		if pq.Len() > 0 {
			// Find the nearest yet to visit node
			p := heap.Pop(&pq).(*Item)
			node = p.pos

			if pre[node[0]][node[1]][0] != -1 && p.distance >= distances[pre[node[0]][node[1]][0]][pre[node[0]][node[1]][1]] {
				continue
			}

			//waruuuuum???
			bothVisited = alreadyVisited2[[2]int{node[0], node[1]}]
			alreadyVisited[[2]int{node[0], node[1]}] = true

			var neighbours = getEdgesPosition(node[0], node[1])
			for k := 0; k <= len(neighbours)-1; k++ {

				var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]))

				var alt = distances[node[0]][node[1]] + distance

				newItem := &Item{pos: neighbours[k], distance: -alt}

				if alt < distances[neighbours[k][0]][neighbours[k][1]] && !found {

					pre[neighbours[k][0]][neighbours[k][1]] = []int{node[0], node[1]}
					distances[neighbours[k][0]][neighbours[k][1]] = alt
					heap.Push(&pq, newItem)

				}

			}

		}

		if pq2.Len() > 0 {
			// Find the nearest yet to visit node
			p2 := heap.Pop(&pq2).(*Item)
			node2 = p2.pos

			if pre2[node2[0]][node2[1]][0] != -1 && p2.distance >= distances2[pre2[node2[0]][node2[1]][0]][pre2[node2[0]][node2[1]][1]] {
				continue
			}

			bothVisited2 = alreadyVisited[[2]int{node2[0], node2[1]}]

			alreadyVisited2[[2]int{node2[0], node2[1]}] = true

			var neighbours2 = getEdgesPosition(node2[0], node2[1])
			for k := 0; k <= len(neighbours2)-1; k++ {

				var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node2[0], node2[1]), getCordsFromArrayPosition(result, neighbours2[k][0], neighbours2[k][1]))

				var alt = distances[node2[0]][node2[1]] + distance

				newItem := &Item{pos: neighbours2[k], distance: -alt}

				if alt < distances[neighbours2[k][0]][neighbours2[k][1]] && !found {

					pre2[neighbours2[k][0]][neighbours2[k][1]] = []int{node2[0], node2[1]}
					distances[neighbours2[k][0]][neighbours2[k][1]] = alt
					heap.Push(&pq, newItem)

				}

			}
		}

		if bothVisited || bothVisited2 || node == node2 {

			found = true
			if distances[node[0]][node[1]]+distances2[node[0]][node[1]] < bestU {
				bestU = distances[node[0]][node[1]] + distances2[node[0]][node[1]]
				bestUp = node
			}
			if distances[node2[0]][node2[1]]+distances2[node2[0]][node2[1]] < bestU {
				bestU = distances[node2[0]][node2[1]] + distances2[node2[0]][node2[1]]
				bestUp = node2
			}

			if bothVisited && pq.Len() == 0 && pq2.Len() == 0 {

				//println(distances[bestUp[0]][bestUp[1]] + distances2[bestUp[0]][bestUp[1]])
				return pre, distances, pre2, distances2, bestUp
			}

		}

	}

	if bothVisited {

		return pre, distances, pre2, distances2, bestUp
	} else {

		return pre, distances, pre2, distances2, bestUp
	}

}

func dijkstra_single(pos1, pos2, pos3, pos4 int) ([][]int, []float64, [2]int, int) {

	var counterPOPS = 0

	distances := make([]float64, len(result)*len(result[0]))
	suchraum := make([]int, len(result)*len(result[0]))
	pq := make(PriorityQueue, 0)
	pre := make([][]int, len(result)*len(result[0]))

	for i := 0; i < len(result); i++ {

		pre[i] = make([]int, len(result[0]))

	}

	for i := 0; i <= len(result)-1; i++ {
		for j := 0; j <= len(result[i])-1; j++ {

			distances[i*(len(result[0]))+j] = math.MaxFloat64
			suchraum[i*(len(result[0]))+j] = 0
			pre[i*(len(result[0]))+j] = []int{-1, -1}

			if i == pos1 && j == pos2 {
				// Insert a new item and then modify its priority.
				item := &Item{
					pos:      [2]int{pos1, pos2},
					distance: 0,
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
		counterPOPS++
		node = p.pos
		var x = suchraum[node[0]*len(result[0])+node[1]]
		x++
		suchraum[node[0]*len(result[0])+node[1]] = x
		if pre[node[0]*len(result[0])+node[1]][0] != -1 && -p.distance > distances[node[0]*len(result[0])+node[1]] {
			//	println("hi")
			continue
		}

		var neighbours = getEdgesPosition(node[0], node[1])
		for k := 0; k <= len(neighbours)-1; k++ {

			var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]))

			var alt = distances[node[0]*(len(result[0]))+node[1]] + distance

			if alt < distances[neighbours[k][0]*(len(result[0]))+neighbours[k][1]] {
				newItem := &Item{pos: neighbours[k], distance: -alt}
				pre[neighbours[k][0]*(len(result[0]))+neighbours[k][1]] = []int{node[0], node[1]}
				distances[neighbours[k][0]*(len(result[0]))+neighbours[k][1]] = alt
				heap.Push(&pq, newItem)

			}

		}

		if node[0] == pos3 && node[1] == pos4 {

			for i := 0; i <= len(result)-1; i = i + 10 {
				for j := 0; j <= len(result[i])-1; j = j + 10 {

					if suchraum[i*len(result[0])+j] > 0 {
						//		print(strconv.Itoa(suchraum[i*len(result[0])+j]))

					} else {
						//		print(" ")
					}
				}
				//	println("")
			}

			return pre, distances, node, counterPOPS

		}

	}

	return pre, distances, node, counterPOPS

}

func a_stern_single_optimized(pos1, pos2, pos3, pos4 int, mapPointSquares map[[2]int][]int, optEgdes [][2][2]int) ([][]int, []float64, [2]int, int) {

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
					pos:       [2]int{pos1, pos2},
					distance:  GreatCircleDistance(getCordsFromArrayPosition(result, pos1, pos2), getCordsFromArrayPosition(result, pos3, pos4)),
					distance2: 0,
				}
				heap.Push(&pq, item)

			}
		}
	}

	distances[pos1*(len(result[0]))+pos2] = 0

	var node [2]int
	//var counter = 0

outer2:
	for pq.Len() > 0 {

		// Find the nearest yet to visit node
		p := heap.Pop(&pq).(*Item)
		counterPOPS++
		node = p.pos
		var isOptimized = p.optimized
		suchraum[node[0]*len(result[0])+node[1]] = true

		if pre[node[0]*len(result[0])+node[1]][0] != -1 && -p.distance2 > distances[node[0]*len(result[0])+node[1]] {

			continue
		}

		if visitedPoints[node[0]*len(result[0])+node[1]] == true {
			println("fehler")
			println(node[0])
			println(node[1])
			continue outer2
		}

		var list = mapPointSquares[[2]int{node[0], node[1]}]

		//var bestSquareDist = math.MaxFloat64
		var k = -1

		if len(list) > 0 {
			//println(len(list))
			k = list[0]
		}

		//node within square and dest pos not within square
		//	println(len(mapPointSquares[[2]int{node[0], node[1]}]))

		if !isOptimized && k != -1 && node[0] >= optEgdes[k][0][0] && node[1] >= optEgdes[k][0][1] && node[0] <= optEgdes[k][1][0] && node[1] <= optEgdes[k][1][1] && ((pos3 < optEgdes[k][0][0] || pos4 < optEgdes[k][0][1]) || (pos3 > optEgdes[k][1][0] || pos4 > optEgdes[k][1][1])) && (node[0] == optEgdes[k][0][0] || node[1] == optEgdes[k][0][1] || node[0] == optEgdes[k][1][0] || node[1] == optEgdes[k][1][1]) {
			if !visitedSquares[k] {
				optimizedAmount += (optEgdes[k][1][1] - optEgdes[k][0][1]) * (optEgdes[k][1][0] - optEgdes[k][0][0])
			}
			visitedSquares[k] = true

			for xAxis := optEgdes[k][0][1]; xAxis <= optEgdes[k][1][1]; xAxis++ {

				if visitedPoints[optEgdes[k][0][0]*len(result[0])+xAxis] == true {
					println("fehler11111111")
					println(node[0])
					println(node[1])

				}

				if visitedPoints[optEdges[k][1][0]*len(result[0])+xAxis] == true {
					println("FEEEEHLER")
					continue
				}

				//upper line
				var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, optEgdes[k][0][0], xAxis))

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				var altHelp = GreatCircleDistance(getCordsFromArrayPosition(result, optEgdes[k][0][0], xAxis), getCordsFromArrayPosition(result, pos3, pos4))

				if alt < distances[optEgdes[k][0][0]*len(result[0])+xAxis] {
					newItem := &Item{pos: [2]int{optEgdes[k][0][0], xAxis}, distance: -alt - altHelp, optimized: true, distance2: -alt}
					pre[optEgdes[k][0][0]*len(result[0])+xAxis] = []int{node[0], node[1]}
					distances[optEgdes[k][0][0]*len(result[0])+xAxis] = alt

					heap.Push(&pq, newItem)
					if visitedPoints[optEgdes[k][0][0]*len(result[0])+xAxis] == true {
						println("fehler11111111")
						println(node[0])
						println(node[1])

					}
				}

				//lower line

				distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, optEgdes[k][1][0], xAxis))

				alt = distances[node[0]*len(result[0])+node[1]] + distance
				altHelp = GreatCircleDistance(getCordsFromArrayPosition(result, optEgdes[k][1][0], xAxis), getCordsFromArrayPosition(result, pos3, pos4))

				if alt < distances[optEgdes[k][1][0]*len(result[0])+xAxis] {
					newItem2 := &Item{pos: [2]int{optEgdes[k][1][0], xAxis}, distance: -alt - altHelp, optimized: true, distance2: -alt}
					pre[optEgdes[k][1][0]*len(result[0])+xAxis] = []int{node[0], node[1]}
					distances[optEgdes[k][1][0]*len(result[0])+xAxis] = alt

					heap.Push(&pq, newItem2)

					if visitedPoints[optEgdes[k][1][0]*len(result[0])+xAxis] == true {
						println("fehler111111112222222222222")
						println(node[0])
						println(node[1])

					}
				}

			}

			//left right
			for yAxis := optEgdes[k][0][0]; yAxis <= optEgdes[k][1][0]; yAxis++ {

				if visitedPoints[(yAxis)*len(result[0])+optEgdes[k][0][1]] == true {
					println("fehler11111111222222222222233333333")
					println(node[0])
					println(node[1])

				}

				if visitedPoints[(yAxis)*len(result[0])+optEgdes[k][1][1]] == true {
					println("fehler11111111222222222222233333333")
					println(node[0])
					println(node[1])

				}

				//left line

				var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, yAxis, optEgdes[k][0][1]))

				var alt = distances[node[0]*len(result[0])+node[1]] + distance
				var altHelp = GreatCircleDistance(getCordsFromArrayPosition(result, yAxis, optEgdes[k][0][1]), getCordsFromArrayPosition(result, pos3, pos4))

				if alt < distances[(yAxis)*len(result[0])+optEgdes[k][0][1]] {
					newItem := &Item{pos: [2]int{yAxis, optEgdes[k][0][1]}, distance: -alt - altHelp, optimized: true, distance2: -alt}
					pre[(yAxis)*len(result[0])+optEgdes[k][0][1]] = []int{node[0], node[1]}
					distances[(yAxis)*len(result[0])+optEgdes[k][0][1]] = alt

					heap.Push(&pq, newItem)

					if visitedPoints[(yAxis)*len(result[0])+optEgdes[k][0][1]] == true {
						println("fehler11111111222222222222233333333")
						println(node[0])
						println(node[1])

					}

				}

				//right line

				distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, yAxis, optEgdes[k][1][1]))

				alt = distances[node[0]*len(result[0])+node[1]] + distance
				altHelp = GreatCircleDistance(getCordsFromArrayPosition(result, yAxis, optEgdes[k][1][1]), getCordsFromArrayPosition(result, pos3, pos4))

				if visitedPoints[(yAxis)*len(result[0])+optEgdes[k][1][1]] == true {
					println("fehler11111111222222222222233333333")
					println(node[0])
					println(node[1])

				}
				if alt < distances[(yAxis)*len(result[0])+optEgdes[k][1][1]] {
					newItem2 := &Item{pos: [2]int{yAxis, optEgdes[k][1][1]}, distance: -alt - altHelp, optimized: true, distance2: -alt}
					pre[(yAxis)*len(result[0])+optEgdes[k][1][1]] = []int{node[0], node[1]}
					distances[(yAxis)*len(result[0])+optEgdes[k][1][1]] = alt

					heap.Push(&pq, newItem2)
					if visitedPoints[(yAxis)*len(result[0])+optEgdes[k][1][1]] == true {
						println("fehler11111111222222222222233333333")
						println(node[0])
						println(node[1])

					}
				}
			}

			//set nodes within negative
			for yAxis := optEgdes[k][0][0] + 1; yAxis < optEgdes[k][1][0]; yAxis++ {
				for xAxis := optEgdes[k][0][1] + 1; xAxis < optEgdes[k][1][1]; xAxis++ {

					//	pre[yAxis*len(result[0])+xAxis] = []int{node[0], node[1]}
					//distances[yAxis*len(result[0])+xAxis] = -3
					visitedPoints[yAxis*len(result[0])+xAxis] = true
				}
			}

			//continue outer2

		}
		/*
			if k != -1 && node[0] >= optEgdes[k][0][0] && node[1] >= optEgdes[k][0][1] && node[0] <= optEgdes[k][1][0] && node[1] <= optEgdes[k][1][1] && (pos3 >= optEgdes[k][0][0] && pos4 >= optEgdes[k][0][1] && pos3 <= optEgdes[k][1][0] && pos4 <= optEgdes[k][1][1]) {

				var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, pos3, pos4))

				var alt = distances[node[0]*len(result[0])+node[1]] + distance
				//var alt2 = alt + GreatCircleDistance(getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]), getCordsFromArrayPosition(result, pos3, pos4))

				if alt <= distances[pos3*len(result[0])+pos4] {

					pre[pos3*len(result[0])+pos4] = []int{node[0], node[1]}
					distances[pos3*len(result[0])+pos4] = alt

				}
				println("optimized amount")
				for i := 0; i <= len(result)-1; i = i + 10 {
					for j := 0; j <= len(result[i])-1; j = j + 10 {

						if suchraum[i*len(result[0])+j] == true {
							print("X")
						} else {
							print(" ")
						}
					}
					println("")
				}
				println("optimized amount")
				println(optimizedAmount)
				return pre, distances, node, counterPOPS
			}
		*/
		var neighbours = getEdgesPosition(node[0], node[1])

	neig:
		for k := 0; k <= len(neighbours)-1; k++ {

			if visitedPoints[neighbours[k][0]*len(result[0])+neighbours[k][1]] == true {

				continue neig
			}

			var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]))

			var alt = distances[node[0]*len(result[0])+node[1]] + distance
			//var alt2 = alt + GreatCircleDistance(getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]), getCordsFromArrayPosition(result, pos3, pos4))

			if alt < distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] {
				newItem := &Item{pos: neighbours[k], distance2: -alt, distance: -alt - GreatCircleDistance(getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]), getCordsFromArrayPosition(result, pos3, pos4))}
				pre[neighbours[k][0]*len(result[0])+neighbours[k][1]] = []int{node[0], node[1]}
				distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] = alt
				if visitedPoints[neighbours[k][0]*len(result[0])+neighbours[k][1]] == true {
					println("fehler11111111222222222222233333333")
					println(node[0])
					println(node[1])

				}
				heap.Push(&pq, newItem)

			}

		}

		if node[0] == pos3 && node[1] == pos4 {

			for i := 0; i <= len(result)-1; i = i + 10 {
				for j := 0; j <= len(result[i])-1; j = j + 10 {

					if suchraum[i*len(result[0])+j] == true {
						//print("X")
					} else {
						//	print(" ")
					}
				}
				//	println("")
			}

			return pre, distances, node, counterPOPS

		}

	}

	return pre, distances, node, counterPOPS

}

func dijkstra_single_optimized(pos1, pos2, pos3, pos4 int, mapPointSquares map[[2]int][]int, optEgdes [][2][2]int, dists [][][]float64) ([][]int, []float64, [2]int, int) {

	var visitedSquares = map[int]bool{}
	suchraum := make([]int, len(result)*len(result[0]))
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

			suchraum[i*len(result[0])+j] = 0

			visitedPoints[i*len(result[0])+j] = false

			pre[i*(len(result[0]))+j] = []int{-1, -1}

			if i == pos1 && j == pos2 {
				// Insert a new item and then modify its priority.
				item := &Item{
					pos:      [2]int{pos1, pos2},
					distance: 0,
				}
				heap.Push(&pq, item)

			}
		}
	}

	distances[pos1*(len(result[0]))+pos2] = 0

	var node [2]int
	//var counter = 0
	var c = 0
outer2:
	for pq.Len() > 0 {

		//println(pq.Len())
		if c%1000 == 0 {
			//	println(pq.Len())

		}
		c++
		// Find the nearest yet to visit node

		p := heap.Pop(&pq).(*Item)

		node = p.pos
		counterPOPS++
		var isOptimized = p.optimized
		if pre[node[0]*len(result[0])+node[1]][0] != -1 && -(p.distance) > distances[node[0]*len(result[0])+node[1]] {
			//	println("hi")
			counterPOPS--
			continue
		}

		if visitedPoints[node[0]*len(result[0])+node[1]] != false {
			counterPOPS--
			println("fehler")
			//println(node[0])
			//println(node[1])

			continue outer2
		}

		var x = suchraum[node[0]*len(result[0])+node[1]]
		x++
		suchraum[node[0]*len(result[0])+node[1]] = x

		var list = mapPointSquares[[2]int{node[0], node[1]}]

		//var bestSquareDist = math.MaxFloat64
		var k = -1

		if len(list) > 0 {
			//println(len(list))
			k = list[0]
		}

		//node within square and dest pos not within square
		if !isOptimized && k != -1 && node[0] >= optEgdes[k][0][0] && node[1] >= optEgdes[k][0][1] && node[0] <= optEgdes[k][1][0] && node[1] <= optEgdes[k][1][1] && ((pos3 < optEgdes[k][0][0] || pos4 < optEgdes[k][0][1]) || (pos3 > optEgdes[k][1][0] || pos4 > optEgdes[k][1][1])) && (node[0] == optEgdes[k][0][0] || node[1] == optEgdes[k][0][1] || node[0] == optEgdes[k][1][0] || node[1] == optEgdes[k][1][1]) && (node[0] == optEgdes[k][0][0] || node[1] == optEgdes[k][0][1] || node[0] == optEgdes[k][1][0] || node[1] == optEgdes[k][1][1]) {

			if !visitedSquares[k] {
				optimizedAmount += (optEgdes[k][1][1] - optEgdes[k][0][1]) * (optEgdes[k][1][0] - optEgdes[k][0][0])
			}
			visitedSquares[k] = true

			var zahl = 0

			var gefunden = false
			//oben
			if !gefunden && node[1] >= optEdges[k][0][1] && node[1] <= optEdges[k][1][1] && node[0] == optEdges[k][0][0] {

				zahl = zahl + node[1] - optEdges[k][0][1]
				gefunden = true
			} else if !gefunden {
				zahl = zahl + optEdges[k][1][1] - optEdges[k][0][1] + 1
			}

			//unten
			if !gefunden && node[1] >= optEdges[k][0][1] && node[1] <= optEdges[k][1][1] && node[0] == optEdges[k][1][0] {

				zahl = zahl + node[1] - optEdges[k][0][1]
				gefunden = true
			} else if !gefunden {
				zahl = zahl + optEdges[k][1][1] - optEdges[k][0][1] + 1
			}

			//links
			if !gefunden && node[0] >= optEdges[k][0][0] && node[0] <= optEdges[k][1][0] && node[1] == optEdges[k][0][1] {

				zahl = zahl + node[0] - optEdges[k][0][0]
				gefunden = true
			} else if !gefunden {
				zahl = zahl + optEdges[k][1][0] - optEdges[k][0][0] + 1
			}

			//rechts
			if !gefunden && node[0] >= optEdges[k][0][0] && node[0] <= optEdges[k][1][0] && node[1] == optEdges[k][1][1] {

				zahl = zahl + node[0] - optEdges[k][0][0]
				gefunden = true
			} else if !gefunden {
				/*	println("ssad")
					println(node[0])
					println(node[1])
					println(optEdges[k][0][0])
					println(optEdges[k][0][0])
					println(optEdges[k][1][0])
					println(optEdges[k][1][1])*/

				zahl = zahl + optEdges[k][1][0] - optEdges[k][0][0]
			}

			var help = dists[k][zahl]

			for i := 0; i <= optEdges[k][1][1]-optEdges[k][0][1]; i++ {

				var distance = help[i]

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[optEgdes[k][0][0]*len(result[0])+(optEdges[k][0][1]+i)] {

					newItem := &Item{pos: [2]int{optEgdes[k][0][0], (optEdges[k][0][1] + i)}, distance: -alt, optimized: true}
					pre[optEgdes[k][0][0]*len(result[0])+(optEdges[k][0][1]+i)] = []int{node[0], node[1]}
					distances[optEgdes[k][0][0]*len(result[0])+(optEdges[k][0][1]+i)] = alt

					heap.Push(&pq, newItem)

				}

			}

			for i := 0; i <= optEdges[k][1][1]-optEdges[k][0][1]; i++ {

				var distance = help[1+optEdges[k][1][1]-optEdges[k][0][1]+i]

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[optEgdes[k][1][0]*len(result[0])+i+optEdges[k][0][1]] {

					newItem := &Item{pos: [2]int{optEgdes[k][1][0], optEdges[k][0][1] + i}, distance: -alt, optimized: true}
					pre[optEgdes[k][1][0]*len(result[0])+optEdges[k][0][1]+i] = []int{node[0], node[1]}
					distances[optEgdes[k][1][0]*len(result[0])+optEdges[k][0][1]+i] = alt

					heap.Push(&pq, newItem)

				}

			}

			//left
			for i := 0; i <= optEdges[k][1][1]-optEdges[k][0][1]; i++ {

				var distance = help[2+2*(optEdges[k][1][1]-optEdges[k][0][1])+i]

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][0][1]] {

					newItem := &Item{pos: [2]int{(optEdges[k][0][0] + i), optEdges[k][0][1]}, distance: -alt, optimized: true}
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

					newItem := &Item{pos: [2]int{(optEdges[k][0][0] + i), optEdges[k][1][1]}, distance: -alt, optimized: true}
					pre[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][1][1]] = []int{node[0], node[1]}
					distances[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][1][1]] = alt

					heap.Push(&pq, newItem)

				}

			}

			/*

				for i := 0; i <= optEdges[k][1][1]-optEdges[k][0][1]; i++ {

					var distance = help[i]

					var alt = distances[node[0]*len(result[0])+node[1]] + distance

					if alt < distances[(optEdges[k][0][1]+i)*len(result[0])+optEdges[k][1][1]] {

						newItem := &Item{pos: [2]int{optEdges[k][0][1] + i, optEdges[k][1][1]}, distance: -alt, optimized: true}
						pre[(optEdges[k][0][1]+i)*len(result[0])+optEdges[k][1][1]] = []int{node[0], node[1]}
						distances[(optEdges[k][0][1]+i)*len(result[0])+optEdges[k][1][1]] = alt

						heap.Push(&pq, newItem)

					}

				}*/

			//set nodes within negative
			for yAxis := optEgdes[k][0][0] + 1; yAxis < optEgdes[k][1][0]; yAxis++ {
				for xAxis := optEgdes[k][0][1] + 1; xAxis < optEgdes[k][1][1]; xAxis++ {

					visitedPoints[yAxis*len(result[0])+xAxis] = true
				}
			}

			//continue outer2

		}
		/*
			if k != -1 && node[0] >= optEgdes[k][0][0] && node[1] >= optEgdes[k][0][1] && node[0] <= optEgdes[k][1][0] && node[1] <= optEgdes[k][1][1] && (pos3 >= optEgdes[k][0][0] && pos4 >= optEgdes[k][0][1] && pos3 <= optEgdes[k][1][0] && pos4 <= optEgdes[k][1][1]) {

				var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, pos3, pos4))

				var alt = distances[node[0]*len(result[0])+node[1]] + distance
				//var alt2 = alt + GreatCircleDistance(getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]), getCordsFromArrayPosition(result, pos3, pos4))

				if alt <= distances[pos3*len(result[0])+pos4] {

					pre[pos3*len(result[0])+pos4] = []int{node[0], node[1]}
					distances[pos3*len(result[0])+pos4] = alt

				}
				println("optimized amount")
				for i := 0; i <= len(result)-1; i = i + 10 {
					for j := 0; j <= len(result[i])-1; j = j + 10 {

						if suchraum[i*len(result[0])+j] == true {
							print("X")
						} else {
							print(" ")
						}
					}
					println("")
				}
				println("optimized amount")
				println(optimizedAmount)
				return pre, distances, node, counterPOPS
			}
		*/
		var neighbours = getEdgesPosition(node[0], node[1])

	neig:
		for k := 0; k <= len(neighbours)-1; k++ {

			if visitedPoints[neighbours[k][0]*len(result[0])+neighbours[k][1]] == true {

				continue neig
			}

			var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]))

			var alt = distances[node[0]*len(result[0])+node[1]] + distance
			//var alt2 = alt + GreatCircleDistance(getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]), getCordsFromArrayPosition(result, pos3, pos4))

			if alt < distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] {
				newItem := &Item{pos: neighbours[k], distance: -alt}
				pre[neighbours[k][0]*len(result[0])+neighbours[k][1]] = []int{node[0], node[1]}
				distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] = alt
				if visitedPoints[neighbours[k][0]*len(result[0])+neighbours[k][1]] == true {
					println("fehler11111111222222222222233333333")
					println(node[0])
					println(node[1])

				}
				heap.Push(&pq, newItem)

			}

		}

		if node[0] == pos3 && node[1] == pos4 {

			for i := 0; i <= len(result)-1; i = i + 10 {
				for j := 0; j <= len(result[i])-1; j = j + 10 {

					if suchraum[i*len(result[0])+j] != 0 {
						//print(suchraum[i*len(result[0])+j])
						//print(" ")
					} else {
						//print(" ")
					}
				}
				//println("")
			}

			return pre, distances, node, counterPOPS

		}

	}

	return pre, distances, node, counterPOPS

}

func a_stern_single(pos1, pos2, pos3, pos4 int, mapPointSquares map[[2]int][]int) ([][]int, []float64, [2]int, int) {

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
					pos:      [2]int{pos1, pos2},
					distance: -GreatCircleDistance(getCordsFromArrayPosition(result, pos1, pos2), getCordsFromArrayPosition(result, pos3, pos4)),
				}
				heap.Push(&pq, item)

			}
		}
	}

	distances[pos1*(len(result[0]))+pos2] = 0

	var node [2]int
	//var counter = 0
	var c = 0
	//outer2:
	for pq.Len() > 0 {

		//println(pq.Len())
		c++
		if c == 0 {
			println(pq.Len())
		}

		p := heap.Pop(&pq).(*Item)
		counterPOPS++
		node = p.pos

		if pre[node[0]*len(result[0])+node[1]][0] != -1 && -(p.distance) > distances[node[0]*len(result[0])+node[1]] {
			/*println("hi")
			println(p.distance)
			println(distances[node[0]*len(result[0])+node[1]])*/
			//continue outer2
		} else {
			///println("not hi")
		}

		var neighbours = getEdgesPosition(node[0], node[1])

		for k := 0; k <= len(neighbours)-1; k++ {

			var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]))

			var alt = distances[node[0]*len(result[0])+node[1]] + distance

			if alt < distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] {

				newItem := &Item{pos: neighbours[k], distance: -alt - GreatCircleDistance(getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]), getCordsFromArrayPosition(result, pos3, pos4))}
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

func a_stern_bi(pos1, pos2, pos3, pos4 int, mapPointSquares map[[2]int][]int) ([][]int, []float64, [][]int, []float64, [2]int, int) {

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
					pos:      [2]int{pos1, pos2},
					distance: -GreatCircleDistance(getCordsFromArrayPosition(result, pos1, pos2), getCordsFromArrayPosition(result, pos3, pos4)),
				}
				heap.Push(&pq, item)

			}

			if i == pos3 && j == pos4 {
				// Insert a new item and then modify its priority.
				item := &Item{
					pos:      [2]int{pos3, pos4},
					distance: -GreatCircleDistance(getCordsFromArrayPosition(result, pos3, pos4), getCordsFromArrayPosition(result, pos1, pos2)),
				}
				heap.Push(&pq2, item)

			}
		}
	}

	distances[pos1*(len(result[0]))+pos2] = 0
	distances2[pos3*(len(result[0]))+pos4] = 0

	var node [2]int
	var node2 [2]int
	//var counter = 0

	for pq.Len() > 0 || pq2.Len() > 0 {

		if pq.Len() > 0 {
			p := heap.Pop(&pq).(*Item)
			counterPOPS++
			node = p.pos

			if pre[node[0]*(len(result[0]))+node[1]][0] != -1 && p.distance >= distances[pre[node[0]*(len(result[0]))+node[1]][0]*len(result[0])+pre[node[0]*(len(result[0]))+node[1]][1]] {
				continue
			}

			bothVisited = alreadyVisited2[[2]int{node[0], node[1]}]
			alreadyVisited[[2]int{node[0], node[1]}] = true

			var neighbours = getEdgesPosition(node[0], node[1])

			for k := 0; k <= len(neighbours)-1; k++ {

				var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]))

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] && !found {

					newItem := &Item{pos: neighbours[k], distance: -alt - GreatCircleDistance(getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]), getCordsFromArrayPosition(result, pos3, pos4))}
					pre[neighbours[k][0]*len(result[0])+neighbours[k][1]] = []int{node[0], node[1]}
					distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] = alt

					heap.Push(&pq, newItem)

				}

			}
		}
		if pq2.Len() > 0 {
			p2 := heap.Pop(&pq2).(*Item)
			counterPOPS++
			node2 = p2.pos

			if pre2[node2[0]*(len(result[0]))+node2[1]][0] != -1 && p2.distance >= distances2[pre2[node2[0]*(len(result[0]))+node2[1]][0]*len(result[0])+pre2[node2[0]*(len(result[0]))+node2[1]][1]] {
				continue
			}

			bothVisited2 = alreadyVisited[[2]int{node2[0], node2[1]}]
			alreadyVisited2[[2]int{node2[0], node2[1]}] = true

			var neighbours = getEdgesPosition(node2[0], node2[1])

			for k := 0; k <= len(neighbours)-1; k++ {

				var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node2[0], node2[1]), getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]))

				var alt = distances2[node2[0]*len(result[0])+node2[1]] + distance

				if alt < distances2[neighbours[k][0]*len(result[0])+neighbours[k][1]] && !found {

					newItem := &Item{pos: neighbours[k], distance: -alt - GreatCircleDistance(getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]), getCordsFromArrayPosition(result, pos1, pos2))}
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
func getShortestPath(desLat int, desLng int, pre [][][]int) [][]int {

	var way [][]int
	way = append(way, []int{desLat, desLng})

	var u []int = []int{desLat, desLng}

	for pre[u[0]][u[1]][0] != -1 {
		//	println("still going")
		u = []int{pre[u[0]][u[1]][0], pre[u[0]][u[1]][1]}
		way = append([][]int{u}, way...)

	}

	return way
}

func getShortestPath2(desLat int, desLng int, pre [][]int) [][]int {

	var way [][]int
	way = append(way, []int{desLat, desLng})

	var u []int = []int{desLat, desLng}

	for pre[(u[0]*(len(result[0])) + u[1])][0] != -1 {

		u = []int{pre[u[0]*(len(result[0]))+u[1]][0], pre[u[0]*(len(result[0]))+u[1]][1]}
		way = append([][]int{u}, way...)

	}

	return way
}

func getEdgesPosition(laInt, lnInt int) [][2]int {

	var edges [][2]int

	var m = -1
	//oben links
	if laInt-1 >= 0 {

		m = modLikePython(lnInt-1, len(result[laInt-1])-1)
	}
	if laInt-1 >= 0 && result[laInt-1][m] == true {

		//	println("oben links")
		edges = append(edges, [2]int{laInt - 1, m})
	}
	//oben
	if laInt-1 >= 0 && result[laInt-1][lnInt] == true {

		edges = append(edges, [2]int{laInt - 1, lnInt})
	}

	//oben rechts
	if laInt-1 >= 0 {

		m = modLikePython(lnInt+1, len(result[laInt-1])-1)
	}
	if laInt-1 >= 0 && result[laInt-1][m] == true {

		edges = append(edges, [2]int{laInt - 1, m})
	}

	//links

	m = modLikePython(lnInt-1, len(result[laInt])-1)

	if result[laInt][m] == true {

		edges = append(edges, [2]int{laInt, m})
	}

	//rechts

	m = modLikePython(lnInt+1, len(result[laInt])-1)

	if result[laInt][m] == true {

		edges = append(edges, [2]int{laInt, m})
	}
	//unten links
	if laInt+1 <= len(result)-1 {

		m = modLikePython(lnInt-1, len(result[laInt+1])-1)
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

		m = modLikePython(lnInt+1, len(result[laInt+1])-1)
	}
	if laInt+1 <= len(result)-1 && result[laInt+1][m] == true {

		edges = append(edges, [2]int{laInt + 1, m})
	}

	return edges
}

//lat lng
func getCordsFromArrayPosition(result [][]bool, pos1 int, pos2 int) [2]float64 {

	return [2]float64{90 - (float64(pos1) / float64(len(result)) * 180), -180 + 360*(float64(pos2)/float64(len(result[0])))}

}

//lat lng
func getArrayPositionFromCords(result [][]bool, lat float64, lng float64) [2]int {

	return [2]int{int(math.Round((lat - 90) / 180 * float64(len(result)) * -1)), int(math.Round((lng + 180) / 360 * float64(len(result[0])-1)))}

}

//https://stackoverflow.com/questions/43018206/modulo-of-negative-integers-in-go
func modLikePython(d, m int) int {
	var res int = d % m
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return res + m
	}
	return res
}

type GeoPoint struct {
	lat float64
	lng float64
}

//https://github.com/kellydunn/golang-geo/blob/master/point.go
// GreatCircleDistance: Calculates the Haversine distance between two points in kilometers.
// Original Implementation from: http://www.movable-type.co.uk/scripts/latlong.html
func GreatCircleDistance(l1 [2]float64, l2 [2]float64) float64 {
	var EARTH_RADIUS = 6371.0
	dLat := (l2[0] - l1[0]) * (math.Pi / 180.0)
	dLon := (l2[1] - l1[1]) * (math.Pi / 180.0)

	lat1 := l1[0] * (math.Pi / 180.0)
	lat2 := l2[0] * (math.Pi / 180.0)

	a1 := math.Sin(dLat/2) * math.Sin(dLat/2)
	a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)

	a := a1 + a2

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EARTH_RADIUS * c
}

//https://www.geeksforgeeks.org/program-for-point-of-intersection-of-two-lines/
func cutting(A, B, C, D GeoPoint) GeoPoint {
	// Line AB represented as a1x + b1y = c1
	var a1 = B.lat - A.lat
	var b1 = A.lng - B.lng
	var c1 = a1*(A.lng) + b1*(A.lat)

	// Line CD represented as a2x + b2y = c2
	var a2 = D.lat - C.lat
	var b2 = C.lng - D.lng
	var c2 = a2*(C.lng) + b2*(C.lat)

	var determinant = a1*b2 - a2*b1

	if determinant == 0 {
		// The lines are parallel. This is simplified
		// by returning a pair of FLT_MAX
		return GeoPoint{math.MaxFloat64, math.MaxFloat64}
	} else {
		var x = (b2*c1 - b1*c2) / determinant
		var y = (a1*c2 - a2*c1) / determinant

		if (y > C.lat && y > D.lat) || (y < C.lat && y < D.lat) || (x > C.lng && x > D.lng) || (x < C.lng && x < D.lng) {
			return GeoPoint{math.MaxFloat64, math.MaxFloat64}
		} else {

			return GeoPoint{y, x}
		}
	}
}

func dijkstra_single_optimized_old(pos1, pos2, pos3, pos4 int, mapPointSquares map[[2]int][]int, optEgdes [][2][2]int) ([][]int, []float64, [2]int, int) {

	var visitedSquares = map[int]bool{}
	suchraum := make([]int, len(result)*len(result[0]))
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

			suchraum[i*len(result[0])+j] = 0

			visitedPoints[i*len(result[0])+j] = false

			pre[i*(len(result[0]))+j] = []int{-1, -1}

			if i == pos1 && j == pos2 {
				// Insert a new item and then modify its priority.
				item := &Item{
					pos:      [2]int{pos1, pos2},
					distance: 0,
				}
				heap.Push(&pq, item)

			}
		}
	}

	distances[pos1*(len(result[0]))+pos2] = 0

	var node [2]int
	//var counter = 0
	var c = 0
outer2:
	for pq.Len() > 0 {

		//println(pq.Len())
		if c%1000 == 0 {
			//	println(pq.Len())

		}
		c++
		// Find the nearest yet to visit node

		p := heap.Pop(&pq).(*Item)

		node = p.pos
		counterPOPS++
		var isOptimized = p.optimized
		if pre[node[0]*len(result[0])+node[1]][0] != -1 && -(p.distance) > distances[node[0]*len(result[0])+node[1]] {
			//	println("hi")
			counterPOPS--
			continue
		}

		if visitedPoints[node[0]*len(result[0])+node[1]] != false {
			counterPOPS--
			println("fehler")
			println(node[0])
			println(node[1])
			continue outer2
		}

		var x = suchraum[node[0]*len(result[0])+node[1]]
		x++
		suchraum[node[0]*len(result[0])+node[1]] = x

		var list = mapPointSquares[[2]int{node[0], node[1]}]

		//var bestSquareDist = math.MaxFloat64
		var k = -1

		if len(list) > 0 {
			//println(len(list))
			k = list[0]
		}

		//node within square and dest pos not within square
		if !isOptimized && k != -1 && node[0] >= optEgdes[k][0][0] && node[1] >= optEgdes[k][0][1] && node[0] <= optEgdes[k][1][0] && node[1] <= optEgdes[k][1][1] && ((pos3 < optEgdes[k][0][0] || pos4 < optEgdes[k][0][1]) || (pos3 > optEgdes[k][1][0] || pos4 > optEgdes[k][1][1])) && (node[0] == optEgdes[k][0][0] || node[1] == optEgdes[k][0][1] || node[0] == optEgdes[k][1][0] || node[1] == optEgdes[k][1][1]) {
			if !visitedSquares[k] {
				optimizedAmount += (optEgdes[k][1][1] - optEgdes[k][0][1]) * (optEgdes[k][1][0] - optEgdes[k][0][0])
			}
			visitedSquares[k] = true

			for xAxis := optEgdes[k][0][1]; xAxis <= optEgdes[k][1][1]; xAxis++ {

				if visitedPoints[optEgdes[k][0][0]*len(result[0])+xAxis] == true {
					println("fehler11111111")
					println(node[0])
					println(node[1])

				}

				if visitedPoints[optEdges[k][1][0]*len(result[0])+xAxis] == true {
					println("FEEEEHLER")
					continue
				}

				//upper line
				var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, optEgdes[k][0][0], xAxis))

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[optEgdes[k][0][0]*len(result[0])+xAxis] {
					newItem := &Item{pos: [2]int{optEgdes[k][0][0], xAxis}, distance: -alt, optimized: true}
					pre[optEgdes[k][0][0]*len(result[0])+xAxis] = []int{node[0], node[1]}
					distances[optEgdes[k][0][0]*len(result[0])+xAxis] = alt

					heap.Push(&pq, newItem)
					if visitedPoints[optEgdes[k][0][0]*len(result[0])+xAxis] == true {
						println("fehler11111111")
						println(node[0])
						println(node[1])

					}
				}

				//lower line

				distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, optEgdes[k][1][0], xAxis))

				alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[optEgdes[k][1][0]*len(result[0])+xAxis] {
					newItem2 := &Item{pos: [2]int{optEgdes[k][1][0], xAxis}, distance: -alt, optimized: true}
					pre[optEgdes[k][1][0]*len(result[0])+xAxis] = []int{node[0], node[1]}
					distances[optEgdes[k][1][0]*len(result[0])+xAxis] = alt

					heap.Push(&pq, newItem2)

					if visitedPoints[optEgdes[k][1][0]*len(result[0])+xAxis] == true {
						println("fehler111111112222222222222")
						println(node[0])
						println(node[1])

					}
				}

			}

			//left right
			for yAxis := optEgdes[k][0][0]; yAxis <= optEgdes[k][1][0]; yAxis++ {

				if visitedPoints[(yAxis)*len(result[0])+optEgdes[k][0][1]] == true {
					println("fehler11111111222222222222233333333")
					println(node[0])
					println(node[1])

				}

				if visitedPoints[(yAxis)*len(result[0])+optEgdes[k][1][1]] == true {
					println("fehler11111111222222222222233333333")
					println(node[0])
					println(node[1])

				}

				var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, yAxis, optEgdes[k][0][1]))

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[(yAxis)*len(result[0])+optEgdes[k][0][1]] {
					newItem := &Item{pos: [2]int{yAxis, optEgdes[k][0][1]}, distance: -alt, optimized: true}
					pre[(yAxis)*len(result[0])+optEgdes[k][0][1]] = []int{node[0], node[1]}
					distances[(yAxis)*len(result[0])+optEgdes[k][0][1]] = alt

					heap.Push(&pq, newItem)

					if visitedPoints[(yAxis)*len(result[0])+optEgdes[k][0][1]] == true {
						println("fehler11111111222222222222233333333")
						println(node[0])
						println(node[1])

					}

				}

				//right line

				distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, yAxis, optEgdes[k][1][1]))

				alt = distances[node[0]*len(result[0])+node[1]] + distance

				if visitedPoints[(yAxis)*len(result[0])+optEgdes[k][1][1]] == true {
					println("fehler11111111222222222222233333333")
					println(node[0])
					println(node[1])

				}
				if alt < distances[(yAxis)*len(result[0])+optEgdes[k][1][1]] {
					newItem2 := &Item{pos: [2]int{yAxis, optEgdes[k][1][1]}, distance: -alt, optimized: true}
					pre[(yAxis)*len(result[0])+optEgdes[k][1][1]] = []int{node[0], node[1]}
					distances[(yAxis)*len(result[0])+optEgdes[k][1][1]] = alt

					heap.Push(&pq, newItem2)
					if visitedPoints[(yAxis)*len(result[0])+optEgdes[k][1][1]] == true {
						println("fehler11111111222222222222233333333")
						println(node[0])
						println(node[1])

					}
				}
			}

			//set nodes within negative
			for yAxis := optEgdes[k][0][0] + 1; yAxis < optEgdes[k][1][0]; yAxis++ {
				for xAxis := optEgdes[k][0][1] + 1; xAxis < optEgdes[k][1][1]; xAxis++ {

					visitedPoints[yAxis*len(result[0])+xAxis] = true
				}
			}

			//continue outer2

		}
		/*
			if k != -1 && node[0] >= optEgdes[k][0][0] && node[1] >= optEgdes[k][0][1] && node[0] <= optEgdes[k][1][0] && node[1] <= optEgdes[k][1][1] && (pos3 >= optEgdes[k][0][0] && pos4 >= optEgdes[k][0][1] && pos3 <= optEgdes[k][1][0] && pos4 <= optEgdes[k][1][1]) {
				var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, pos3, pos4))
				var alt = distances[node[0]*len(result[0])+node[1]] + distance
				//var alt2 = alt + GreatCircleDistance(getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]), getCordsFromArrayPosition(result, pos3, pos4))
				if alt <= distances[pos3*len(result[0])+pos4] {
					pre[pos3*len(result[0])+pos4] = []int{node[0], node[1]}
					distances[pos3*len(result[0])+pos4] = alt
				}
				println("optimized amount")
				for i := 0; i <= len(result)-1; i = i + 10 {
					for j := 0; j <= len(result[i])-1; j = j + 10 {
						if suchraum[i*len(result[0])+j] == true {
							print("X")
						} else {
							print(" ")
						}
					}
					println("")
				}
				println("optimized amount")
				println(optimizedAmount)
				return pre, distances, node, counterPOPS
			}
		*/
		var neighbours = getEdgesPosition(node[0], node[1])

	neig:
		for k := 0; k <= len(neighbours)-1; k++ {

			if visitedPoints[neighbours[k][0]*len(result[0])+neighbours[k][1]] == true {

				continue neig
			}

			var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]))

			var alt = distances[node[0]*len(result[0])+node[1]] + distance
			//var alt2 = alt + GreatCircleDistance(getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]), getCordsFromArrayPosition(result, pos3, pos4))

			if alt < distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] {
				newItem := &Item{pos: neighbours[k], distance: -alt}
				pre[neighbours[k][0]*len(result[0])+neighbours[k][1]] = []int{node[0], node[1]}
				distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] = alt
				if visitedPoints[neighbours[k][0]*len(result[0])+neighbours[k][1]] == true {
					println("fehler11111111222222222222233333333")
					println(node[0])
					println(node[1])

				}
				heap.Push(&pq, newItem)

			}

		}

		if node[0] == pos3 && node[1] == pos4 {

			for i := 0; i <= len(result)-1; i = i + 10 {
				for j := 0; j <= len(result[i])-1; j = j + 10 {

					if suchraum[i*len(result[0])+j] != 0 {
						//print(suchraum[i*len(result[0])+j])
						//print(" ")
					} else {
						//print(" ")
					}
				}
				//println("")
			}

			return pre, distances, node, counterPOPS

		}

	}

	return pre, distances, node, counterPOPS

}

func a_stern_single_optimized_with_pre(pos1, pos2, pos3, pos4 int, mapPointSquares map[[2]int][]int, optEgdes [][2][2]int, dists [][][]float64) ([][]int, []float64, [2]int, int) {

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
					pos:       [2]int{pos1, pos2},
					distance:  GreatCircleDistance(getCordsFromArrayPosition(result, pos1, pos2), getCordsFromArrayPosition(result, pos3, pos4)),
					distance2: 0,
				}
				heap.Push(&pq, item)

			}
		}
	}

	distances[pos1*(len(result[0]))+pos2] = 0

	var node [2]int
	//var counter = 0

outer2:
	for pq.Len() > 0 {

		// Find the nearest yet to visit node
		p := heap.Pop(&pq).(*Item)
		counterPOPS++
		node = p.pos
		var isOptimized = p.optimized
		suchraum[node[0]*len(result[0])+node[1]] = true

		if pre[node[0]*len(result[0])+node[1]][0] != -1 && -p.distance2 > distances[node[0]*len(result[0])+node[1]] {

			continue
		}

		if visitedPoints[node[0]*len(result[0])+node[1]] == true {
			//println("fehler")
			//println(node[0])
			//println(node[1])
			continue outer2
		}

		var list = mapPointSquares[[2]int{node[0], node[1]}]

		//var bestSquareDist = math.MaxFloat64
		var k = -1

		if len(list) > 0 {
			//println(len(list))
			k = list[0]
		}

		//node within square and dest pos not within square

		if !isOptimized && k != -1 && node[0] >= optEgdes[k][0][0] && node[1] >= optEgdes[k][0][1] && node[0] <= optEgdes[k][1][0] && node[1] <= optEgdes[k][1][1] && ((pos3 < optEgdes[k][0][0] || pos4 < optEgdes[k][0][1]) || (pos3 > optEgdes[k][1][0] || pos4 > optEgdes[k][1][1])) && (node[0] == optEgdes[k][0][0] || node[1] == optEgdes[k][0][1] || node[0] == optEgdes[k][1][0] || node[1] == optEgdes[k][1][1]) {
			if !visitedSquares[k] {
				optimizedAmount += (optEgdes[k][1][1] - optEgdes[k][0][1]) * (optEgdes[k][1][0] - optEgdes[k][0][0])
			}
			visitedSquares[k] = true

			var zahl = 0

			var gefunden = false
			//oben
			if !gefunden && node[1] >= optEdges[k][0][1] && node[1] <= optEdges[k][1][1] && node[0] == optEdges[k][0][0] {

				zahl = zahl + node[1] - optEdges[k][0][1]
				gefunden = true
			} else if !gefunden {
				zahl = zahl + optEdges[k][1][1] - optEdges[k][0][1] + 1
			}

			//unten
			if !gefunden && node[1] >= optEdges[k][0][1] && node[1] <= optEdges[k][1][1] && node[0] == optEdges[k][1][0] {

				zahl = zahl + node[1] - optEdges[k][0][1]
				gefunden = true
			} else if !gefunden {
				zahl = zahl + optEdges[k][1][1] - optEdges[k][0][1] + 1
			}

			//links
			if !gefunden && node[0] >= optEdges[k][0][0] && node[0] <= optEdges[k][1][0] && node[1] == optEdges[k][0][1] {

				zahl = zahl + node[0] - optEdges[k][0][0]
				gefunden = true
			} else if !gefunden {

				zahl = zahl + optEdges[k][1][0] - optEdges[k][0][0] + 1
			}

			//rechts
			if !gefunden && node[0] >= optEdges[k][0][0] && node[0] <= optEdges[k][1][0] && node[1] == optEdges[k][1][1] {

				zahl = zahl + node[0] - optEdges[k][0][0]
				gefunden = true
			} else if !gefunden {
				println("ssad")
				println(node[0])
				println(node[1])
				println(optEdges[k][0][0])
				println(optEdges[k][0][0])
				println(optEdges[k][1][0])
				println(optEdges[k][1][1])

				zahl = zahl + optEdges[k][1][0] - optEdges[k][0][0]

			}

			var help = dists[k][zahl]

			for i := 0; i <= optEdges[k][1][1]-optEdges[k][0][1]; i++ {

				var distance = help[i]

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[optEgdes[k][0][0]*len(result[0])+(optEdges[k][0][1]+i)] {

					var altHelp = GreatCircleDistance(getCordsFromArrayPosition(result, optEgdes[k][0][0], optEdges[k][0][1]+i), getCordsFromArrayPosition(result, pos3, pos4))

					newItem := &Item{pos: [2]int{optEgdes[k][0][0], (optEdges[k][0][1] + i)}, distance2: -alt, optimized: true, distance: -alt - altHelp}
					pre[optEgdes[k][0][0]*len(result[0])+(optEdges[k][0][1]+i)] = []int{node[0], node[1]}
					distances[optEgdes[k][0][0]*len(result[0])+(optEdges[k][0][1]+i)] = alt

					heap.Push(&pq, newItem)

				}

			}

			for i := 0; i <= optEdges[k][1][1]-optEdges[k][0][1]; i++ {

				var distance = help[1+optEdges[k][1][1]-optEdges[k][0][1]+i]

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[optEgdes[k][1][0]*len(result[0])+i+optEdges[k][0][1]] {

					var altHelp = GreatCircleDistance(getCordsFromArrayPosition(result, optEgdes[k][1][0], optEdges[k][0][1]+i), getCordsFromArrayPosition(result, pos3, pos4))

					newItem := &Item{pos: [2]int{optEgdes[k][1][0], optEdges[k][0][1] + i}, distance2: -alt, optimized: true, distance: -alt - altHelp}
					pre[optEgdes[k][1][0]*len(result[0])+optEdges[k][0][1]+i] = []int{node[0], node[1]}
					distances[optEgdes[k][1][0]*len(result[0])+optEdges[k][0][1]+i] = alt

					heap.Push(&pq, newItem)

				}

			}

			//left
			for i := 0; i <= optEdges[k][1][1]-optEdges[k][0][1]; i++ {

				var distance = help[2+2*(optEdges[k][1][1]-optEdges[k][0][1])+i]

				var alt = distances[node[0]*len(result[0])+node[1]] + distance

				if alt < distances[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][0][1]] {

					var altHelp = GreatCircleDistance(getCordsFromArrayPosition(result, optEgdes[k][0][0]+i, optEdges[k][0][1]), getCordsFromArrayPosition(result, pos3, pos4))

					newItem := &Item{pos: [2]int{(optEdges[k][0][0] + i), optEdges[k][0][1]}, distance2: -alt, optimized: true, distance: -alt - altHelp}
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

					var altHelp = GreatCircleDistance(getCordsFromArrayPosition(result, optEgdes[k][0][0]+i, optEdges[k][1][1]), getCordsFromArrayPosition(result, pos3, pos4))

					newItem := &Item{pos: [2]int{(optEdges[k][0][0] + i), optEdges[k][1][1]}, distance2: -alt, optimized: true, distance: -alt - altHelp}
					pre[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][1][1]] = []int{node[0], node[1]}
					distances[(optEdges[k][0][0]+i)*len(result[0])+optEdges[k][1][1]] = alt

					heap.Push(&pq, newItem)

				}

			}

			//set nodes within negative
			for yAxis := optEgdes[k][0][0] + 1; yAxis < optEgdes[k][1][0]; yAxis++ {
				for xAxis := optEgdes[k][0][1] + 1; xAxis < optEgdes[k][1][1]; xAxis++ {

					//	pre[yAxis*len(result[0])+xAxis] = []int{node[0], node[1]}
					//distances[yAxis*len(result[0])+xAxis] = -3
					visitedPoints[yAxis*len(result[0])+xAxis] = true
				}
			}

			//continue outer2

		}
		/*
			if k != -1 && node[0] >= optEgdes[k][0][0] && node[1] >= optEgdes[k][0][1] && node[0] <= optEgdes[k][1][0] && node[1] <= optEgdes[k][1][1] && (pos3 >= optEgdes[k][0][0] && pos4 >= optEgdes[k][0][1] && pos3 <= optEgdes[k][1][0] && pos4 <= optEgdes[k][1][1]) {

				var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, pos3, pos4))

				var alt = distances[node[0]*len(result[0])+node[1]] + distance
				//var alt2 = alt + GreatCircleDistance(getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]), getCordsFromArrayPosition(result, pos3, pos4))

				if alt <= distances[pos3*len(result[0])+pos4] {

					pre[pos3*len(result[0])+pos4] = []int{node[0], node[1]}
					distances[pos3*len(result[0])+pos4] = alt

				}
				println("optimized amount")
				for i := 0; i <= len(result)-1; i = i + 10 {
					for j := 0; j <= len(result[i])-1; j = j + 10 {

						if suchraum[i*len(result[0])+j] == true {
							print("X")
						} else {
							print(" ")
						}
					}
					println("")
				}
				println("optimized amount")
				println(optimizedAmount)
				return pre, distances, node, counterPOPS
			}
		*/
		var neighbours = getEdgesPosition(node[0], node[1])

	neig:
		for k := 0; k <= len(neighbours)-1; k++ {

			if visitedPoints[neighbours[k][0]*len(result[0])+neighbours[k][1]] == true {

				continue neig
			}

			var distance = GreatCircleDistance(getCordsFromArrayPosition(result, node[0], node[1]), getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]))

			var alt = distances[node[0]*len(result[0])+node[1]] + distance
			//var alt2 = alt + GreatCircleDistance(getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]), getCordsFromArrayPosition(result, pos3, pos4))

			if alt < distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] {
				newItem := &Item{pos: neighbours[k], distance2: -alt, distance: -alt - GreatCircleDistance(getCordsFromArrayPosition(result, neighbours[k][0], neighbours[k][1]), getCordsFromArrayPosition(result, pos3, pos4))}
				pre[neighbours[k][0]*len(result[0])+neighbours[k][1]] = []int{node[0], node[1]}
				distances[neighbours[k][0]*len(result[0])+neighbours[k][1]] = alt
				if visitedPoints[neighbours[k][0]*len(result[0])+neighbours[k][1]] == true {
					println("fehler11111111222222222222233333333")
					println(node[0])
					println(node[1])

				}
				heap.Push(&pq, newItem)

			}

		}

		if node[0] == pos3 && node[1] == pos4 {

			for i := 0; i <= len(result)-1; i = i + 10 {
				for j := 0; j <= len(result[i])-1; j = j + 10 {

					if suchraum[i*len(result[0])+j] == true {
						//print("X")
					} else {
						//	print(" ")
					}
				}
				//	println("")
			}

			return pre, distances, node, counterPOPS

		}

	}

	return pre, distances, node, counterPOPS

}
