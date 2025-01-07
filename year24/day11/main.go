package main

import (
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

type stone int

func parse(fileName string) []stone {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	nb_slice := strings.Split(strings.TrimSpace(string(file)), " ")
	result := make([]stone, len(nb_slice))
	for i, nb_str := range nb_slice {
		nb, err := strconv.Atoi(nb_str)
		if err != nil {
			panic(err)
		}
		result[i] = stone(nb)
	}
	return result
}

func conditionRule1(s stone) bool{
	return s == 0
}

func conditionRule2(s stone) bool{
	length := int(math.Log10(float64(int(s)))) + 1
	return length & 1 == 0
}

func conditionRule3(s stone) bool {
	return true
}

func rule1(s stone) stone {
	return stone(1)
}

func rule2(s stone) (stone, stone) {
	numDigits := int(math.Log10(float64(s))) + 1

	mid := numDigits / 2

	divisor := int(math.Pow(10, float64(mid)))

	left := stone(int(s) / divisor)
	right := stone(int(s) % divisor)

	return left, right
}

func rule3(s stone) stone {
	return s * 2024
}

func blink(stones []stone) []stone{
	for i := 0; i < len(stones); i++ {

		if conditionRule1(stones[i]) {
			stones[i] =  rule1(stones[i])

		} else if  conditionRule2(stones[i]) {
			l, r := rule2(stones[i])
			stones = slices.Replace(stones, i, i + 1, l, r)
			i++

		} else if conditionRule3(stones[i]) {
			stones[i] = rule3(stones[i]) 
		}
	}
	return stones
}

func sol1() {
	stones := parse("input.txt")
	blink_numbers := 25
	for i := 0; i < blink_numbers; i++ {
		stones = blink(stones)
	}
	fmt.Println(len(stones))
}

func blink2(stones []stone, blinks int) int{
	mem := make(map[stone][]int, 1 << 10) // 2^10 (might be slightly off)
	var rec func(s stone, b int) int
	rec = func(s stone, b int) int {
		sum := 0
		if b <= 0 {
			sum = 1
		} else if values, present := mem[s]; present && values[b] > 0{
			return values[b]

		}	else if conditionRule1(s) {
			sum = rec(rule1(s), b - 1)

		} else if  conditionRule2(s) {
			l, r := rule2(s)
			sum = rec(l, b - 1) + rec(r, b - 1)

		} else if conditionRule3(s) {
			sum = rec(rule3(s), b - 1)
		}

		if _, present := mem[s]; !present {
			mem[s] = make([]int, blinks + 1)
		}

		mem[s][b] = sum
		return sum
	}
	sum := 0
	for _, s := range stones {
		sum += rec(s, blinks)
	}
	return sum
}

func sol2() {
	stones := parse("input.txt")
	blinks := 75
	fmt.Println(blink2(stones, blinks))
}

func blink3(stones []stone, blinks int) int{

	mem := make(map[stone]int, 1 << 10)
	for _, s := range stones {
		mem[s]++
	}
	newMem := make(map[stone]int, 1 << 10)


	for i := 0; i < blinks; i++ {
		for key := range newMem {
			delete(newMem, key)
		}
		for k, v := range mem {
			if conditionRule1(k) {
				newMem[rule1(k)] += v
			} else if conditionRule2(k) {
				l, r := rule2(k)
				newMem[l] += v
				newMem[r] += v
			} else if conditionRule3(k) {
				newMem[rule3(k)] += v
			}
		}
		mem, newMem = newMem, mem
	}

	sum := 0
	for _, v := range mem {
		sum += v
	}

	return sum
}

func sol3() {
	stones := parse("input.txt")
	blinks := 75
	fmt.Println(blink3(stones, blinks))
}

func main() {
	sol1()
	sol2()
	sol3()
}
