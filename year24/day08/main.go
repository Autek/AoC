package main

import (
	"fmt"
	"os"
	"strings"
)

type position struct {
	x, y int
}

type antenna struct {
	name string
	pos []position
}

func getAntena(antennaName string, antennas []antenna) *antenna {
	for i, a := range antennas {
		if a.name == antennaName {
			return &antennas[i]
		}
	}
	return nil
}

func parse(fileName string) ([]antenna, position) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
	antennas := []antenna{}
	for i, line := range lines {
		chars := strings.Split(strings.TrimSpace(line), "")
		for j, char := range chars {
			if char == "." {
				continue 
			}
			a := getAntena(char, antennas)
			pos := position{x: j, y: i}
			if a == nil {
				newAntenna := antenna{name: char, pos: []position{pos}}
				antennas = append(antennas, newAntenna)
			} else {
				a.pos = append(a.pos, pos)
			}
		}
	}
	maxPos := position{x: len(lines[0]) - 1, y: len(lines) - 1}
	return antennas, maxPos
}

func getAntiNodes(a antenna) map[position]struct{}{
	positions := map[position]struct{}{}
	for i, p1 := range a.pos[:len(a.pos)-1] {
		for _, p2 := range a.pos[i+1:] {
			diff := position{x: p1.x - p2.x, y: p1.y - p2.y}
			a1 := position{x: p1.x + diff.x, y: p1.y + diff.y}
			a2 := position{x: p2.x - diff.x, y: p2.y - diff.y}
			positions[a1] = struct{}{}
			positions[a2] = struct{}{}
		}
		
	}
	return positions
}

func inRange(p position, maxPos position) bool {
	return p.x >= 0 && p.x <= maxPos.x && p.y >= 0 && p.y <= maxPos.y
}

func pruneOutsideRange(positions map[position]struct{}, maxPos position) {
	for p := range positions {
		if !inRange(p, maxPos) {
			delete(positions, p)
		}
	}
}

func union[T comparable](m1, m2 map[T]struct{}) {
	for k := range m2 {
		m1[k] = struct{}{}
	}
}

func sol1() {
	antennas, maxPos := parse("input.txt")
	antiNodes := map[position]struct{}{}
	for _, a := range antennas {
		union(antiNodes, getAntiNodes(a))
	}
	pruneOutsideRange(antiNodes, maxPos)
	fmt.Println(len(antiNodes))
}

func getAntiNodes2(a antenna, maxPos position) map[position]struct{}{
	positions := map[position]struct{}{}
	for i, p1 := range a.pos[:len(a.pos)-1] {
		for _, p2 := range a.pos[i+1:] {
			diff := position{x: p1.x - p2.x, y: p1.y - p2.y}
			minAnode := p1

			// find the first harmonic antinode outside range
			for inRange(minAnode, maxPos) {
				minAnode.x += diff.x
				minAnode.y += diff.y
			}

			// add all the harmonic antinodes in range to the set
			aNode := position{minAnode.x - diff.x, minAnode.y - diff.y}
			for inRange(aNode, maxPos) {
				positions[aNode] = struct{}{}
				aNode.x -= diff.x
				aNode.y -= diff.y
			}
		}
	}
	return positions
}

func sol2() {
	antennas, maxPos := parse("input.txt")
	antiNodes := map[position]struct{}{}
	for _, a := range antennas {
		union(antiNodes, getAntiNodes2(a, maxPos))
	}
	fmt.Println(len(antiNodes))
}

func main() {
	sol1()
	sol2()
}
