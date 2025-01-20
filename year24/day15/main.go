package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	empty = iota
	box
	wall
	robot
)

type cell int

type vec struct {
	x, y int
}

func (v1 *vec) add(v2 vec) {
	v1.x += v2.x
	v1.y += v2.y
}

func (v1 *vec) sub(v2 vec) {
	v1.x -= v2.x
	v1.y -= v2.y
}

func parse(fileName string) ([][]cell, []rune) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	warehouseStr, movesStr, _ := strings.Cut(string(file), "\n\n")

	lines := strings.Split(strings.TrimSpace(warehouseStr), "\n")
	wareHouse := make([][]cell, len(lines))
	for i, line := range lines {
		row := make([]cell, len(line))
		for j, r := range line {
			if r == '#' {
				row[j] = wall
			} else if r == '@' {
				row[j] = robot
			} else if r == 'O' {
				row[j] = box
			} else {
				row[j] = empty
			}
		}
		wareHouse[i] = row
	}

	moves := make([]rune, 0, len(movesStr))
	for _, r := range movesStr {
		if r != '\n' {
			moves = append(moves, r)
		}
	}

	return wareHouse, moves
}

func performMove(posRobot *vec, wareHouse [][]cell, move rune) {
	var moves = map[rune]vec{
		'<': vec{-1, 0},
		'^': vec{0, -1},
		'>': vec{1, 0},
		'v': vec{0, 1},
	}

	posPointer := *posRobot
	for wareHouse[posPointer.y][posPointer.x] != empty {
		posPointer.add(moves[move])
		if wareHouse[posPointer.y][posPointer.x] == wall {
			return
		}
	}

	for wareHouse[posPointer.y][posPointer.x] != robot {
		lastCell := &wareHouse[posPointer.y][posPointer.x]
		posPointer.sub(moves[move])
		*lastCell = wareHouse[posPointer.y][posPointer.x]
	}

	wareHouse[posRobot.y][posPointer.x] = empty
	posRobot.add(moves[move])
}

func findRobot(wareHouse [][]cell) *vec {
	for i := range wareHouse {
		for j := range wareHouse[i] {
			if wareHouse[i][j] == robot {
				return &vec{j, i}
			}
		}
	}
	return nil
}

func sumGPS(wareHouse [][]cell) int {
	sum := 0
	for i := range wareHouse {
		for j := range wareHouse[i] {
			if wareHouse[i][j] == box {
				sum += i*100 + j
			}
		}
	}
	return sum
}

func sol1() {
	wareHouse, moves := parse("input.txt")
	robotPos := findRobot(wareHouse)
	for _, m := range moves {
		performMove(robotPos, wareHouse, m)
	}
	sum := sumGPS(wareHouse)
	fmt.Println(sum)
}

type cell2 struct {
	isSimple bool
	simple   cell
	box      *cell2
}

type Node[T any] struct {
	val         T
	left, right *Node[T]
}

func parse2(fileName string) ([][]cell2, []rune) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	warehouseStr, movesStr, _ := strings.Cut(string(file), "\n\n")

	lines := strings.Split(strings.TrimSpace(warehouseStr), "\n")
	wareHouse := make([][]cell2, len(lines))
	for i, line := range lines {
		row := make([]cell2, len(line)*2)
		for j, r := range line {
			if r == '#' {
				row[j*2] = cell2{isSimple: true, simple: wall}
				row[j*2+1] = cell2{isSimple: true, simple: wall}
			} else if r == '@' {
				row[j*2] = cell2{isSimple: true, simple: robot}
				row[j*2+1] = cell2{isSimple: true, simple: empty}
			} else if r == 'O' {
				row[j*2] = cell2{isSimple: false}
				row[j*2+1] = cell2{isSimple: false, box: &row[j*2]}
				row[j*2].box = &row[j*2+1]
			} else {
				row[j*2] = cell2{isSimple: true, simple: empty}
				row[j*2+1] = cell2{isSimple: true, simple: empty}
			}
		}
		wareHouse[i] = row
	}

	moves := make([]rune, 0, len(movesStr))
	for _, r := range movesStr {
		if r != '\n' {
			moves = append(moves, r)
		}
	}

	return wareHouse, moves
}

func findRobot2(wareHouse [][]cell2) *vec {
	for i := range wareHouse {
		for j := range wareHouse[i] {
			if wareHouse[i][j].isSimple && wareHouse[i][j].simple == robot {
				return &vec{j, i}
			}
		}
	}
	return nil
}

func sumGPS2(wareHouse [][]cell2) int {
	sum := 0
	for i := range wareHouse {
		for j := 0; j < len(wareHouse[i]); j++ {
			if !wareHouse[i][j].isSimple {
				sum += i*100 + j
				j++
			}
		}
	}
	return sum
}

func getMoves(wareHouse [][]cell2, pos vec, move vec) (bool, *Node[vec]) {
	nextPos := vec{pos.x + move.x, pos.y + move.y}
	nextCell := wareHouse[nextPos.y][nextPos.x]
	if nextCell.isSimple &&
		nextCell.simple == wall {
		return false, nil
	}

	root := Node[vec]{val: pos}
	if nextCell.isSimple && nextCell.simple == empty {
		return true, &root
	}

	canMove, moves := getMoves(wareHouse, nextPos, move)
	root.left = moves

	// if the movement is horizontal, the other part of the box is getting
	// moved by pushing, no need to move it by "linking"
	if move.y == 0 {
		return canMove, &root
	}

	// nextPos will be the other part of the box, it is either to
	// the left or right
	if &wareHouse[nextPos.y][nextPos.x-1] == nextCell.box {
		nextPos.sub(vec{1, 0})
	} else {
		nextPos.add(vec{1, 0})
	}

	canMove2, moves2 := getMoves(wareHouse, nextPos, move)
	root.right = moves2

	return canMove && canMove2, &root
}

func performMoves(wareHouse [][]cell2, toMove *Node[vec], move vec) {
	if toMove.left != nil {
		performMoves(wareHouse, toMove.left, move)
	}
	if toMove.right != nil {
		performMoves(wareHouse, toMove.right, move)
	}
	pos := toMove.val
	wareHouse[pos.y+move.y][pos.x+move.x] = wareHouse[pos.y][pos.x]
	if !wareHouse[pos.y][pos.x].isSimple {
		wareHouse[pos.y+move.y][pos.x+move.x].box.box = &wareHouse[pos.y+move.y][pos.x+move.x]
	}
	wareHouse[pos.y][pos.x] = cell2{isSimple: true, simple: empty}
}

// remove duplicates from the moves, could probably
// be replaced by a visited set in getMoves
func removeDuplicates(moves *Node[vec]) {
	reached := map[vec]struct{}{}
	var rec func(n *Node[vec])
	rec = func(n *Node[vec]) {
		c1 := n.left
		c2 := n.right
		if c1 != nil {
			if _, ok := reached[c1.val]; ok {
				n.left = nil
			} else {
				reached[c1.val] = struct{}{}
				rec(c1)
			}
		}
		if c2 != nil {
			if _, ok := reached[c2.val]; ok {
				n.right = nil
			} else {
				reached[c2.val] = struct{}{}
				rec(c2)
			}
		}
	}
	rec(moves)
}

func move2(posRobot *vec, wareHouse [][]cell2, move rune) {
	moves := map[rune]vec{
		'^': vec{0, -1},
		'v': vec{0, 1},
		'<': vec{-1, 0},
		'>': vec{1, 0},
	}

	possible, toMove := getMoves(wareHouse, *posRobot, moves[move])
	if !possible {
		return
	}
	removeDuplicates(toMove)
	performMoves(wareHouse, toMove, moves[move])
	posRobot.add(moves[move])
}

func sol2() {
	wareHouse, moves := parse2("input.txt")
	robotPos := findRobot2(wareHouse)
	for _, m := range moves {
		move2(robotPos, wareHouse, m)
	}
	sum := sumGPS2(wareHouse)
	fmt.Println(sum)
}

func main() {
	sol1()
	sol2()
}

// ----------------- Debugging -----------------

func printWarehouse(wareHouse [][]cell2) {
	str := ""
	for i := range wareHouse {
		for j := range wareHouse[i] {
			if !wareHouse[i][j].isSimple {
				if wareHouse[i][j].box == &wareHouse[i][j-1] {
					str += "]"
				} else {
					str += "["
				}
			} else if wareHouse[i][j].simple == empty {
				str += "."
			} else if wareHouse[i][j].simple == wall {
				str += "#"
			} else if wareHouse[i][j].simple == robot {
				str += "@"
			}
		}
		str += "\n"
	}
	fmt.Print(str)
}

func testGrid(w [][]cell2) bool {
	for i := range w {
		for j := range w[i] {
			if !w[i][j].isSimple {
				if w[i][j].box != &w[i][j-1] && w[i][j].box != &w[i][j+1] {
					return false
				}
			}
		}
	}
	return true
}

func printMoves(m Node[vec]) {
	moves := []Node[vec]{m}
	for i := 0; len(moves) != 0; i++ {
		fmt.Println("level", i)
		str := ""
		newMoves := []Node[vec]{}
		for _, n := range moves {
			str += fmt.Sprintf("%v", n.val)
			if n.left != nil {
				newMoves = append(newMoves, *n.left)
			}
			if n.right != nil {
				newMoves = append(newMoves, *n.right)
			}
		}
		moves = newMoves
		fmt.Println(str)
	}
}
