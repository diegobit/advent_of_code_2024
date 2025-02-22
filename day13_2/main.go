package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*********
** FUNCS
**********/
type Machine struct {
	Ax     int
	Ay     int
	Bx     int
	By     int
	prizeX int
	prizeY int
}

func (m Machine) String() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("(Ax=%v, Ay=%v, ", m.Ax, m.Ay))
	s.WriteString(fmt.Sprintf("Bx=%v, By=%v, ", m.Bx, m.By))
	s.WriteString(fmt.Sprintf("Px=%v, Py=%v),", m.prizeX, m.prizeY))
	return s.String()
}

func readFile(path string) (machines []*Machine) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	reBtn := regexp.MustCompile(`.*X\+(\d+), Y\+(\d+)`)
	rePrz := regexp.MustCompile(`.*X\=(\d+), Y\=(\d+)`)

	getXYvalues := func(line string, re *regexp.Regexp) (x int, y int) {
		var err1, err2 error
		match := re.FindStringSubmatch(line)
		x, err1 = strconv.Atoi(match[1])
		y, err2 = strconv.Atoi(match[2])
		if err1 != nil || err2 != nil {
			fmt.Printf("ERROR: %v, %v", err1, err2)
		}
		return
	}

	scanner := bufio.NewScanner(file)
	currMachine := &Machine{}
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "Button A:") {
			Ax, Ay := getXYvalues(line, reBtn)
			currMachine.Ax, currMachine.Ay = Ax, Ay
		} else if strings.HasPrefix(line, "Button B:") {
			Bx, By := getXYvalues(line, reBtn)
			currMachine.Bx, currMachine.By = Bx, By
		} else if strings.HasPrefix(line, "Prize:") {
			Px, Py := getXYvalues(line, rePrz)
			currMachine.prizeX, currMachine.prizeY = 10000000000000+Px, 10000000000000+Py
			machines = append(machines, currMachine)
			currMachine = &Machine{}
		} else {
			continue
		}
	}

	return
}

func getSolutions(machines []*Machine) (minTokens []int) {
	for i, m := range machines {
		found, x, y := SolveLinearEquation(m)
		if found {
			fmt.Printf("%d (Ax=%d): Solved! x=%d, y=%d\n", i+1, m.Ax, x, y)
			tokens := 3*x + y
			minTokens = append(minTokens, tokens)
		} else {
			fmt.Printf("%d (Ax=%d): No solution\n", i+1, m.Ax)
		}
	}
	return
}

func getTotal(minTokens []int) (total int) {
	for _, toks := range minTokens {
		total += toks
	}
	return
}

func SolveLinearEquation(m *Machine) (bool, int, int) {
	// Using Cramer's Rule to solve the system of equations
	// x = Dx/D, where Dx = determinant of numerator of solution of x
	// y = Dx/D, where Dy = determinant of numerator of solution of y
	//				   D  = determinant of the coefficient matrix
	// https://math.libretexts.org/Bookshelves/Precalculus/Precalculus_1e_(OpenStax)/09%3A_Systems_of_Equations_and_Inequalities/9.08%3A_Solving_Systems_with_Cramer's_Rule
	// D_x = c1 * b2 - c2 * b1
	// D_y = a1 * c2 - a2 * c1
	// D = a1 * b2 - a2 * b1

	Dx := (m.prizeX * m.By) - (m.prizeY * m.Bx)
	Dy := (m.prizeY * m.Ax) - (m.prizeX * m.Ay)
	D := (m.Ax * m.By) - (m.Ay * m.Bx)
	x := Dx / D
	y := Dy / D
	if Dx%D == 0 && Dy%D == 0 {
		return true, x, y
	}
	return false, x, y
}

func Problem(path string) {
	machines := readFile(path)
	fmt.Println(machines)
	fmt.Println()

	minTokens := getSolutions(machines)
	fmt.Printf("\nminTokens list: %v\n\n", minTokens)

	fmt.Printf("TOTAL: %v\n", getTotal(minTokens))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no path")
		os.Exit(1)
	}
	path := os.Args[1]
	Problem(path)
}
