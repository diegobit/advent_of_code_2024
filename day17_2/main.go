package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	// "slices"
	"strconv"
	"strings"
)

type Set map[int64]struct{}

func (s Set) Add(key int64) {
	s[key] = struct{}{}
}

func (s Set) Contains(key int64) bool {
	_, ok := s[key]
	return ok
}

func (s Set) Min() int64 {
	m := int64(math.MaxInt64)
	for v := range s {
		if v < m {
			m = v
		}
	}
	return m
}

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

func combo2int(op int8, A, B, C int) int {
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

func string2binary(program string) string {
	chunks := strings.Split(program, ",")
	res := []string{}
	for _, el := range chunks {
		i, _ := strconv.Atoi(el)
		res = append(res, fmt.Sprintf("%b", i))
	}
	return strings.Join(res, ",")
}

// --- combo-operand helper --------------------------------------------------
func combo(op int8, A, B, C int64) int64 {
	switch op {
	case 0, 1, 2, 3:
		return int64(op)          // literal 0‥3
	case 4:
		return A                  // value of A
	case 5:
		return B                  // value of B
	case 6:
		return C                  // value of C
	default:
		panic("operand 7 is reserved")
	}
}

func run_program(A, B, C int, program []int8) (output string) {
	p := 0
	outStack := []string{}

	for p < len(program) {
		i := program[p]
		op := program[p+1]
		p += 2 // pre-increase it, will be overridden by jumps
		switch i {
		case adv:
			comboOp := combo2int(op, A, B, C)
			den := int(math.Pow(2.0, float64(comboOp)))
			A = int(A / den)
		case bxl:
			litOp := int(op)
			B = B ^ litOp
		case bst:
			B = combo2int(op, A, B, C) & 7 // keep last 3 bits, ie. do modulo 8
		case jnz:
			if A != 0 {
				litOp := int(op)
				p = litOp
			}
		case bxc:
			B = B ^ C
		case out:
			comboOp := combo2int(op, A, B, C) & 7
			outStack = append(outStack, strconv.Itoa(comboOp))
		case bdv:
			comboOp := combo2int(op, A, B, C)
			den := int(math.Pow(2.0, float64(comboOp)))
			B = int(A / den)
		case cdv:
			comboOp := combo2int(op, A, B, C)
			den := int(math.Pow(2.0, float64(comboOp)))
			C = int(A / den)
		}
	}

	output = strings.Join(outStack, ",")
	return
}

func firstOutput(program []int8, A0 int64) int8 {
	A, B, C := A0, int64(0), int64(0)
	ip := 0
	for ip < len(program) {
		opcode, operand := program[ip], program[ip+1]
		ip += 2
		switch opcode {
		case adv:
			shift := combo(operand, A, B, C)
			if shift >= 63 {
				A = 0
			} else {
				A >>= shift
			}
		case bxl:
			B ^= int64(operand)
		case bst:
			B = combo(operand, A, B, C) & 7
		case jnz:
			if A != 0 {
				ip = int(operand) // jump target is *token* index
			}
		case bxc:
			B ^= C
		case out:
			return int8(combo(operand, A, B, C) & 7)
		case bdv:
			shift := combo(operand, A, B, C)
			if shift >= 63 {
				B = 0
			} else {
				B = A >> shift
			}
		case cdv:
			shift := combo(operand, A, B, C)
			if shift >= 63 {
				C = 0
			} else {
				C = A >> shift
			}
		default:
			panic("bad opcode")
		}
	}
	panic("program produced no output")
}

func run(program []int8) string {

	candidates := make(Set)
	candidates.Add(0)

	for i := len(program) - 1; i >= 0; i-- {
		digit := program[i]
		next_cands := make(Set)
		for c := range candidates {
			base := int64(c) << 3 // left shift of three (new three bits empty)
			for i := 0; i < 8; i++ {
				guess := base | int64(i) // add i as new least significant digit
				if firstOutput(program, guess) == digit {
					next_cands.Add(guess)
				}
			}
		candidates = next_cands
		}
	}

	solution := candidates.Min()

	execution := run_program(int(solution), 0, 0, program)
	fmt.Printf("found A: %d, program executed with this A: %v\n", solution, execution)

	return strconv.FormatInt(solution, 10)
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
