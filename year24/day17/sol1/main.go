package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var instructions = []instruction{
	adv{},
	bxl{},
	bst{},
	jnz{},
	bxc{},
	out{},
	bdv{},
	cdv{},
}

type register int
const(
	a = register(iota)
	b
	c
)

type instruction interface{
	apply(int, *int, *[]int, []int)
}

func mapOperand(i int, memory []int) int{
	switch {
	case i >= 0 && i <= 3:
		return i
	case i == 4:
		return memory[a]
	case i == 5:
		return memory[b]
	case i == 6:
		return memory[c]
	}
	panic("unknown operand")
}

type adv struct {}
func (adv)apply(operand int, ipt *int, _ *[]int, memory []int) {
	memory[a] >>= mapOperand(operand, memory)
	*ipt += 2
}

type bxl struct {}
func (bxl)apply(operand int, ipt *int, _ *[]int, memory []int) {
	memory[b] ^= operand
	*ipt += 2
}

type bst struct {}
func (bst)apply(operand int, ipt *int, _ *[]int, memory []int) {
	memory[b] = mapOperand(operand, memory) & 0b111
	*ipt += 2
}

type jnz struct {}
func (jnz)apply(operand int, ipt *int, _ *[]int, memory []int) {
	if memory[a] == 0 {
		*ipt += 2
		return
	}
	*ipt = operand
}

type bxc struct {}
func (bxc)apply(_ int, ipt *int, _ *[]int, memory []int) {
	memory[b] ^= memory[c]
	*ipt += 2
}

type out struct {}
func (out)apply(operand int, ipt *int, output *[]int, memory []int) {
	*output = append(*output, mapOperand(operand, memory) & 0b111)
	*ipt += 2
}

type bdv struct {}
func (bdv)apply(operand int, ipt *int, _ *[]int, memory []int) {
	memory[b] = memory[a] >> mapOperand(operand, memory)
	*ipt += 2
}

type cdv struct {}
func (cdv)apply(operand int, ipt *int, _ *[]int, memory []int) {
	memory[c] = memory[a] >> mapOperand(operand, memory)
	*ipt += 2
}

func parse(fileName string) (int, int, int, []int) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	var a, b, c int
	p := ""

	format := 
		`Register A: %d
		Register B: %d
		Register C: %d

		Program: %s`

	fmt.Sscanf(strings.TrimSpace(string(file)), format, &a, &b, &c, &p)
	nbs := strings.Split(p, ",")
	program := make([]int, len(nbs))
	for i, c := range nbs {
		nb, err := strconv.Atoi(c)
		if err != nil {
			panic(err)
		}
		program[i] = nb
	}
	return a, b, c, program
}

func runProgram(program []int, output *[]int, memory []int) {
	iptr := 0
	programLen := len(program)
	for iptr < programLen {
		inst := instructions[program[iptr]]
		ptr := &iptr
		val := program[iptr+1]
		inst.apply(val, ptr, output, memory)
	}
}

func stdout(output []int) {
	out := strconv.Itoa(output[0])
	for _, val := range output[1:] {
		out += ","
		out += strconv.Itoa(val)
	}
	fmt.Println(out)
}

func sol() {
	A, B, C, p := parse("../input.txt")
	memory := make([]int, 3)
	memory[a] = A
	memory[b] = B
	memory[c] = C
	output := []int{}
	runProgram(p, &output, memory)
	stdout(output)
}

func main() {
	sol()
}
