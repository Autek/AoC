package main

import (
	"fmt"
	"os"
	"strings"
)

type cell int
const(
	free = cell(iota)
	corrupted
)

type position struct {
	x, y int
}

func parse(fileName string) []position {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
	bytes := make([]position, len(lines))

	for i, line := range lines {
		var x, y int
		fmt.Sscanf(line, "%d,%d", &x, &y)
		bytes[i] = position{x, y}
	}
	return bytes
}

func getNeighbours(p position, maze [][]cell) []position {
	directions := []position{
		position{0, 1}, 
		position{1, 0},
		position{0, -1}, 
		position{-1, 0},
	}
	
	neigbours := make([]position, 0, len(directions))
	for _, dir := range directions {
		x := p.x + dir.x
		y := p.y + dir.y
		if x >= 0 && x < len(maze[0]) && y >= 0 && y < len(maze) {
			if maze[y][x] == free {
				neigbours = append(neigbours, position{x, y})
			}
		}
	}
	return neigbours
}

func BFSSolve(maze [][]cell, start, finish position) int {
	reached := make([][]bool, len(maze))
	for i := range reached {
		reached[i] = make([]bool, len(maze[0]))
	}
	currentLevel := []position{start}
	reached[start.y][start.x] = true
	nextLevel := []position{}
	distance := 0
	for len(currentLevel) != 0 {
		for _, p := range currentLevel {
			if p == finish {
				return distance
			}

			for _, n := range getNeighbours(p, maze) {
				if !reached[n.y][n.x] {
					reached[n.y][n.x] = true
					nextLevel = append(nextLevel, n)	
				}
			}
		}
		distance++
		currentLevel, nextLevel = nextLevel, currentLevel[:0]
	}
	return -1
}

func setMaze(maze [][]cell, bytes []position, n int) {
	for i, p := range bytes {
		if i < n {
			maze[p.y][p.x] = corrupted
		} else {
			maze[p.y][p.x] = free
		}
	}
}

func binSearch(bytes []position, size int, start, finish position) position{
	maze := make([][]cell, size)
	for i := range maze {
		maze[i] = make([]cell, size)
	}

	left := 0
	right := len(bytes) - 1
	var mid int
	for left + 1 < right {
		mid = (left + right) / 2
		setMaze(maze, bytes, mid)
		if BFSSolve(maze, start, finish) == -1 {
			right = mid
		} else {
			left = mid
		}
	}
	return bytes[right - 1]
}

func sol() {
	size := 71
	start := position{0, 0}
	finish := position{size-1, size-1}
	bytes := parse("../input.txt")

	ans := binSearch(bytes, size, start, finish)
	fmt.Printf("%d,%d\n", ans.x, ans.y)
}

func main() {
	sol()
}

//-----------------DEBUG-----------------

func printMaze(maze [][]cell) {
	str := ""
	for _, row := range maze {
		for _, c := range row {
			if c == free {
				str += "."
			} else {
				str += "#"
			}
		}
		str += "\n"
	}
	fmt.Print(str)
}
