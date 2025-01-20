package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type relation struct {
	from, to int
}

type update []int

func parse(fileName string) (map[relation]struct{}, []update) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	relationsUpdates := strings.Split(strings.TrimSpace(string(file)), "\n\n")
	relations := strings.TrimSpace(relationsUpdates[0])
	updates := strings.TrimSpace(relationsUpdates[1])
	splitRelations := strings.Split(relations, "\n")
	relMap := map[relation]struct{}{}
	for _, rel_str := range splitRelations {
		relNb := strings.Split(rel_str, "|")
		from, err := strconv.Atoi(relNb[0])
		if err != nil {
			panic(err)
		}
		to, err := strconv.Atoi(relNb[1])
		if err != nil {
			panic(err)
		}
		rel := relation{from, to}
		relMap[rel] = struct{}{}
	}

	splitUpdates := strings.Split(updates, "\n")
	updateList := make([]update, len(splitUpdates))
	for i, upd_str := range splitUpdates {
		nbList := strings.Split(strings.TrimSpace(upd_str), ",")
		upd := make(update, len(nbList))
		for j, nb_str := range nbList {
			nb, err := strconv.Atoi(nb_str)
			if err != nil {
				panic(err)
			}
			upd[j] = nb
		}
		updateList[i] = upd
	}
	return relMap, updateList
}

func sol1() {
	relations, updates := parse("input.txt")
	sum := 0
	comparator := func(a int, b int) int {
		if _, present := relations[relation{a, b}]; present {
			return -1
		}
		if _, present := relations[relation{b, a}]; present {
			return 1
		}
		return 0
	}

	for _, u := range updates {
		if slices.IsSortedFunc(u, comparator) {
			sum += u[(len(u)-1)/2]

		}
	}
	fmt.Println(sum)
}

func sol2() {
	relations, updates := parse("input.txt")
	sum := 0
	comparator := func(a int, b int) int {
		if _, present := relations[relation{a, b}]; present {
			return -1
		}
		if _, present := relations[relation{b, a}]; present {
			return 1
		}
		return 0
	}

	for _, u := range updates {
		if !slices.IsSortedFunc(u, comparator) {
			slices.SortFunc(u, comparator)
			sum += u[(len(u)-1)/2]
		}
	}
	fmt.Println(sum)
}

func main() {
	sol1()
	sol2()
}
