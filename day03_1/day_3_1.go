package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

type State string
const (
	InitialWaitM					State = "m"
	WaitU							State = "u"
	WaitL							State = "l"
	WaitOpenBracket					State = "("
	WaitLeftDigit					State = "l_d"
	WaitLeftDigitOrComma			State = "l_dc"
	WaitLeftDigitOrCommaFinal		State = "l_dc_f"
	WaitComma						State = "c"
	WaitRightDigit					State = "r_d"
	WaitRightDigitOrBracket			State = "r_d)"
	WaitRightDigitOrBracketFinal	State = "r_d)_f"
	WaitCloseBracket				State = ")"
	Final							State = "final"
)

type Operation string
const (
	OpMul Operation = "mul"
	OpNo Operation = "no"
)

func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

func parseState(c rune, state State, stack []rune, lh int) (State, []rune, Operation, int, int) {
	// c := string(char)

	switch state {
	case InitialWaitM:
		if c == 'm' {
			return WaitU, nil, OpNo, lh, 0
		}
	case WaitU:
		if c == 'u' {
			return WaitL, nil, OpNo, lh, 0
		}
	case WaitL:
		if c == 'l' {
			return WaitOpenBracket, nil, OpNo, lh, 0
		}
	case WaitOpenBracket:
		if c == '(' {
			return WaitLeftDigit, nil, OpMul, lh, 0
		}
	case WaitLeftDigit:
		if isDigit(c) {
			return WaitLeftDigitOrComma, []rune{c}, OpMul, lh, 0
		}
	case WaitLeftDigitOrComma:
		if isDigit(c) {
			return WaitLeftDigitOrCommaFinal, append(stack, c), OpMul, lh, 0
		} else if c == ',' {
			newLh, _ := strconv.Atoi(string(stack))
			return WaitRightDigit, nil, OpMul, newLh, 0
		}
	case WaitLeftDigitOrCommaFinal:
		if isDigit(c) {
			return WaitComma, append(stack, c), OpMul, lh, 0
		} else if c == ',' {
			newLh, _ := strconv.Atoi(string(stack))
			return WaitRightDigit, nil, OpMul, newLh, 0
		}
	case WaitComma:
		if c == ',' {
			newLh, _ := strconv.Atoi(string(stack))
			return WaitRightDigit, nil, OpMul, newLh, 0
		}
	case WaitRightDigit:
		if isDigit(c) {
			return WaitRightDigitOrBracket, []rune{c}, OpMul, lh, 0
		}
	case WaitRightDigitOrBracket:
		if isDigit(c) {
			return WaitRightDigitOrBracketFinal, append(stack, c), OpMul, lh, 0
		} else if c == ')' {
			newRh, _ := strconv.Atoi(string(stack))
			return Final, nil, OpMul, lh, newRh
		}
	case WaitRightDigitOrBracketFinal:
		if isDigit(c) {
			return WaitCloseBracket, append(stack, c), OpMul, lh, 0
		} else if c == ')' {
			newRh, _ := strconv.Atoi(string(stack))
			return Final, nil, OpMul, lh, newRh
		}
	case WaitCloseBracket:
		if c == ')' {
			newRh, _ := strconv.Atoi(string(stack))
			return Final, nil, OpMul, lh, newRh
		}
	case Final:
		fmt.Println("Final?")
	}
	return InitialWaitM, nil, OpNo, 0, 0
}

func executeOperation(op Operation, lh int, rh int) int {
	switch op {
	case OpMul:
		return lh*rh
	}
	return 0
}

func MyMain() {
	path := "../data/input_3.txt"
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	total := 0
	state := InitialWaitM
	stack := []rune{}
	lh := 0

	reader := bufio.NewReader(file)
	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		// fmt.Printf("%c", char)
		newState, newStack, op, newLh, rh := parseState(char, state, stack, lh)
		// fmt.Printf("%s, %v, %s, %v, %d, %d\n", string(char), newState, newStack, op, newLh, rh)

		if newState == Final {
			// fmt.Println("final")
			total += executeOperation(op, newLh, rh)
			state = InitialWaitM
			stack = nil
			lh = 0
		} else {
			// fmt.Println("else")
			state = newState
			stack = newStack
			lh = newLh
		}
		// fmt.Println("---")
	}
	// fmt.Println(total)
}

func main() {
	MyMain()
}
