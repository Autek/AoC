package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type logic int

const(
	OR = iota
	AND
	XOR
)

func GetlogicFromStr(name string) logic {
	switch name {
	case "OR": return OR
	case "AND": return AND
	case "XOR": return XOR
	default: panic("name not found \"" + name + "\"")
	}
}

type wireValue int

const(
	ONE = iota
	ZERO
	UNDEFINED
)

func GetWireValueFromStr(name string) wireValue {
	switch name {
	case "1": return ONE
	case "0": return ZERO
	default: panic("name not found \"" + name + "\"")
	}
}

type wire struct{
	name string
	val wireValue
}

type gate struct{
	left, right, output *wire
	logic
}

func getWireFromName(name string, wires []wire) *wire{
	for _, w := range wires {
		if w.name == name {
			return &w
		}
	}
	return nil
}

func getOrCreateWire(key string, wires map[string]*wire) *wire {
	if w, present := wires[key]; present {
		return w
	}
	newWire := &wire{name: key, val: UNDEFINED}
	wires[key] = newWire
	return newWire
}

func parseInput(fileName string) ([]gate, map[string]*wire){
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	file_str := strings.TrimSpace(string(file))
	split_file_str := strings.Split(file_str, "\n\n")

	if (len(split_file_str) != 2) {
		panic("the file was not split in 2")
	}
	wires_str := strings.TrimSpace(split_file_str[0])
	gates_str := strings.TrimSpace(split_file_str[1])
	wires_slice := strings.Split(wires_str, "\n")
	gates_slice := strings.Split(gates_str, "\n")
	
	wires := make(map[string]*wire, 0)
	gates := make([]gate, 0)
	for _, g := range gates_slice {
		gateSlice := strings.Split(g, " -> ")
		if (len(gateSlice) != 2) {
			panic("couldn't split gate in two")
		}
		inputsAndLogic := strings.TrimSpace(gateSlice[0])
		outputStr := strings.TrimSpace(gateSlice[1])
		inputsAndLogicSlice := strings.Split(inputsAndLogic, " ")
		if (len(inputsAndLogicSlice) != 3) {
			panic("couldn't split gate input in three")
		}
		leftStr := strings.TrimSpace(inputsAndLogicSlice[0])
		logicStr := strings.TrimSpace(inputsAndLogicSlice[1])
		rightStr := strings.TrimSpace(inputsAndLogicSlice[2])

		left := getOrCreateWire(leftStr, wires)
		right := getOrCreateWire(rightStr, wires)
		output := getOrCreateWire(outputStr, wires)
		logic := GetlogicFromStr(logicStr)
		builtGate := gate{left: left, right: right, output: output, logic: logic}
		gates = append(gates, builtGate)
	}

	for _, w := range wires_slice {
		wireSlice := strings.Split(w, ": ")
		if (len(wireSlice) != 2) {
			panic("couldn't split wire in two")
		}
		wireName := strings.TrimSpace(wireSlice[0])
		wireValueStr := strings.TrimSpace(wireSlice[1])
		currentWire := wires[wireName]
		wireValue := GetWireValueFromStr(wireValueStr)
		currentWire.val = wireValue
	}
	return gates, wires
}

func update (gates []gate) {
	finished := false
	for !finished{
		finished = true
		for _, g := range gates {
			if updateGate(&g) {
				finished = false
			}
		}
	}
}

func updateGate(g *gate) bool{
	if g == nil {
		panic("recieved nil pointer")
	}
	if g.left == nil {
		panic("left is nil")
	}
	if g.right == nil {
		panic("right is nil")
	}
	if g.left.val != UNDEFINED && 
	g.right.val != UNDEFINED && 
	g.output.val == UNDEFINED {
		switch g.logic {
			case AND: {
				if g.left.val == ONE && g.right.val == ONE {
					g.output.val = ONE 
				} else {
					g.output.val = ZERO
				}
			}
			case OR: {
				if g.left.val == ONE || g.right.val == ONE {
					g.output.val = ONE 
				} else {
					g.output.val = ZERO
				}
			}
			case XOR: {
				if g.left.val != g.right.val {
					g.output.val = ONE 
				} else {
					g.output.val = ZERO
				}
			}
		}
		return true
	}
	return false
}

func createOutput(wires map[string]*wire) int {
	output := 0
	for key, wire := range wires {
		if indexStr, hasPrefix := strings.CutPrefix(key, "z"); hasPrefix {
			if wire.val == UNDEFINED {
				panic("wire is undefined")
			}
			if wire.val == ONE {
				index, err := strconv.Atoi(indexStr)
				if err != nil {
					panic("conversion failed")
				}
				output |= (1 << index)
			}
		}
	}
	return output
}


func sol1() {
	gates, wires := parseInput("input.txt")
	update(gates)
	output := createOutput(wires)
	fmt.Printf("the bits %b make the number %d\n", output, output)
}

//let's try a greedy implementation

func hammingDistance(a int, b int) int{
	nbDifferentBits := 0
	for different := a ^ b; different > 0; different >>= 1 {
		nbDifferentBits += different & 1
	}
	return nbDifferentBits
}

func wireToInt(wires map[string]*wire, wirePrefix string) (int, bool) {
	output := 0
	for key, wire := range wires {
		if indexStr, hasPrefix := strings.CutPrefix(key, wirePrefix); hasPrefix {
			if wire.val == UNDEFINED {
				return 0, true
			}
			if wire.val == ONE {
				index, err := strconv.Atoi(indexStr)
				if err != nil {
					panic("conversion failed")
				}
				output |= (1 << index)
			}
		}
	}
	return output,  false
}

func sol2() {
	gates, wires := parseInput("input.txt")
	x, _ := wireToInt(wires, "x")
	y, _ := wireToInt(wires, "y")
	expectedZ := x + y
	update(gates)
	output, _ := wireToInt(wires, "z")
	_, _, better := bestChange()
	fmt.Printf("expected: %b, broken: %b, better? %b\n", expectedZ, output, better)
	fmt.Println(hammingDistance(output, expectedZ))
	fmt.Println(hammingDistance(better, expectedZ))
}

func copyWires(a map[string]*wire, b map[string]*wire) {
	for key, _ := range a {
		*b[key] = *a[key]
	}
}

func bestChange() (*gate, *gate, int){
	var gate1, gate2 *gate
	bestDist := 64
	bestOutput := 0
	gates, wires := parseInput("input.txt")
	x, _ := wireToInt(wires, "x")
	y, _ := wireToInt(wires, "y")
	expectedZ := x + y

	initialWires := make(map[string]*wire, 0) 
	for key, val := range wires {
		cp := *val
		initialWires[key] = &cp
	}

	for i, g1 := range gates {
		for j, g2 := range gates[i:] {
			copyWires(initialWires, wires)
			tmp := g1.output
			g1.output = g2.output
			g2.output = tmp
			gates[i] = g1
			gates[j] = g2
			update(gates)
			output, undef := wireToInt(wires, "z")
			if undef {
				continue
			}
			fmt.Printf("bits: %b, hamming: %d\n", output, hammingDistance(expectedZ, output))
			dist := hammingDistance(expectedZ, output)
			if dist < bestDist {
				bestDist = dist
				gate1 = &g1
				gate2 = &g2
				bestOutput = output
			}
		}
	}
	return gate1, gate2, bestOutput
}

func main() {
	sol2()
}
