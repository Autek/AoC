package main

import (
	"fmt"
	"math"
	"slices"
)

func runProgram2(program []int, memory []int, output *[]int) {
	goodOutput := true
	iptr := 0
	getCombo := func(operand int) int {
		switch operand {
		case 4:
			return memory[0]
		case 5:
			return memory[1]
		case 6:
			return memory[2]
		default:
			return operand
		}
	}
	instructions := []func(int){
		func(operand int) {
			memory[a] >>= getCombo(operand)
		},
		func(operand int) {
			memory[b] ^= operand
		},
		func(operand int) {
			memory[b] = getCombo(operand) & 0x7
		},
		func(operand int) {
			if memory[a] != 0 {
				iptr = operand - 2
			}
		},
		func(operand int) {
			memory[b] ^= memory[c]
		},
		func(operand int) {
			val := getCombo(operand) & 0x7
			*output = append(*output, val)
			goodOutput = val == program[len(*output) - 1]
		},
		func(operand int) {
			memory[b] = memory[a] >> getCombo(operand)
		},
		func(operand int) {
			memory[c] = memory[a] >> getCombo(operand)
		},
	}
	for goodOutput && iptr < len(program) {
		instructions[program[iptr]](program[iptr + 1])
		iptr += 2
	}
}

func findRegisterA(program[]int, B, C int) int{
	nbRoutines := 8
	solution := make(chan int)
	stop := make(chan struct{})
	for i := 0; i <= nbRoutines; i++ {
		go func(start int) {
			output := []int{}
			memory := make([]int, 3)
			for j := start; ;j += nbRoutines {
				select {
				case <- stop:
					return
				default:
					memory[a] = j
					memory[b] = B
					memory[c] = C
					runProgram2(program, memory, &output)
					if j % 1000000000 == 0 {
						fmt.Printf("%.4g%% done.\n", float64(j)/math.MaxInt64 * 100)
					}
					if slices.Equal(program, output) {
						solution <- j
						close(stop)
						return
					}
					output = output[:0]
				}
			}
		}(i)
	}
	return <- solution
}

func sol2() {
	_, B, C, p := parse("input.txt")
	fmt.Println(findRegisterA(p, B, C))
}
