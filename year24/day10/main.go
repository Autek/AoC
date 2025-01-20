package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type queue[T any] []T

func (q *queue[T]) enqueue(value T) {
	*q = append(*q, value)
}

func (q *queue[T]) dequeue() T {
	front := (*q)[0]
	*q = (*q)[1:]
	return front
}

func (q queue[T]) isEmpty() bool {
	return len(q) == 0
}

type graph[T any] [][]vertex[T]

type vertex[T any] struct {
	val     T
	reached bool
}

func (g graph[T]) adjacentEdges(indices []int) [][]int {
	maxI := len(g)
	maxJ := len(g[0])
	steps := [][]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	adjacent := make([][]int, 0, len(steps))
	for _, step := range steps {
		edgeI := indices[0] + step[0]
		edgeJ := indices[1] + step[1]

		if edgeI >= 0 && edgeI < maxI && edgeJ >= 0 && edgeJ < maxJ {
			adjacent = append(adjacent, []int{edgeI, edgeJ})
		}
	}
	return adjacent
}

func parse(fileName string) graph[int] {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.TrimSpace(string(file)), "\n")
	grid := make(graph[int], len(lines))
	for i, line := range lines {
		row := make([]vertex[int], len(line))
		for j, r := range line {
			nb, err := strconv.Atoi(string(r))
			if err != nil {
				panic(err)
			}
			v := vertex[int]{val: nb, reached: false}
			row[j] = v
		}
		grid[i] = row
	}
	return grid
}

func score(g graph[int], i, j int) int {
	for _, row := range g {
		for i := range row {
			row[i].reached = false
		}
	}
	summits := 0
	g[i][j].reached = true
	q := queue[[]int]{}
	q.enqueue([]int{i, j})
	for !q.isEmpty() {
		u := q.dequeue()
		for _, indexes := range g.adjacentEdges(u) {
			if !g[indexes[0]][indexes[1]].reached {
				if g[u[0]][u[1]].val+1 == g[indexes[0]][indexes[1]].val {
					g[indexes[0]][indexes[1]].reached = true
					if g[indexes[0]][indexes[1]].val == 9 {
						summits++
					} else {
						q.enqueue(indexes)
					}
				}
			}
		}
	}
	return summits
}

func sol1() {
	graph := parse("input.txt")
	sum := 0
	for i, row := range graph {
		for j, v := range row {
			if v.val == 0 {
				sum += score(graph, i, j)
			}
		}
	}
	fmt.Println(sum)
}

func rating(g graph[int], i, j int) int {
	rating := 0
	q := queue[[]int]{}
	q.enqueue([]int{i, j})
	for !q.isEmpty() {
		u := q.dequeue()
		for _, n := range g.adjacentEdges(u) {
			if g[u[0]][u[1]].val+1 == g[n[0]][n[1]].val {
				if g[n[0]][n[1]].val == 9 {
					rating++
				} else {
					q.enqueue(n)
				}
			}
		}
	}
	return rating
}

func sol2() {
	graph := parse("input.txt")
	sum := 0
	for i, row := range graph {
		for j, v := range row {
			if v.val == 0 {
				sum += rating(graph, i, j)
			}
		}
	}
	fmt.Println(sum)
}

func main() {
	sol1()
	sol2()
}
