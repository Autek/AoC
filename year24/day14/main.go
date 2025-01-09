package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type vec struct {
	x, y int
}

type robot struct {
	pos, speed vec
}

func parse(fileName string) []robot {
	file, err := os.ReadFile(fileName) 
	if err != nil {
		panic(err)
	}
	
	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
	robots := make([]robot, len(lines))
	for i, line := range lines {
		var posX, posY, speedX, speedY int
		fmt.Sscanf(line, "p=%d,%d v=%d,%d", &posX, &posY, &speedX, &speedY)
		robots[i] = robot{vec{posX, posY}, vec{speedX, speedY}}
	}

	return robots
}

func posAfterN(r *robot, seconds, maxX, maxY int) {
	// we use (a % b + b) % b to compute the modulo and stay positive
	modX := maxX + 1
	modY := maxY + 1
	r.pos.x = ((r.speed.x * seconds + r.pos.x) % modX + modX) % modX
	r.pos.y = ((r.speed.y * seconds + r.pos.y) % modY + modY) % modY
}

func computeSafetyFactor(robots []robot, maxX, maxY int) int {
	topLeft, topRight, botLeft, botRight := 0, 0, 0, 0
	for _, r := range robots {
		if r.pos.x < (maxX + 1) / 2 && r.pos.y < (maxY + 1) / 2 {
			topLeft++
		} else if r.pos.x > maxX / 2 && r.pos.y < (maxY + 1) / 2 {
			topRight++
		} else if r.pos.x < (maxX + 1) / 2 && r.pos.y > maxY / 2 {
			botLeft++
		} else if r.pos.x > maxX / 2 && r.pos.y > maxY / 2 {
			botRight++
		}

	}
	return  topLeft * topRight * botLeft * botRight
}

func printRobots(robots []robot, maxX, maxY int) {
	grid := make([][]int, maxY + 1)
	for i := range grid {
		grid[i] = make([]int, maxX + 1)
	}

	for _, r := range robots {
		i, j := r.pos.y, r.pos.x
		grid[i][j]++
	}

	str := ""
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == 0 {
				str += "."
			} else {
				str += strconv.Itoa(grid[i][j])
			}
		}
		str += "\n"
	}
	fmt.Println(str)
	fmt.Println()
}

func sol1() {
	robots := parse("input.txt")
	seconds := 100
	maxX := 100
	maxY := 102
	for i := range robots {
		posAfterN(&robots[i], seconds, maxX, maxY)
	}
	fmt.Println(computeSafetyFactor(robots, maxX, maxY))
}

func sol2() {
	robots := parse("input.txt")
	maxX := 100
	maxY := 102
	firstPseudoPattern := 68
	pseudoPeriod := 101
	for i := range robots {
		posAfterN(&robots[i], firstPseudoPattern, maxX, maxY)
	}
	printRobots(robots, maxX, maxY)
	time.Sleep(time.Second)
	for i := 1; ; i++{
		for i := range robots {
			posAfterN(&robots[i], pseudoPeriod, maxX, maxY)
		}
		fmt.Println("seconds:", i * pseudoPeriod + firstPseudoPattern)
		printRobots(robots, maxX, maxY)
		time.Sleep(time.Second)
	}
}

func main() {
	sol1()
	//sol2()
}
