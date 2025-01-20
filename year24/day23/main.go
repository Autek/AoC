package main

import (
	"fmt"
	"os"
	"strings"
)

type connection struct {
	a, b string
}

type connections map[connection]struct{}

func newConnections() connections {
	return make(connections)
}

func (connection *connection) has_same_endpoint(c *connection) bool {
	c1, c2 := *connection, *c
	return c1.a == c2.a || c1.a == c2.b || c1.b == c2.a || c1.b == c2.b
}

func nameStartingWith(network []connections, c string) int {
	counter := 0
	for _, connections := range network {
		for co := range connections {
			if strings.HasPrefix(co.a, c) || strings.HasPrefix(co.b, c) {
				counter += 1
				break
			}
		}
	}
	// each triangle is two times in the list so the result is divided by two
	return counter / 2
}

func (connections *connections) add(c *connection) {
	(*connections)[*c] = struct{}{}
}

func (connections *connections) contains_same_endpoint(c *connection) []*connection {
	c1 := *c
	matches := make([]*connection, 0)
	for c2 := range *connections {
		if c1.has_same_endpoint(&c2) {
			matches = append(matches, &c2)
		}
	}
	return matches
}

func (connections *connections) contains(c *connection) bool {
	_, contained := (*connections)[*c]
	return contained
}

func newConnection(a, b string) connection {
	if a < b {
		return connection{a, b}
	}
	return connection{b, a}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func sol1() {
	results := make([]connections, 100)
	input, err := os.ReadFile("input.txt")
	check(err)
	slice_input := strings.Split(strings.TrimSpace(string(input)), "\n")
	connections := newConnections()

	for _, c1_string := range slice_input {
		c1_array := strings.Split(strings.TrimSpace(c1_string), "-")
		c1 := newConnection(c1_array[0], c1_array[1])
		potentials := connections.contains_same_endpoint(&c1)
		for _, c2 := range potentials {
			var e1 string
			if c1.a != c2.a && c1.a != c2.b {
				e1 = c1.a
			} else {
				e1 = c1.b
			}
			var e2 string
			if c2.a != c1.a && c2.a != c1.b {
				e2 = c2.a
			} else {
				e2 = c2.b
			}
			c3 := newConnection(e1, e2)
			if connections.contains(&c3) {
				result := newConnections()
				result.add(&c1)
				result.add(c2)
				result.add(&c3)
				results = append(results, result)
			}
		}
		connections.add(&c1)
	}
	fmt.Println(nameStartingWith(results, "t"))
}

func main() {
	sol1()
}
