package main

import (
	"fmt"
	"math"
	"sort"
)

type point struct {
	x, y int
}

type node struct {
	p      point
	dist   float64
	parent *node
}

type astar struct {
	matrix               [][]int
	startPoint, endPoint point
	h                    func(point, point) float64
}

func (algo *astar) run() *node {
	if algo.h == nil {
		algo.h = func(p1 point, p2 point) float64 {
			return math.Sqrt(math.Pow(float64(p1.x-p2.x), 2.0) + math.Pow(float64(p1.y-p2.y), 2.0))
		}
	}

	var nodes []*node

	start := node{algo.startPoint, 0, nil}
	nodes = append(nodes, &start)

	var currentNode *node
	parentNode := &start

	for len(nodes) > 0 {
		parentNode = nodes[0]
		if len(nodes) > 1 {
			nodes = nodes[1:]
		} else {
			nodes = nodes[0:0]
		}

		if parentNode.p.x == algo.endPoint.x && parentNode.p.y == algo.endPoint.y {
			return parentNode
		}

		if algo.matrix[parentNode.p.x][parentNode.p.y+1] == 0 {
			currentNode = &node{
				p:      point{parentNode.p.x, parentNode.p.y + 1},
				parent: parentNode}
			currentNode.dist = parentNode.dist + 1

			nodes = append(nodes, currentNode)
		}

		if algo.matrix[parentNode.p.x][parentNode.p.y-1] == 0 {
			currentNode = &node{
				p:      point{parentNode.p.x, parentNode.p.y - 1},
				parent: parentNode}
			currentNode.dist = parentNode.dist + 1
			nodes = append(nodes, currentNode)
		}

		if algo.matrix[parentNode.p.x+1][parentNode.p.y] == 0 {
			currentNode = &node{
				p:      point{parentNode.p.x + 1, parentNode.p.y},
				parent: parentNode}
			currentNode.dist = parentNode.dist + 1
			nodes = append(nodes, currentNode)
		}

		if algo.matrix[parentNode.p.x-1][parentNode.p.y] == 0 {
			currentNode = &node{
				p:      point{parentNode.p.x - 1, parentNode.p.y},
				parent: parentNode}
			currentNode.dist = parentNode.dist + 1
			nodes = append(nodes, currentNode)
		}

		sort.Slice(nodes, func(i, j int) bool {
			a := nodes[i].dist + algo.h(nodes[i].p, algo.endPoint)
			b := nodes[j].dist + algo.h(nodes[j].p, algo.endPoint)

			return a < b
		})
	}

	return nil
}

func main() {
	var mat = [][]int{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0}}

	var algo = astar{
		matrix:     mat,
		startPoint: point{0, 0},
		endPoint:   point{4, 4}}

	var dist = algo.run()

	fmt.Println(dist)
}
