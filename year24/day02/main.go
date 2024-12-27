package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func parseInput(fileName string) [][]int {

	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
	reports := make([][]int, len(lines))
	for i, line := range lines {
		splitLine := strings.Split(strings.TrimSpace(line), " ")
		levels := make([]int, len(splitLine))
		for j, level_str := range splitLine {
			levels[j], err = strconv.Atoi(level_str)
			if err != nil {
				panic(err)
			}
		}
		reports[i] = levels
	}
	return reports
}

func isSafe(report []int) bool{
	if len(report) < 2 {
		panic("report is too small")
	}

	lastDelta := 0
	for i := range report[:len(report) -1] {
		delta := report[i] - report[i+1]
		abs_delta := int(math.Abs(float64(delta)))
		if abs_delta < 1 || abs_delta > 3 {
			return false
		}
		if lastDelta * delta < 0 {
			return false
		}
		lastDelta = delta
	}
	return true
}

func sol1(){
	reports := parseInput("input.txt")
	nbSafeReports := 0
	for _, report := range reports{
		if (isSafe(report)) {
			nbSafeReports++
		}
	}
	fmt.Println(nbSafeReports)
}

// this can probably be better implemented but as a quick and dirty way to do it
// let's just remove all the samples one by one (and also without removing)
func isSafeDampener(report []int) bool{
	if len(report) < 2 {
		panic("report is too small")
	}
	if isSafe(report) {
		return true
	}
	tempReport := make([]int, 0, len(report) -1)
	for i := range report {
		tempReport = tempReport[:0]
		tempReport = append(tempReport, report[:i]...)
		tempReport = append(tempReport, report[i+1:]...)
		if isSafe(tempReport) {
			return true
		}
	}
	return false
}

func sol2(){
	reports := parseInput("input.txt")
	nbSafeReports := 0
	for _, report := range reports{
		if (isSafeDampener(report)) {
			nbSafeReports++
		}
	}
	fmt.Println(nbSafeReports)
}

func main() {
	sol2()
}
