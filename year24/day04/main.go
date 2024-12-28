package main

import (
	"fmt"
	"os"
	"strings"
)

type grid = [][]string

func parseInput(fileName string) grid{
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
	g := make(grid, len(lines))
	for i, line := range lines {
		for _, r := range strings.TrimSpace(line) {
			g[i] = append(g[i], string(r))
		}
	}
	return g
}

func checkXmas(i int, j int, g grid) int {
	letters := "XMAS"
	directions := [][]int{{0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1}, {-1, 0}, {-1, 1}}
	xmas := 0
	outterloop:
	for _, direction := range directions {
		for index, letter := range letters {
			newI := i + index * direction[0]
			newJ := j + index * direction[1]
			if newI < 0 || newI >= len(g) || newJ < 0 || newJ >= len(g[0]) {
				continue outterloop
			}
			if g[newI][newJ] != string(letter) {
				continue outterloop
			}
		}
		xmas += 1
	}
	return xmas
}

func sol1(){
	g := parseInput("input.txt")
	sum := 0
	for i, line := range g{
		for j, char := range line {
			if char == "X" {
				sum += checkXmas(i, j, g)
			}
		}
	}
	fmt.Println(sum)
}

func checkXmas2(i int, j int, g grid) bool {
	if g[i][j] != "A" {
		return false
	}
	if i < 1 || i >= len(g) - 1 || j < 1 || j >= len(g[0]) - 1 {
		return false
	}
	topRight := g[i-1][j+1]
	topLeft := g[i-1][j-1]
	botRight := g[i+1][j+1]
	botLeft := g[i+1][j-1]

	if ((topRight == "M" && botLeft == "S") || 
	   (topRight == "S" && botLeft == "M")) &&
	   ((topLeft == "M" && botRight == "S") || 
	   (topLeft == "S" && botRight == "M")) {
		   return true
	   }
	return false
}

func sol2() {
	g := parseInput("input.txt")
	sum := 0
	for i, line := range g{
		for j := range line {
			if checkXmas2(i, j, g){
				sum++
			}
		}
	}
	fmt.Println(sum)
}

func main(){
	sol1()
	sol2()
}
