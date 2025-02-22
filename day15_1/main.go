package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// MATRIX
type Matrix [][]string

func (m *Matrix) String() string {
	s := strings.Builder{}
	for _, row := range *m {
		s.WriteString(fmt.Sprintf("%v\n", row))
	}
	return s.String()
}

func (m *Matrix) move(ysrc, xsrc, ydst, xdst int) {
	(*m)[ydst][xdst] = "@"
	(*m)[ysrc][xsrc] = "."
}

// WAREHOUSE
type Warehouse struct {
	m     *Matrix
	rx    int
	ry    int
	moves *[]string
}

func (w *Warehouse) String() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("Warehouse(rx=%d, ry=%d,\nm=\n%v\nmoves=%v)\n", w.rx, w.ry, w.m, w.moves))
	return s.String()
}

func (w *Warehouse) move(ydst, xdst int) {
	w.m.move(w.ry, w.rx, ydst, xdst)
	w.ry = ydst
	w.rx = xdst
}

func readFile(path string) (warehouse *Warehouse) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	warehouse = &Warehouse{}
	warehouse.m = &Matrix{}
	warehouse.moves = &[]string{}

	scanner := bufio.NewScanner(file)
	y := 0
	mode := "matrix" // matrix or move
	for scanner.Scan() {
		line := scanner.Text()
		if line == "\n" || line == "" {
			mode = "move"
			continue
		}

		switch mode {
		case "matrix":
			mRow := make([]string, 0)
			for x, cell := range line {
				mRow = append(mRow, string(cell))
				if cell == '@' {
					warehouse.rx = x
					warehouse.ry = y
				}
			}
			*warehouse.m = append(*warehouse.m, mRow)
			y++
		case "move":
			for _, insRune := range line {
				*warehouse.moves = append(*warehouse.moves, string(insRune))
			}
		}
	}

	return
}

func getAdjacentYX(y int, x int, move string) (newy int, newx int) {
	switch move {
	case "^":
		newy = y - 1
		newx = x
	case ">":
		newy = y
		newx = x + 1
	case "v":
		newy = y + 1
		newx = x
	case "<":
		newy = y
		newx = x - 1
	}
	return
}

func followOs(m *Matrix, y, x int, move string) (finalY, finalX int, finalCell string) {
	nextY, nextX := y, x
	for (*m)[nextY][nextX] == "O" {
		nextY, nextX = getAdjacentYX(nextY, nextX, move)
	}

	return nextY, nextX, (*m)[nextY][nextX]
}

func simulate(warehouse *Warehouse) {
	m := warehouse.m
	var adjX, adjY int
	for i, move := range *warehouse.moves {
		adjY, adjX = getAdjacentYX(warehouse.ry, warehouse.rx, move)
		switch (*m)[adjY][adjX] {
		case ".":
			fmt.Printf("%d: [MOVE][%s] - adjY: %d, adjX: %d\n", i, move, adjY, adjX)
			warehouse.move(adjY, adjX)
			// fmt.Printf("x,r: %d, %d\n", warehouse.ry, warehouse.rx)
			// fmt.Println(warehouse.m)
		case "O":
			finalY, finalX, finalCell := followOs(m, adjY, adjX, move)
			fmt.Printf("%d: [FL %s][%s] - finY: %d, finX: %d,\n", i, finalCell, move, finalY, finalX)
			if finalCell == "." {
				// Move first O to final place of .
				(*m)[finalY][finalX] = "O"
				warehouse.move(adjY, adjX)
			}
			// fmt.Printf("x,r: %d, %d\n", warehouse.ry, warehouse.rx)
			// fmt.Println(warehouse.m)
		default:
			fmt.Printf("%d: [ST %s][%s] - ry: %d, rx: %d\n", i, (*m)[adjY][adjX], move, warehouse.ry, warehouse.rx)
		}
	}
}

func sumBoxesGPS(warehouse *Warehouse) int {
	sum := 0
	for y, row := range *warehouse.m {
		for x, cell := range row {
			if cell == "O" {
				sum += 100*y + x
			}
		}
	}
	return sum
}

func Problem(path string) {
	warehouse := readFile(path)
	fmt.Println(warehouse)
	fmt.Println()

	simulate(warehouse)
	fmt.Println(warehouse)
	fmt.Println()

	fmt.Printf("TOTAL: %d\n", sumBoxesGPS(warehouse))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no path")
		os.Exit(1)
	}
	path := os.Args[1]
	Problem(path)
}
