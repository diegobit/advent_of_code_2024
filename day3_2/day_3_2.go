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
	Initial							State = "init"
	InitialDisabled					State = "init_disabled"
	D_WaitO							State = "d_o"
	D_WaitOpenOrN					State = "d_(n"
	D_WaitClose						State = "d_)"
	Dn_WaitQuote					State = "dn_'"
	Dn_WaitT						State = "dn_t"
	Dn_WaitOpen						State = "dn_("
	Dn_WaitClose					State = "dn_)"
	M_WaitU							State = "m_u"
	M_WaitL							State = "m_l"
	WaitOpenBracket					State = "("
	WaitLeftDigit					State = "digit_left"
	WaitLeftDigitOrComma			State = "digit_,_left"
	WaitLeftDigitOrCommaFinal		State = "digit_,_left_final"
	WaitComma						State = "c"
	WaitRightDigit					State = "digit_right"
	WaitRightDigitOrBracket			State = "digit_)_right"
	WaitRightDigitOrBracketFinal	State = "digit_)_right_final"
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
	case Initial:
		if c == 'd' {
			return D_WaitO, nil, OpNo, 0, 0
		} else if c == 'm' {
			return M_WaitU, nil, OpNo, 0, 0
		}
		return Initial, nil, OpNo, 0, 0
	case InitialDisabled:
		if c == 'd' {
			return D_WaitO, nil, OpNo, 0, 0
		}
		return InitialDisabled, nil, OpNo, 0, 0
	case D_WaitO:
		if c == 'o' {
			return D_WaitOpenOrN, nil, OpNo, 0, 0
		} else if c == 'm' {
			return M_WaitU, nil, OpNo, 0, 0
		}
	case D_WaitOpenOrN:
		if c == '(' {
			return D_WaitClose, nil, OpNo, 0, 0
		} else if c == 'n' {
			return Dn_WaitQuote, nil, OpNo, 0, 0
		} else if c == 'm' {
			return M_WaitU, nil, OpNo, 0, 0
		}
	case D_WaitClose:
		if c == ')' {
			return Initial, nil, OpNo, 0, 0
		} else if c == 'm' {
			return M_WaitU, nil, OpNo, 0, 0
		}
	case Dn_WaitQuote:
		if c == '\'' {
			return Dn_WaitT, nil, OpNo, 0, 0
		} else if c == 'm' {
			return M_WaitU, nil, OpNo, 0, 0
		}
	case Dn_WaitT:
		if c == 't' {
			return Dn_WaitOpen, nil, OpNo, 0, 0
		} else if c == 'm' {
			return M_WaitU, nil, OpNo, 0, 0
		}
	case Dn_WaitOpen:
		if c == '(' {
			return Dn_WaitClose, nil, OpNo, 0, 0
		} else if c == 'm' {
			return M_WaitU, nil, OpNo, 0, 0
		}
	case Dn_WaitClose:
		if c == ')' {
			return InitialDisabled, nil, OpNo, 0, 0
		} else if c == 'm' {
			return M_WaitU, nil, OpNo, 0, 0
		}
	case M_WaitU:
		if c == 'u' {
			return M_WaitL, nil, OpNo, 0, 0
		}
	case M_WaitL:
		if c == 'l' {
			return WaitOpenBracket, nil, OpNo, 0, 0
		}
	case WaitOpenBracket:
		if c == '(' {
			return WaitLeftDigit, nil, OpMul, 0, 0
		}
	case WaitLeftDigit:
		if isDigit(c) {
			return WaitLeftDigitOrComma, []rune{c}, OpMul, 0, 0
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
	return Initial, nil, OpNo, 0, 0
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
	state := Initial
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
			state = Initial
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
	fmt.Println(total)
}

func main() {
	MyMain()
}

