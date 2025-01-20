// in this first day of the advent of code I had fun with functional programing
// concepts in go and this made me use a lot of generics yay!
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseInput(fileName string) ([]int, []int) {
	l1 := make([]int, 0)
	l2 := make([]int, 0)

	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
	for _, line := range lines {
		splitLine := strings.Split(strings.TrimSpace(line), "   ")
		nb1, err := strconv.Atoi(strings.TrimSpace(splitLine[0]))
		if err != nil {
			panic(err)
		}

		nb2, err := strconv.Atoi(strings.TrimSpace(splitLine[1]))
		if err != nil {
			panic(err)
		}
		l1 = append(l1, nb1)
		l2 = append(l2, nb2)
	}
	return l1, l2
}

func transform[T1 any, T2 any](l1 []T1, f func(T1) T2) []T2 {
	l2 := make([]T2, len(l1))
	for i, val := range l1 {
		l2[i] = f(val)
	}
	return l2
}

func filter[T any](l1 []T, f func(T) bool) []T {
	l2 := make([]T, 0)
	for _, val := range l1 {
		if f(val) {
			l2 = append(l2, val)
		}
	}
	return l2
}

func mergeSort(l []int) []int {
	merge := func(l1 []int, l2 []int) []int {
		i, j := 0, 0
		l3 := make([]int, 0, len(l1)+len(l2))
		for i < len(l1) || j < len(l2) {
			if i >= len(l1) {
				l3 = append(l3, l2[j])
				j++
			} else if j >= len(l2) {
				l3 = append(l3, l1[i])
				i++
			} else if l1[i] <= l2[j] {
				l3 = append(l3, l1[i])
				i++
			} else {
				l3 = append(l3, l2[j])
				j++
			}
		}
		return l3
	}
	base := transform(l, func(e int) []int { return []int{e} })
	for len(base) != 1 {
		for i := 0; i < len(base); i += 2 {
			a1 := base[i]
			var a2 []int
			if i+1 >= len(base) {
				a2 = nil
			} else {
				a2 = base[i+1]
			}
			base[i/2] = merge(a1, a2)
		}
		base = base[:(len(base)+1)/2]
	}
	return base[0]
}

func zip[T any](l1 []T, l2 []T) [][]T {
	if len(l1) != len(l2) {
		panic("l1 and l2 have different length")
	}
	l3 := make([][]T, len(l1))
	for i := range l3 {
		l3[i] = []T{l1[i], l2[i]}
	}
	return l3
}

func abs(nb int) int {
	if nb < 0 {
		return -nb
	}
	return nb
}

func foldLeft[T any, I any](l []T, init I, f func(T, I) I) I {
	acc := init
	for _, val := range l {
		acc = f(val, acc)
	}
	return acc
}

func group[T comparable](l []T) map[T]int {
	f := func(e T, m map[T]int) map[T]int {
		m[e]++
		return m
	}
	return foldLeft(l, map[T]int{}, f)
}

func sol1() {
	l1, l2 := parseInput("input.txt")
	l1 = mergeSort(l1)
	l2 = mergeSort(l2)
	l3 := zip(l1, l2)
	l4 := transform(l3, func(tuple []int) int { return tuple[0] - tuple[1] })
	l5 := transform(l4, abs)
	result := foldLeft(l5, 0, func(a int, b int) int { return a + b })
	fmt.Println(result)
}

func sol2() {
	l1, l2 := parseInput("input.txt")
	group1 := group(l1)
	group2 := group(l2)
	sum := 0
	for key := range group1 {
		sum += key * group2[key]
	}
	fmt.Println(sum)
}

func main() {
	sol2()
}
