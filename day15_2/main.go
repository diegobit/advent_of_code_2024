package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Point struct {
	y int
	x int
}

type Set map[string]struct{}

func (s Set) Add(key Point) {
	k := fmt.Sprintf("y:%d,x:%d", key.y, key.x)
	s[k] = struct{}{}
}

func (s Set) Contains(key Point) bool {
	k := fmt.Sprintf("y:%d,x:%d", key.y, key.x)
	_, ok := s[k]
	return ok
}

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
	(*m)[ydst][xdst] = (*m)[ysrc][xsrc]
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
			x := 0
			for _, cell := range line {
				switch cell {
				case '.':
					mRow = append(mRow, ".")
					mRow = append(mRow, ".")
				case '#':
					mRow = append(mRow, "#")
					mRow = append(mRow, "#")
				case 'O':
					mRow = append(mRow, "[")
					mRow = append(mRow, "]")
				case '@':
					mRow = append(mRow, "@")
					mRow = append(mRow, ".")
					if cell == '@' {
						warehouse.rx = x
						warehouse.ry = y
					}
				}
				// mRow = append(mRow, string(cell))
				x += 2
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

func followOs(m *Matrix, y, x int, move string) (boxes []Point) {
	// Build list of boxes ordered st.: starting from the end, I can move them one by one
	// NOTE: Return boxes only if they CAN be moved, ie. there is space after each cell to be moved
	switch move {
	case ">", "<":
		nextY, nextX := y, x
		for (*m)[nextY][nextX] == "[" || (*m)[nextY][nextX] == "]" {
			boxes = append(boxes, Point{nextY, nextX})
			nextY, nextX = getAdjacentYX(nextY, nextX, move)
		}
		// Check if CAN be moved
		if (*m)[nextY][nextX] != "." {
			return []Point{}
		}
	case "^", "v":
		intMove := 1
		if move == "^" {
			intMove = -1
		}

		visited := Set{}
		visited.Add(Point{y, x})
		currLevel := []Point{{y, x}}
		nextLevel := []Point{}

		for len(currLevel) > 0 {
			// Extend current level
			currIdx := 0
			for currIdx < len(currLevel) { // do it until I find no new boxes
				frozenLen := len(currLevel)
				for currIdx < frozenLen {
					pt := currLevel[currIdx]
					// visited.Add(pt)
					cell := (*m)[pt.y][pt.x]
					// prevPt := Point{0, 0} // Always '#'
					// if currIdx > 0 {
					// 	prevPt = currLevel[currIdx-1]
					// }
					// prevCell := (*m)[prevPt.y][prevPt.]

					ptR := Point{pt.y, pt.x + 1}
					ptL := Point{pt.y, pt.x - 1}
					if cell == "[" && (*m)[ptR.y][ptR.x] == "]" && !visited.Contains(ptR) {
						visited.Add(ptR)
						currLevel = append(currLevel, ptR)
					} else if cell == "]" && (*m)[ptL.y][ptL.x] == "[" && !visited.Contains(ptL) {
						visited.Add(ptL)
						currLevel = append(currLevel, ptL)
					}
					currIdx++
				}
			}
			// copy currLevel into boxes
			boxes = append(boxes, currLevel...)
			// get nextLevel
			for _, pt := range currLevel {
				if (*m)[pt.y+intMove][pt.x] == "[" || (*m)[pt.y+intMove][pt.x] == "]" {
					nextLevel = append(nextLevel, Point{pt.y + intMove, pt.x})
				}
			}
			// reset for next iter
			currLevel = nextLevel
			for _, pt := range currLevel {
				visited.Add(pt)
			}
			nextLevel = []Point{}
		}
		// Check if can be moved
		for _, pt := range boxes {
			if (*m)[pt.y+intMove][pt.x] == "#" {
				return []Point{}
			}
		}
	}

	return
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
		case "[", "]":
			boxes := followOs(m, adjY, adjX, move)
			fmt.Printf("%d: [FL  ][%s] - Boxes: %v\n", i, move, boxes)
			if len(boxes) > 0 {
				for j := len(boxes) - 1; j >= 0; j-- {
					y := boxes[j].y
					x := boxes[j].x
					emptyY, emptyX := getAdjacentYX(y, x, move)
					m.move(y, x, emptyY, emptyX)
				}
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
	// maxY := len(*warehouse.m)-1
	// maxX := len((*warehouse.m)[0])-1
	sum := 0
	for y, row := range *warehouse.m {
		for x, cell := range row {
			if cell == "[" {
				sum += 100*y+x
				// minDistanceXEdge := x
				// xFromLeft := x
				// xFromRight := maxX - (x+1) // x+1 because of the "]"
				// if xFromRight < xFromLeft {
				// 	minDistanceXEdge = xFromRight
				// } else {
				// 	minDistanceXEdge = xFromLeft
				// }
				//
				// minDistanceYEdge := y
				// yFromTop := y
				// yFromBottom := maxY - y
				// if yFromBottom < yFromTop {
				// 	minDistanceYEdge = yFromBottom
				// }
				//
				// sum += 100*minDistanceYEdge + minDistanceXEdge
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
