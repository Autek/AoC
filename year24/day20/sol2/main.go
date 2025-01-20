package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type cell int

const (
	wall = cell(iota)
	path
)

type position struct {
	x, y int
}

type pathID struct {
	start, end position
}

func parse(fileName string) ([][]cell, position, position) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	file_slice := strings.Split(strings.TrimSpace(string(file)), "\n")
	grid := make([][]cell, len(file_slice))
	var start, end position
	for i, line := range file_slice {
		grid[i] = make([]cell, len(line))
		for j, r := range line {
			switch r {
			case '#':
				grid[i][j] = wall
			case '.':
				grid[i][j] = path
			case 'E':
				grid[i][j] = path
				end = position{j, i}
			case 'S':
				grid[i][j] = path
				start = position{j, i}
			}
		}
	}
	return grid, start, end
}

func getNeighbours(grid [][]cell, p position) []position {
	directions := []position{
		{1, 0},
		{-1, 0},
		{0, 1},
		{0, -1},
	}
	neighbours := make([]position, 0, 4)
	for _, dir := range directions {
		x := p.x + dir.x
		y := p.y + dir.y
		if y >= 0 && y < len(grid) && x >= 0 && x < len(grid[0]) {
			neighbours = append(neighbours, position{x, y})
		}
	}
	return neighbours
}

func findPathCost(grid [][]cell, start, end position) map[position]int {
	visited := map[position]struct{}{start: {}}
	previous := map[position]position{}
	currentLevel := []position{start}
	nextLevel := []position{}

	for len(currentLevel) != 0 {
		for _, pos := range currentLevel {
			for _, neighbour := range getNeighbours(grid, pos) {
				_, ok := visited[neighbour]
				if !ok && grid[neighbour.y][neighbour.x] == path {
					nextLevel = append(nextLevel, neighbour)
					visited[neighbour] = struct{}{}
					previous[neighbour] = pos
				}
			}
		}
		currentLevel, nextLevel = nextLevel, currentLevel[:0]
	}

	posToDistance := map[position]int{}
	distFromEnd := 0
	pos := end
	for pos != start {
		posToDistance[pos] = distFromEnd
		pos = previous[pos]
		distFromEnd++
	}
	posToDistance[start] = distFromEnd
	return posToDistance
}

func BFSfree(grid [][]cell, pathCost map[position]int, shortcutLength int,
	start position) ([]int, []position) {

	visited := map[position]struct{}{start: {}}
	currentLevel := []position{start}
	nextLevel := []position{}

	initialCost := pathCost[start]
	length := 1
	shortcuts := []int{}
	shortcutEnds := []position{}
	for len(currentLevel) != 0 && length <= shortcutLength {
		for _, pos := range currentLevel {
			for _, neighbour := range getNeighbours(grid, pos) {
				if _, ok := visited[neighbour]; !ok {
					nextLevel = append(nextLevel, neighbour)
					visited[neighbour] = struct{}{}
					if grid[neighbour.y][neighbour.x] == path {
						shortcut := initialCost - length - pathCost[neighbour]
						shortcuts = append(shortcuts, shortcut)
						shortcutEnds = append(shortcutEnds, pos)
					}
				}
			}
		}
		length++
		currentLevel, nextLevel = nextLevel, currentLevel[:0]
	}
	return shortcuts, shortcutEnds
}

func findShortCuts(grid [][]cell, pathCost map[position]int,
	threshold, shortcutLength int) []pathID {

	paths := []pathID{}
	for position := range pathCost {
		shortcuts, ends := BFSfree(grid, pathCost, shortcutLength, position)
		for i, shortcut := range shortcuts {
			if shortcut >= threshold {
				paths = append(paths, pathID{position, ends[i]})
			}
		}
	}
	return paths
}

func sol() {
	threshold := 100
	shortcutLength := 20
	grid, start, end := parse("../input.txt")
	posMap := findPathCost(grid, start, end)
	bestShorts := findShortCuts(grid, posMap, threshold, shortcutLength)
	fmt.Println(len(bestShorts))
}

func main() {
	sol()
}

//------------------DEBUGGING------------------

func printGrid(grid [][]cell, posMap map[position]int, shortcut []position) {
	str := ""
	for i, row := range grid {
		for j, c := range row {
			if v, ok := posMap[position{j, i}]; ok {
				str += strconv.Itoa(v % 10)
			} else if slices.Contains(shortcut, position{j, i}) {
				str += "^"
			} else if c == wall {
				str += "#"
			} else if c == path {
				str += "."
			}
		}
		str += "\n"
	}
	fmt.Print(str)
}
