package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Op = int8

const (
	adv Op = iota
	bxl
	bst
	jnz
	bxc
	out
	bdv
	cdv
)

var A int
var B int
var C int

func combo2int(op int8) int {
	switch op {
	case 0, 1, 2, 3:
		return int(op)
	case 4:
		return A
	case 5:
		return B
	case 6:
		return C
	case 7:
		panic("7 is reserved")
	}
	panic(fmt.Sprint("invalid instruction: %d", op))
}

func readFile(path string) (program []int8) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// registers
	scanner.Scan()
	line := scanner.Text()
	chunks := strings.Split(line, ":")
	A, _ = strconv.Atoi(strings.TrimSpace(chunks[1]))
	scanner.Scan()
	line = scanner.Text()
	chunks = strings.Split(line, ":")
	B, _ = strconv.Atoi(strings.TrimSpace(chunks[1]))
	scanner.Scan()
	line = scanner.Text()
	chunks = strings.Split(line, ":")
	C, _ = strconv.Atoi(strings.TrimSpace(chunks[1]))

	// program
	scanner.Scan() // whitespace
	scanner.Scan()
	line = scanner.Text()
	chunks = strings.Split(line, ":")

	// Convert program
	var code int
	instructions := strings.Split(strings.TrimSpace(chunks[1]), ",")
	for _, codeStr := range instructions {
		code, _ = strconv.Atoi(codeStr)
		program = append(program, int8(code))
	}

	return
}

func run(program []int8) (output string) {
	p := 0
	outStack := []string{}

	for p < len(program) {
		i := program[p]
		op := program[p+1]
		p += 2 // pre-increase it, will be overridden by jumps
		switch i {
		case adv:
			comboOp := combo2int(op)
			den := int(math.Pow(2.0, float64(comboOp)))
			A = int(A / den)
		case bxl:
			litOp := int(op)
			B = B ^ litOp
		case bst:
			B = combo2int(op) & 7 // keep last 3 bits, ie. do modulo 8
		case jnz:
			if A != 0 {
				litOp := int(op)
				p = litOp
			}
		case bxc:
			B = B ^ C
		case out:
			comboOp := combo2int(op) & 7
			outStack = append(outStack, strconv.Itoa(comboOp))
		case bdv:
			comboOp := combo2int(op)
			den := int(math.Pow(2.0, float64(comboOp)))
			B = int(A / den)
		case cdv:
			comboOp := combo2int(op)
			den := int(math.Pow(2.0, float64(comboOp)))
			C = int(A / den)
		}
	}

	output = strings.Join(outStack, ",")
	return
}

func Problem(path string) (score string) {
	data := readFile(path)
	fmt.Printf("A: %d\nB: %d\nC: %d\n", A, B, C)
	fmt.Printf("Data: %v\n", data)
	score = run(data)
	fmt.Printf("Score: %s\n", score)
	return
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no path")
		os.Exit(1)
	}
	path := os.Args[1]
	Problem(path)
}
