package main

import (
	"fmt"
	"os"
	"strings"
)

type stripe int

const (
	white = stripe(iota)
	blue
	black
	red
	green
)

type pattern []stripe

func parse(fileName string) ([]pattern, []pattern) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	stripeMap := map[rune]stripe{
		'w': white,
		'u': blue,
		'b': black,
		'r': red,
		'g': green,
	}

	p, d, _ := strings.Cut(strings.TrimSpace(string(file)), "\n\n")

	patternsStr := strings.Split(p, ", ")
	patterns := make([]pattern, len(patternsStr))
	for i, patternStr := range patternsStr {
		patt := make(pattern, len(patternStr))
		for j, r := range patternStr {
			patt[j] = stripeMap[r]
		}
		patterns[i] = patt
	}

	designsStr := strings.Split(d, "\n")
	designs := make([]pattern, len(designsStr))
	for i, designStr := range designsStr {
		design := make(pattern, len(designStr))
		for j, r := range designStr {
			design[j] = stripeMap[r]
		}
		designs[i] = design
	}
	return patterns, designs
}

func startsWidth[T comparable](s []T, prefix []T) bool {
	if len(s) < len(prefix) {
		return false
	}

	for i := range prefix {
		if s[i] != prefix[i] {
			return false
		}
	}
	return true
}

func solve(patterns []pattern, designs []pattern) []int {
	res := make([]int, len(designs))
	for i, design := range designs {
		acc := make([]int, len(design)+1)
		acc[0] = 1
		for j := range acc {
			for _, patt := range patterns {
				if startsWidth(design[j:], patt) {
					acc[j+len(patt)] += acc[j]
				}
			}
		}
		res[i] = acc[len(design)]
	}
	return res
}

func sol1() {
	patterns, designs := parse("input.txt")
	waysToReach := solve(patterns, designs)
	sum := 0
	for _, val := range waysToReach {
		if val != 0 {
			sum++
		}
	}
	fmt.Println(sum)
}

func sol2() {
	patterns, designs := parse("input.txt")
	waysToReach := solve(patterns, designs)
	sum := 0
	for _, val := range waysToReach {
		sum += val
	}
	fmt.Println(sum)
}

func main() {
	sol1()
	sol2()
}
