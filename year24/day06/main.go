package main

import (
	"fmt"
	"os"
	"strings"
)

type cell struct {
	isReached, isObstacle bool
}

func parse(fileName string) [][]cell {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
	grid := make([][]cell, len(lines))
	for i, line := range lines {
		chars := strings.Split(strings.TrimSpace(line), "")
		row := make([]cell, len(chars))
		for j, char := range chars {
			if char == "." {
				row[j] = cell{false, false}
			} else if char == "#" {
				row[j] = cell{false, true}
			} else if char == "^" {
				row[j] = cell{true, false}
			}
		}
		grid[i] = row
	}
	return grid
}

func findPlayer(grid [][]cell) (int, int) {
	for i, row := range grid {
		for j, c := range row {
			if c.isReached {
				return i, j
			}
		}
	}
	return -1, -1
}

func isInBounds(grid [][]cell, i int, j int) bool {
	return i >= 0 && i < len(grid) && j >= 0 && j < len(grid[0])
}

func runSimulation(grid [][]cell) [][]cell {
	directions := []direction{direction{-1, 0},
		direction{0, 1},
		direction{1, 0},
		direction{0, -1}}
	i, j := findPlayer(grid)
	lastI, lastJ := -1, -1
	directionIndex := 0
	for true {
		lastI = i
		lastJ = j
		i += directions[directionIndex%len(directions)].i
		j += directions[directionIndex%len(directions)].j
		for isInBounds(grid, i, j) && grid[i][j].isObstacle {
			directionIndex++
			i = lastI + directions[directionIndex%len(directions)].i
			j = lastJ + directions[directionIndex%len(directions)].j
		}
		if !isInBounds(grid, i, j) {
			break
		}
		grid[i][j].isReached = true
	}
	return grid
}

func countReached(grid [][]cell) int {
	sum := 0
	for _, row := range grid {
		for _, c := range row {
			if c.isReached {
				sum++
			}
		}
	}
	return sum
}

func printGrid(grid [][]cell) {
	str := ""
	for _, row := range grid {
		for _, c := range row {
			if c.isReached {
				str += "X"
			} else if c.isObstacle {
				str += "#"
			} else {
				str += "."
			}

		}
		str += "\n"
	}
	fmt.Print(str)
}

// the first way I can see to do this is to build a 2D slice of custom type cell
// then we run the simulation on the array and we just count the cells reached
// this should run in O(number of cells) which seems reasonable
func sol1() {
	grid := parse("input.txt")
	grid = runSimulation(grid)
	reached := countReached(grid)
	fmt.Println(reached)
}

type direction struct {
	i, j int
}

type cell2 struct {
	isReached, isObstacle bool
	dir                   direction
}

func parse2(fileName string) [][]cell2 {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
	grid := make([][]cell2, len(lines))
	for i, line := range lines {
		chars := strings.Split(strings.TrimSpace(line), "")
		row := make([]cell2, len(chars))
		for j, char := range chars {
			if char == "." {
				row[j] = cell2{false, false, direction{}}
			} else if char == "#" {
				row[j] = cell2{false, true, direction{}}
			} else if char == "^" {
				row[j] = cell2{true, false, direction{-1, 0}}
			}
		}
		grid[i] = row
	}
	return grid
}

func findPlayer2(grid [][]cell2) (int, int) {
	for i, row := range grid {
		for j, c := range row {
			if c.isReached {
				return i, j
			}
		}
	}
	return -1, -1
}

func isInBounds2(grid [][]cell2, i int, j int) bool {
	return i >= 0 && i < len(grid) && j >= 0 && j < len(grid[0])
}

func checkLoop(grid [][]cell2) bool {
	directions := []direction{direction{-1, 0},
		direction{0, 1},
		direction{1, 0},
		direction{0, -1}}
	i, j := findPlayer2(grid)
	lastI, lastJ := -1, -1
	directionIndex := 0
	for isInBounds2(grid, i, j) {
		lastI = i
		lastJ = j
		i += directions[directionIndex%len(directions)].i
		j += directions[directionIndex%len(directions)].j
		for isInBounds2(grid, i, j) && grid[i][j].isObstacle {
			directionIndex++
			i = lastI + directions[directionIndex%len(directions)].i
			j = lastJ + directions[directionIndex%len(directions)].j
		}
		if isInBounds2(grid, i, j) {
			if grid[i][j].isReached && grid[i][j].dir == directions[directionIndex%len(directions)] {
				return true
			}
			grid[i][j].isReached = true
			grid[i][j].dir = directions[directionIndex%len(directions)]
		}
	}
	return false
}

func printGrid2(grid [][]cell2) {
	str := ""
	for _, row := range grid {
		for _, c := range row {
			if c.isReached {
				str += "X"
			} else if c.isObstacle {
				str += "#"
			} else {
				str += "."
			}

		}
		str += "\n"
	}
	fmt.Print(str)
}

func copyGrid(grid [][]cell2) [][]cell2 {
	newGrid := make([][]cell2, len(grid))
	for i, row := range grid {
		newGrid[i] = make([]cell2, len(row))
		copy(newGrid[i], row)
	}
	return newGrid
}

// to do this part for now I can think of doing it in a bruteForce manner
// let's check if there is a loop after adding an obstacle to every Point
// in the grid (this can be optimised to put an obstacle on the points reached
// but the complexity would stay the same so let's do it the dirty way)
func sol2() {
	grid := parse2("input.txt")
	sum := 0
	for i, row := range grid {
		for j, _ := range row {
			testGrid := copyGrid(grid)
			//we can't place an obstacle on the person
			if testGrid[i][j].isReached {
				continue
			}
			testGrid[i][j].isObstacle = true
			if checkLoop(testGrid) {
				sum++
			}
		}
	}
	fmt.Println(sum)
}

func main() {
	sol1()
	sol2()
}
