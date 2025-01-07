package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

type vec struct {
	a, b float64
}

type matrix struct {
	a, b, c, d float64
}

func (m *matrix)inverse() matrix {
	det := m.a * m.d - m.b * m.c    
	if math.Abs(det) < 1e-9 {
        panic("Matrix is singular or nearly singular")
    }

	a := m.d / det
	b := -m.b / det
	c := -m.c / det
	d := m.a / det
	return matrix{a, b, c, d}
}

func matXVec(m *matrix, v *vec) vec {
	a := m.a * v.a + m.b * v.b
	b := m.c * v.a + m.d * v.b
	return vec{a, b}
}

func matFromVec(v1, v2 *vec) matrix {
	return matrix{v1.a, v2.a, v1.b, v2.b}
}

func parse(filename string) ([]vec, []vec, []vec) {
	file, err := os.ReadFile(filename) 
	if err != nil {
		panic(err)
	}
	machines := strings.Split(strings.TrimSpace(string(file)), "\n\n")
	v1s := []vec{}
	v2s := []vec{}
	goals := []vec{}
	for _, machineStr := range machines {
		var v1a, v1b, v2a, v2b, goala, goalb float64
		format := "Button A: X+%f, Y+%f\nButton B: X+%f, Y+%f\nPrize: X=%f, Y=%f"
		fmt.Sscanf(machineStr, format, &v1a, &v1b, &v2a, &v2b, &goala, &goalb)
		v1s = append(v1s, vec{v1a, v1b})
		v2s = append(v2s, vec{v2a, v2b})
		goals = append(goals, vec{goala, goalb})
	}
	return v1s, v2s, goals
}

func closeToInt(val, threshold float64) (int, bool) {
	if (math.Abs(math.Round(val) - val) < threshold) {
		return int(math.Round(val)), true
	}
	return -1.0, false 
}

// computes the change of basis from the custom basis v1, v2 to the 
// cannonical basis
func computePresses(v1, v2, goal vec) (int, int, bool) {
	mat := matFromVec(&v1, &v2)
	matInv := mat.inverse()
	res := matXVec(&matInv, &goal)
	coeffAFloat := res.a
	coeffBFloat := res.b
	if coeffAFloat >= 0 && coeffBFloat >= 0 {
		threshold := 1e-4
		coeffA, isClose1 := closeToInt(coeffAFloat, threshold)
		coeffB, isClose2 := closeToInt(coeffBFloat, threshold)
		if isClose1 && isClose2 {
			return coeffA, coeffB, true 
		}
	}
	return -1, -1, false
}

func sol1() {
	v1, v2, goal := parse("input.txt")
	sum := 0
	for i := range v1 {
		a, b, possible := computePresses(v1[i], v2[i], goal[i])
		if possible {
			sum += 3*a + b
		}
	}

	fmt.Println(sum)
}

func sol2() {
	v1, v2, goal := parse("input.txt")
	conversionError := float64(10000000000000)
	sum := 0
	for i := range v1 {
		g := vec{goal[i].a + conversionError, goal[i].b + conversionError}
		a, b, possible := computePresses(v1[i], v2[i], g)
		if possible {
			sum += 3*a + b
		}
	}
	fmt.Println(sum)
}

func main() {
	sol1()
	sol2()
}
