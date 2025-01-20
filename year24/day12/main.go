package main

import (
	"fmt"
	"os"
	"strings"
)

type gardenMap [][]string

type region struct {
	perimeter, area int
}

type position struct {
	i, j int
}

const (
	none = edge(iota)
	hortop
	horbot
	vertleft
	vertright
)

type edge int

func parse(fileName string) gardenMap {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
	garden := make(gardenMap, len(lines))
	for i, line := range lines {
		row := make([]string, len(line))
		for j, r := range line {
			row[j] = string(r)
		}
		garden[i] = row
	}
	return garden
}

func neighbours(i, j, maxI, maxJ int) []position {
	directions := []position{position{0, 1}, position{1, 0}, position{0, -1}, position{-1, 0}}
	positions := make([]position, 0, len(directions))
	for _, dir := range directions {
		newI, newJ := i+dir.i, j+dir.j
		if newI >= 0 && newI <= maxI && newJ >= 0 && newJ <= maxJ {
			positions = append(positions, position{newI, newJ})
		}
	}
	return positions
}

func getRegions(garden gardenMap) []region {
	reached := make([][]bool, len(garden))
	for i := range reached {
		reached[i] = make([]bool, len(garden[0]))
	}

	var recRegion func(int, int) (perimeter, area int)
	recRegion = func(i, j int) (perimeter, area int) {
		reached[i][j] = true
		nei := neighbours(i, j, len(garden)-1, len(garden[0])-1)
		perimeter = 4 - len(nei)
		area = 1
		for _, n := range nei {
			if garden[i][j] != garden[n.i][n.j] {
				perimeter++
			} else if !reached[n.i][n.j] {
				p, a := recRegion(n.i, n.j)
				perimeter += p
				area += a
			}
		}
		return perimeter, area
	}

	regions := make([]region, 0, 1<<10)
	for i, line := range garden {
		for j := range line {
			if !reached[i][j] {
				p, a := recRegion(i, j)
				regions = append(regions, region{p, a})
			}
		}
	}
	return regions
}

func sol1() {
	garden := parse("input.txt")
	regions := getRegions(garden)
	sum := 0
	for _, r := range regions {
		sum += r.area * r.perimeter
	}
	fmt.Println(sum)
}

type region2 [][]bool

func countEdges(array [][]edge) int {
	sum := 0
	for i := range array {
		currentEdge := none
		for j := range array[i] {
			if array[i][j] != currentEdge {
				if array[i][j] == horbot || array[i][j] == hortop {
					sum++
				}
				currentEdge = array[i][j]
			}
		}
	}

	for i := range array[0] {
		currentEdge := none
		for j := range array {
			if array[j][i] != currentEdge {
				if array[j][i] == vertleft || array[j][i] == vertright {
					sum++
				}
				currentEdge = array[j][i]
			}
		}
	}
	return sum
}

func processEdges(array region2) ([][]edge, [][]edge) {
	vertedges := make2DArray[edge](len(array), len(array[0])+1)
	horedges := make2DArray[edge](len(array)+1, len(array[0]))

	vertKernel := func(i1, i2 bool) edge {
		if i1 != i2 {
			if i2 == true {
				return vertleft
			}
			return vertright
		}
		return none
	}

	for i := range array {
		for j := range array[i][:len(array[0])-1] {
			vertedges[i][j+1] = vertKernel(array[i][j], array[i][j+1])
		}
		if array[i][0] == true {
			vertedges[i][0] = vertleft
		}
		if array[i][len(array[i])-1] == true {
			vertedges[i][len(array[i])] = vertright
		}
	}

	horKernel := func(i1, i2 bool) edge {
		if i1 != i2 {
			if i2 == true {
				return hortop
			}
			return horbot
		}
		return none
	}
	for j := range array[0] {
		for i := range array[:len(array)-1] {
			horedges[i+1][j] = horKernel(array[i][j], array[i+1][j])
		}
		if array[0][j] == true {
			horedges[0][j] = hortop
		}
		if array[len(array)-1][j] == true {
			horedges[len(array)][j] = horbot
		}
	}
	return horedges, vertedges
}

func getArea(r region2) int {
	sum := 0
	for i := range r {
		for j := range r[i] {
			if r[i][j] {
				sum++
			}
		}
	}
	return sum
}

func getRegions2(garden gardenMap) []region2 {
	reached := make2DArray[bool](len(garden), len(garden[0]))

	var recRegion func(int, int, region2)
	recRegion = func(i, j int, r region2) {
		reached[i][j] = true
		r[i][j] = true
		nei := neighbours(i, j, len(garden)-1, len(garden[0])-1)
		for _, n := range nei {
			if garden[i][j] == garden[n.i][n.j] && !reached[n.i][n.j] {
				recRegion(n.i, n.j, r)
			}
		}
	}

	regions := make([]region2, 0, 1<<10)
	for i, line := range garden {
		for j := range line {
			if !reached[i][j] {
				r := make2DArray[bool](len(garden), len(garden[0]))
				recRegion(i, j, r)
				regions = append(regions, r)
			}
		}
	}
	return regions
}

func make2DArray[T any](height, width int) [][]T {
	lines := make([][]T, height)
	for i := range lines {
		lines[i] = make([]T, width)
	}
	return lines
}

func sol2() {
	garden := parse("input.txt")
	regions := getRegions2(garden)
	sum := 0
	for _, r := range regions {
		horizontalEdges, verticalEdges := processEdges(r)
		edges := countEdges(horizontalEdges) + countEdges(verticalEdges)
		area := getArea(r)
		sum += area * edges
	}
	fmt.Println(sum)
}

func main() {

	sol1()
	sol2()
}
