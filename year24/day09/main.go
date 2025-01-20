package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ID int

type option[T any] struct {
	value T
	valid bool
}

func some[T any](value T) option[T] {
	return option[T]{value: value, valid: true}
}

func none[T any]() option[T] {
	return option[T]{valid: false}
}

func (o option[T]) isSome() bool {
	return o.valid
}

func (o option[T]) isNone() bool {
	return !o.valid
}

func (o option[T]) unwrap() T {
	if o.isNone() {
		panic("unwrap: none")
	}
	return o.value
}

func (o option[T]) unwrapOr(def T) T {
	if o.isSome() {
		return o.value
	}
	return def
}

func parse(fileName string) []option[ID] {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	isFile := true
	id := ID(0)
	disk := []option[ID]{}
	for _, r := range strings.TrimSpace(string(file)) {
		nb, err := strconv.Atoi(string(r))
		if err != nil {
			panic(err)
		}
		var e option[ID]
		if isFile {
			e = some(id)
			id++
		} else {
			e = none[ID]()
		}
		for i := 0; i < nb; i++ {
			disk = append(disk, e)
		}
		isFile = !isFile
	}
	return disk
}

func printDisk(disk []option[ID]) {
	str := ""
	for _, o := range disk {
		if o.isSome() {
			str += strconv.Itoa((int(o.unwrap())))
		} else {
			str += "."
		}
	}
	fmt.Println(str)
}

func compact(disk []option[ID]) {
	leftStart := 0
	for i := len(disk) - 1; i > leftStart; i-- {
		if disk[i].isNone() {
			continue
		}
		for j := leftStart; j < i; j++ {
			if disk[j].isNone() {
				disk[j] = disk[i]
				disk[i] = none[ID]()
				leftStart = j + 1
				break
			}
		}
	}
}

func checkSum(disk []option[ID]) int {
	sum := 0
	for i, o := range disk {
		if o.isNone() {
			break
		}
		sum += i * int(o.unwrap())
	}
	return sum
}

func sol1() {
	disk := parse("input.txt")
	compact(disk)
	fmt.Println(checkSum(disk))
}

func getSize(disk []option[ID], dir int, start int, isFile bool) int {
	size := 0
	val := none[ID]()
	for i := start; i < len(disk) && i >= 0; i += dir {
		if i < 0 || i >= len(disk) || disk[i].isSome() != isFile {
			break
		}
		if isFile && val.isSome() && val.unwrap() != disk[i].unwrap() {
			break
		}
		val = disk[i]
		size++
	}
	return size
}

func swap(disk []option[ID], from int, to int, size int) {
	for i := 0; i < size; i++ {
		disk[to+i] = disk[from+i]
		disk[from+i] = none[ID]()
	}
}

func compact2(disk []option[ID]) {
	for i := len(disk) - 1; i > 0; {
		if disk[i].isNone() {
			i -= getSize(disk, -1, i, false)
			continue
		}
		fileSize := getSize(disk, -1, i, true)
		for j := 0; j < i; {
			if disk[j].isSome() {
				j += getSize(disk, 1, j, true)
				continue
			}
			gapSize := getSize(disk, 1, j, false)
			if gapSize >= fileSize {
				swap(disk, i-fileSize+1, j, fileSize)
				break
			}
			j += gapSize
		}
		i -= fileSize
	}
}

func checkSum2(disk []option[ID]) int {
	sum := 0
	for i, o := range disk {
		sum += i * int(o.unwrapOr(0))
	}
	return sum
}

func sol2() {
	disk := parse("input.txt")
	compact2(disk)
	fmt.Println(checkSum2(disk))
}

func main() {
	sol1()
	sol2()
}
