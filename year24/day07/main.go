package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type eq struct {
	result int
	numbers []int
}

func parse(fileName string) []eq {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	fileStr := string(file)
	lines := strings.Split(string(strings.TrimSpace(fileStr)), "\n")
	result := make([]eq, len(lines))
	for i, line := range lines {
		splitLine := strings.Split(line, ": ")
		res, err := strconv.Atoi(splitLine[0])
		if err != nil {
			panic(err)
		}
		numbersStr := strings.Split(splitLine[1], " ")
		numbers := make([]int, len(numbersStr))
		for j, numbersStr := range numbersStr {
			nb, err := strconv.Atoi(numbersStr)
			if err != nil {
				panic(err)
			}
			numbers[j] = nb
		}
		result[i] = eq{result: res, numbers: numbers}
	}
	return result
}

func isValid(e eq) bool {
	var rec func([]int, int) bool
	rec = func(n []int, acc int) bool {
		if acc == e.result && len(n) == 0 {
			return true
		}
		if len(n) == 0 || acc > e.result {
			return false
		}
		return rec(n[1:], acc + n[0]) || rec(n[1:], acc * n[0])
	}
	return rec(e.numbers[1:], e.numbers[0])
}

func sol1() {
	input := parse("input.txt")
	sum := 0
	for _, v := range input {
		if isValid(v) {
			sum += v.result
		}
	}
	fmt.Println(sum)
}

func concat(a, b int) int {
	digits := int(math.Log10(float64(b))) + 1
	return a * int(math.Pow(10, float64(digits))) + b
}

func isValid2(e eq) bool {
	var rec func([]int, int) bool
	rec = func(n []int, acc int) bool {
		if acc == e.result && len(n) == 0 {
			return true
		}
		if len(n) == 0 || acc > e.result {
			return false
		}
		return rec(n[1:], acc + n[0]) || rec(n[1:], acc * n[0]) || rec(n[1:], concat(acc, n[0]))
	}
	return rec(e.numbers[1:], e.numbers[0])
}

func sol2() {
	input := parse("input.txt")
	sum := 0
	for _, v := range input {
		if isValid2(v) {
			sum += v.result
		}
	}
	fmt.Println(sum)
}

func main() {
	sol1()
	sol2()
}
