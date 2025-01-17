package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func parse(fileName string) []int {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	pString := strings.TrimSpace(strings.Split(string(file), "\n\n")[1])
	nbs := strings.Split(strings.TrimPrefix(pString, "Program: "), ",")
	program := make([]int, len(nbs))
	for i, c := range nbs {
		nb, err := strconv.Atoi(c)
		if err != nil {
			panic(err)
		}
		program[i] = nb
	}
	return program
}

// closed form formula for the program without the loop at the end 
// and without updating the value of the register A (which is A >> 3)
func closedForm(A int) int{
	B := (A & 0b111) ^ 3
	C := A >> B
	tmp := (B ^ 5 ^ C)
	return tmp & 0b111
}

func findA(program []int) int {
	outputs := make([][]int, 8, 8)

	// find all possible values of A that lead to a given output 
	// we only need to check the first 11 bits of A since after that the 
	// bits are not used in the program
	for i := range outputs {
		for A := 0; A < 1 << 11; A++ {
			value := closedForm(A)
			if value == i {
				outputs[i] = append(outputs[i], A)
			}
		}
	}

	workingNumbers := make(map[int]struct{})
	// initialize the working numbers with the possible values 
	// of the register A to get the last value of the program
	for _, nb := range outputs[program[len(program)-1]] {	
		workingNumbers[nb] = struct{}{}
	}

	for i := len(program) - 2; i >= 0; i-- {
		newNumbers := make(map[int]struct{})
		for nb := range workingNumbers {
			for _, next := range outputs[program[i]] {
				potentialA := nb << 3 | next
				tag := true
				// check if the value of the register A 
				// does lead to the correct output.
				// we are shifting the bits of A to the right 3 times
				// because between each iteration of the loop the value of A
				// is shifted to the right 3 times
				for j := i; j < len(program); j++ {
					if closedForm(potentialA >> (3 * (j-i))) != program[j] {
						tag = false
					}
				}
				// check if the the value of register A outputs a value of
				// the correct length (0 means that the loop is finished)
				if potentialA >> (3 * (len(program))) != 0 {
					tag = false
				}
				if tag {
					newNumbers[potentialA] = struct{}{}
				}
			}
		}
		workingNumbers = newNumbers
	}

	m := math.MaxInt64
	for k := range workingNumbers {
		if k < m {
			m = k
		}
	}
	return m
}

func sol() {
	p := parse("../input.txt")
	fmt.Println(findA(p))
}

func main() {
	sol()
}
