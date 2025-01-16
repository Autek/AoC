package main

import (
	"fmt"
	"os"
	"strings"
)

type cell int
const(
	wall = cell(iota)
	path
	reachedHor
	reachedVert
)

type axis int
const(
	horizontal = axis(iota)
	vertical
)

const (
	stepcost = 1
	turnCost = 1000
)

type edge struct {
	cost int
	from, to pos
	facing axis
}

type pos struct{
	x, y int
}

type direction struct {
	x, y int
	a axis
}

type posAndAxis struct {
	p pos
	a axis
}

//--------------------MinQueue implementation--------------------

type item[T any] struct {
	val T
	prio int
}

type minQueue[T any] []item[T] 

func (q * minQueue[T])push(e T, prio int) {
	l := append(*q, item[T]{e, prio})
	index := len(l) - 1
	parent := (index - 1) >> 1
	for parent >= 0 && prio < l[parent].prio {
		l[index], l[parent] = l[parent], l[index]
		index = parent
		parent = (parent - 1) >> 1
	}
	*q = l
}

func (q * minQueue[T])pop() T {
	if q.isEmpty() {
		panic("pop on empty queue")
	}

	l := *q
	popped := l[0]
	size := len(l) - 1
	l[0] = l[size]
	l = l[:size]

	index := 0
	child := 1
	for child < size {
		if child+1 < size && l[child].prio > l[child+1].prio {
			child++
		}
		if l[index].prio <= l[child].prio {
			break
		}

		l[index], l[child] = l[child], l[index]
		index = child
		child = (index << 1) + 1
	}
	*q = l

	return popped.val
}

func (q minQueue[T])isEmpty() bool {
	return len(q) == 0
}

//--------------------Solution functions--------------------

func parse(fileName string) [][]cell {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	file_slice := strings.Split(strings.TrimSpace(string(file)), "\n")
	grid := make([][]cell, len(file_slice))
	for i, line := range file_slice {
		grid[i] = make([]cell, len(line))
		for j, r := range line {
			if r == '#' {
				grid[i][j] = wall
			} else {
				grid[i][j] = path
			}
		}
	}
	return grid
}

func getNeighbours(grid [][]cell, p pos) []posAndAxis{
	neighbours := make([]posAndAxis, 0, 4)
	directions := []direction{
		{1, 0, horizontal},
		{-1, 0, horizontal},
		{0, 1, vertical},
		{0, -1, vertical},
	}
	for _, dir := range directions {
		nextPos := pos{p.x + dir.x, p.y + dir.y}

		if grid[nextPos.y][nextPos.x] == path {
			neighbours = append(neighbours, posAndAxis{nextPos, dir.a})
		}
	}
	return neighbours
}

func dijkstra(grid [][]cell, start, end pos) int{
	distances := make([][][]int, len(grid))
	for i := range distances {
		distances[i] = make([][]int, len(grid[0]))
		for j := range distances[i] {
			distances[i][j] = make([]int, 2)
			distances[i][j][0] = -1
			distances[i][j][1] = -1
		}
	}
	previous := make(map[posAndAxis]posAndAxis, len(grid) * len(grid[0]) * 2)
	q := make(minQueue[posAndAxis], 0, len(grid) * len(grid[0]))
	s := posAndAxis{start, horizontal}
	var r posAndAxis
	q.push(s, 0)
	distances[s.p.x][s.p.y][s.a] = 0
	for !q.isEmpty() {
		u := q.pop()
		for _, v := range getNeighbours(grid, u.p) {
			step := stepcost
			if v.a != u.a {
				step += turnCost
			}
			alt := distances[u.p.x][u.p.y][u.a] + step
			oldDist := distances[v.p.x][v.p.y][v.a]
			if oldDist == -1 || alt < oldDist {
				previous[v] = u
				distances[v.p.x][v.p.y][v.a] = alt
				q.push(v, alt)
				rVal := distances[r.p.x][r.p.y][r.a]
				if v.p == end && (rVal == -1 || alt < rVal) {
					r = v
				}
			}
		}
	}
	return distances[r.p.x][r.p.y][r.a]
}


func sol1() {
	grid := parse("input.txt")
	start := pos{1, len(grid) - 2}
	end := pos{len(grid[0]) - 2, 1}
	fmt.Println(dijkstra(grid, start, end))
}

func dijkstra2(grid [][]cell, start, end pos) (map[posAndAxis][]posAndAxis, posAndAxis, posAndAxis){
	distances := make([][][]int, len(grid))
	for i := range distances {
		distances[i] = make([][]int, len(grid[0]))
		for j := range distances[i] {
			distances[i][j] = make([]int, 2)
			distances[i][j][0] = -1
			distances[i][j][1] = -1
		}
	}
	previous := make(map[posAndAxis][]posAndAxis, len(grid) * len(grid[0]) * 2)
	q := make(minQueue[posAndAxis], 0, len(grid) * len(grid[0]))
	s := posAndAxis{start, horizontal}
	var r posAndAxis
	q.push(s, 0)
	distances[s.p.x][s.p.y][s.a] = 0
	for !q.isEmpty() {
		u := q.pop()
		for _, v := range getNeighbours(grid, u.p) {
			step := stepcost
			if v.a != u.a {
				step += turnCost
			}
			alt := distances[u.p.x][u.p.y][u.a] + step
			oldDist := distances[v.p.x][v.p.y][v.a]
			if oldDist == -1 || alt <= oldDist {
				if oldDist == -1 || alt == oldDist {
					previous[v] = append(previous[v], u)
				} else {
					previous[v] = previous[v][:1]
					previous[v][0] = u
				}
				distances[v.p.x][v.p.y][v.a] = alt
				q.push(v, alt)
				rVal := distances[r.p.x][r.p.y][r.a]
				if v.p == end && (rVal == -1 || alt < rVal) {
					r = v
				}
			}
		}
	}
	return previous, r, s
}

func countBestPathTiles(previous map[posAndAxis][]posAndAxis, start, end posAndAxis) int {
	currentCells := map[posAndAxis]struct{}{end: struct{}{}}
	posSet := map[pos]struct{}{}
	for len(currentCells) != 0 {
		for c := range currentCells {
			delete(currentCells, c)
			for _, n := range previous[c] {
				currentCells[n] = struct{}{}
			}
			posSet[c.p] = struct{}{}
		}
	}
	return len(posSet)
}

func sol2() {
	grid := parse("input.txt")
	start := pos{1, len(grid) - 2}
	end := pos{len(grid[0]) - 2, 1}
	previous, e, s := dijkstra2(grid, start, end)
	fmt.Println(countBestPathTiles(previous, s, e))
}

func main() {
	sol1()
	sol2()
}

//--------------------PrintPath debugging functions--------------------

func PrintPath2(grid [][]cell, prev map[posAndAxis][]posAndAxis, to, from posAndAxis) {
	currentCells := map[posAndAxis]struct{}{to: struct{}{}}
	for len(currentCells) != 0 {
		for c := range currentCells {
			for _, n := range prev[c] {
				currentCells[n] = struct{}{}
			}
			if c.a == horizontal {
				grid[c.p.y][c.p.x] = reachedHor
			}else {
				grid[c.p.y][c.p.x] = reachedVert
			}
			delete(currentCells, c)
		}
	}
	s := ""

	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == reachedHor {
				s += "-"
			} else if grid[i][j] == reachedVert {
				s += "|"
			} else if grid[i][j] == wall {
				s += "#"
			} else if grid[i][j] == path {
				s += "."
			}
		}
		s += "\n"
	}

	fmt.Print(s)
}

func PrintPath(grid [][]cell, prev map[posAndAxis]posAndAxis, to, from posAndAxis) {
	start := to
	for start != from {
		if start.a == horizontal {
			grid[start.p.y][start.p.x] = reachedHor
		}else {
			grid[start.p.y][start.p.x] = reachedVert
		}
		start = prev[start]
	}
	s := ""

	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == reachedHor {
				s += "-"
			} else if grid[i][j] == reachedVert {
				s += "|"
			} else if grid[i][j] == wall {
				s += "#"
			} else if grid[i][j] == path {
				s += "."
			}
		}
		s += "\n"
	}

	fmt.Print(s)
}
