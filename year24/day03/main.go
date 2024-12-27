// this puzzle seems to be typical use case for a regex so let's try to use
// regex in go !


package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)


func filterInput(input string) []string{
	regex := regexp.MustCompile(`(mul\([0-9]{1,3},[0-9]{1,3}\))`)
	return regex.FindAllString(input, -1)
}

func evaluateMul(input string) int{
	regex := regexp.MustCompile(`([0-9]{1,3})`)
	numbers_str := regex.FindAllString(input, 2)
	nb1, err := strconv.Atoi(numbers_str[0]) 
	if err != nil {
		panic(err)
	}
	nb2, err := strconv.Atoi(numbers_str[1]) 
	if err != nil {
		panic(err)
	}
	return nb1 * nb2
}

func sol1() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	filtered := filterInput(string(input))
	sum := 0
	for _, mul := range filtered {
		sum += evaluateMul(mul)
	}
	fmt.Println(sum)
}

func filterInput2(input string) []string{
	regex := regexp.MustCompile(`(mul\([0-9]{1,3},[0-9]{1,3}\))|(don't)|(do)`)
	return regex.FindAllString(input, -1)
}

func sol2() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	filtered := filterInput2(string(input))
	sum := 0
	do := true
	for _, opp := range filtered {
		if opp == "do" {
			do = true
		} else if opp == "don't" {
			do = false
		} else if do {
			sum += evaluateMul(opp)
		}
	}
	fmt.Println(sum)
}

func main() {
	sol2()
}
