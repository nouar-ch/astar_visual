package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type point struct {
	x, y int
}

type node struct {
	p      point
	parent *node
	dist   float64
}

func (n *node) neighbors(c, l int, mat [][]int) []*node {
	var neighbors []*node

	if n.p.x > 0 && mat[n.p.x-1][n.p.y] == 0 {
		neighbors = append(neighbors, &node{point{x: n.p.x - 1, y: n.p.y}, n, n.dist + 1})
	}
	if n.p.y > 0 && mat[n.p.x][n.p.y-1] == 0 {
		neighbors = append(neighbors, &node{point{x: n.p.x, y: n.p.y - 1}, n, n.dist + 1})
	}
	if n.p.x < c-1 && mat[n.p.x+1][n.p.y] == 0 {
		neighbors = append(neighbors, &node{point{x: n.p.x + 1, y: n.p.y}, n, n.dist + 1})
	}
	if n.p.y < l-1 && mat[n.p.x][n.p.y+1] == 0 {
		neighbors = append(neighbors, &node{point{x: n.p.x, y: n.p.y + 1}, n, n.dist + 1})
	}
	/*if n.p.x < c && n.p.y < l && mat[n.p.x+1][n.p.y+1] == 0 {
		neighbors = append(neighbors, &node{point{x: n.p.x + 1, y: n.p.y + 1}, n, n.dist + 1})
	}
	if n.p.x < c && n.p.y > 0 && mat[n.p.x+1][n.p.y-1] == 0 {
		neighbors = append(neighbors, &node{point{x: n.p.x + 1, y: n.p.y - 1}, n, n.dist + 1})
	}
	if n.p.x > 0 && n.p.y < l && mat[n.p.x-1][n.p.y+1] == 0 {
		neighbors = append(neighbors, &node{point{x: n.p.x - 1, y: n.p.y + 1}, n, n.dist + 1})
	}
	if n.p.x > 0 && n.p.y > 0 && mat[n.p.x-1][n.p.y-1] == 0 {
		neighbors = append(neighbors, &node{point{x: n.p.x - 1, y: n.p.y - 1}, n, n.dist + 1})
	}*/
	return neighbors
}

type astar struct {
	matrix               *[][]int
	l, c                 int
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

	start := node{algo.startPoint, nil, 0}
	nodes = append(nodes, &start)

	parentNode := &start

	for len(nodes) > 0 {
		time.Sleep(time.Millisecond)
		parentNode = nodes[0]
		(*algo.matrix)[parentNode.p.x][parentNode.p.y] = 1
		if len(nodes) > 1 {
			nodes = nodes[1:]
		} else {
			nodes = nodes[0:0]
		}

		if parentNode.p.x == algo.endPoint.x && parentNode.p.y == algo.endPoint.y {
			return parentNode
		}

		nodes = append(nodes, parentNode.neighbors(algo.c, algo.l, *algo.matrix)...)

		sort.Slice(nodes, func(i, j int) bool {
			a := nodes[i].dist + algo.h(nodes[i].p, algo.endPoint)
			b := nodes[j].dist + algo.h(nodes[j].p, algo.endPoint)

			return a < b
		})
		/*for i := range nodes {
			fmt.Println(*nodes[i])
		}
		break*/
	}

	return nil
}

func printMatrix(mat *[][]int) {
	m := len(*mat)
	for i := 0; i < m; i++ {
		fmt.Println((*mat)[i])
	}
}

func gui(matrix *[][]int, m int32, n int32) {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 800, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	rw := surface.W / int32(n)
	rh := surface.H / int32(m)

	running := true
	var i, j int32
	var rect sdl.Rect

	for running {
		surface.FillRect(nil, 0)
		for i = 0; i < m; i++ {
			for j = 0; j < n; j++ {
				rect = sdl.Rect{i*rh + 1, j*rw + 1, rw - 1, rh - 1}
				if (*matrix)[i][j] == 1 {
					surface.FillRect(&rect, 0xff00ff00)
				} else if (*matrix)[i][j] == 2 {
					surface.FillRect(&rect, 0x00000000)
				} else if (*matrix)[i][j] == 3 {
					surface.FillRect(&rect, 0xff0000ff)
				} else {
					surface.FillRect(&rect, 0xffffffff)
				}
			}
		}
		window.UpdateSurface()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}

func init() {
	runtime.LockOSThread()
}

func main() {
	mat := make([][]int, 25)
	for i := 0; i < 25; i++ {
		mat[i] = make([]int, 25)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 25; i++ {
		for j := 0; j < 25; j++ {
			if r.Float32() < 0.2 {
				mat[i][j] = 2
			} else {
				mat[i][j] = 0
			}
		}
	}
	/*{
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}*/

	printMatrix(&mat)

	var algo = astar{
		matrix:     &mat,
		l:          25,
		c:          25,
		startPoint: point{0, 0},
		endPoint:   point{24, 24},
		h: func(p1, p2 point) float64 {
			return math.Abs(float64(p1.x-p2.x)) + math.Abs(float64(p1.y-p2.y))
		}}

	var result *node = nil
	go func() {
		result = algo.run()
		for result != nil {
			fmt.Println(result.p)
			mat[result.p.x][result.p.y] = 3
			result = result.parent
		}
	}()
	gui(&mat, 25, 25)
}
