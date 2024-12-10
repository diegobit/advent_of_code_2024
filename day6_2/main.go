package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"time"
)

// Assumptions:
// - 1 guard
// - no conflicts, eg. corners in which tht guard has to turn twice

type GameMap = [][]Cell

type Cell int

const (
	Empty Cell = iota
	Obstacle
	GuardUp
	GuardDown
	GuardLeft
	GuardRight
	NewObstacle
	VisitedVertical
	VisitedHorizontal
	VisitedCorner
)

func (c Cell) String() string {
	names := map[Cell]string{
		Empty:             ".",
		Obstacle:          "#",
		GuardUp:           "^",
		GuardDown:         "v",
		GuardLeft:         "<",
		GuardRight:        ">",
		NewObstacle:       "O",
		VisitedVertical:   "|",
		VisitedHorizontal: "—",
		VisitedCorner:     "+",
	}

	if name, ok := names[c]; ok {
		return name
	}
	return fmt.Sprintf("Unknown(%d)", c)
}

type Dir int

const (
	Up = iota
	Down
	Left
	Right
)

func (d Dir) String() string {
	names := map[Dir]string{
		Up:    "^",
		Down:  "v",
		Left:  "<",
		Right: ">",
	}

	if name, ok := names[d]; ok {
		return name
	}
	return fmt.Sprintf("Unknown(%d)", d)
}

type Guard struct {
	x   int
	y   int
	dir Dir
}

func printLevel(level *GameMap, guard *Guard, foundLoop bool) {
	for _, row := range *level {
		fmt.Println(row)
	}

	steps := 1
	for _, row := range *level {
		for _, val := range row {
			if slices.Contains([]Cell{GuardUp, GuardDown, GuardLeft, GuardRight, VisitedVertical, VisitedHorizontal, VisitedCorner}, val) {
				steps += 1
			}
		}
	}
	fmt.Printf("STATUS: Steps=%d, Loop=%v\n", steps, foundLoop)
	fmt.Println("")
}

func isOnBorder (guard Guard, level GameMap) bool {
	lastRow := len(level)-1
	lastCol := len(level[0])
	return (guard.y == 0 && guard.dir == Up) ||
	(guard.y == lastRow && guard.dir == Down) ||
	(guard.x == 0 && guard.dir == Left) ||
	(guard.x == lastCol && guard.dir == Right)
}

func dir2Cell(dir Dir) Cell {
	switch dir {
	case Up:
		return GuardUp
	case Down:
		return GuardDown
	case Left:
		return GuardLeft
	case Right:
		return GuardRight
	default:
		return Empty
	}
}

func tracePath(levelOrig *GameMap, guardOrig Guard) (*GameMap, *Guard, bool) {
	// CLONE VARS
	nRows := len(*levelOrig)
	nCols := len((*levelOrig)[0])
	level := make(GameMap, nRows)
	for i := range level {
		level[i] = make([]Cell, nCols)
		copy(level[i], (*levelOrig)[i])
	}
	guard := Guard{x: guardOrig.x, y: guardOrig.y, dir: guardOrig.dir}

	// CODE
	for {
		printLevel(&level, &guard, false)
		time.Sleep(500*time.Millisecond)
		if isOnBorder(guard, level) {
			// EXIT MAP!
			level[guardOrig.y][guardOrig.x] = dir2Cell(guardOrig.dir)
			return &level, &guard, false
		}

		switch guard.dir {
		case Up:
			if level[guard.y-1][guard.x] == VisitedVertical || level[guard.y-1][guard.x] == VisitedCorner {
				// LOOP!
				level[guardOrig.y][guardOrig.x] = dir2Cell(guardOrig.dir)
				return &level, &guard, true
			} else if (level[guard.y-1][guard.x] == Obstacle || level[guard.y-1][guard.x] == NewObstacle) {
				// TURN
				if level[guard.y][guard.x+1] == VisitedVertical || level[guard.y][guard.x+1] == VisitedCorner {
					level[guard.y][guard.x+1] = VisitedCorner
				} else {
					level[guard.y][guard.x+1] = VisitedHorizontal
				}
				level[guard.y][guard.x] = VisitedCorner
				guard.x += 1
				guard.dir = Right
			} else {
				// CONTINUE
				if level[guard.y-1][guard.x] == VisitedHorizontal {
					level[guard.y-1][guard.x] = VisitedCorner
				} else {
					level[guard.y-1][guard.x] = VisitedVertical
				}
				// level[guard.y][guard.x] = VisitedVertical
				guard.y -= 1
			}
		case Down:
			if level[guard.y+1][guard.x] == VisitedVertical || level[guard.y+1][guard.x] == VisitedCorner {
				// LOOP!
				level[guardOrig.y][guardOrig.x] = dir2Cell(guardOrig.dir)
				return &level, &guard, true
			} else if (level[guard.y+1][guard.x] == Obstacle || level[guard.y+1][guard.x] == NewObstacle) {
				// TURN
				if level[guard.y][guard.x-1] == VisitedVertical || level[guard.y][guard.x-1] == VisitedCorner {
					level[guard.y][guard.x-1] = VisitedCorner
				} else {
					level[guard.y][guard.x-1] = VisitedHorizontal
				}
				level[guard.y][guard.x] = VisitedCorner
				guard.x -= 1
				guard.dir = Left
			} else {
				// CONTINUE
				if level[guard.y+1][guard.x] == VisitedHorizontal {
					level[guard.y+1][guard.x] = VisitedCorner
				} else {
					level[guard.y+1][guard.x] = VisitedVertical
				}
				guard.y += 1
			}
		case Left:
			if level[guard.y][guard.x-1] == VisitedHorizontal || level[guard.y][guard.x-1] == VisitedCorner {
				// LOOP!
				level[guardOrig.y][guardOrig.x] = dir2Cell(guardOrig.dir)
				return &level, &guard, true
			} else if (level[guard.y][guard.x-1] == Obstacle || level[guard.y][guard.x-1] == NewObstacle) {
				// TURN
				if level[guard.y-1][guard.x] == VisitedHorizontal || level[guard.y-1][guard.x] == VisitedCorner {
					level[guard.y-1][guard.x] = VisitedCorner
				} else {
					level[guard.y-1][guard.x] = VisitedVertical
				}
				level[guard.y][guard.x] = VisitedCorner
				guard.y -= 1
				guard.dir = Up
			} else {
				// CONTINUE
				if level[guard.y][guard.x-1] == VisitedVertical {
					level[guard.y][guard.x-1] = VisitedCorner
				} else {
					level[guard.y][guard.x-1] = VisitedHorizontal
				}
				guard.x -= 1
			}
		case Right:
			if level[guard.y][guard.x+1] == VisitedHorizontal || level[guard.y][guard.x+1] == VisitedCorner {
				// LOOP!
				level[guardOrig.y][guardOrig.x] = dir2Cell(guardOrig.dir)
				return &level, &guard, true
			} else if (level[guard.y][guard.x+1] == Obstacle || level[guard.y][guard.x+1] == NewObstacle) {
				// TURN
				if level[guard.y+1][guard.x] == VisitedHorizontal || level[guard.y+1][guard.x] == VisitedCorner {
					level[guard.y+1][guard.x] = VisitedCorner
				} else {
					level[guard.y+1][guard.x] = VisitedVertical
				}
				level[guard.y][guard.x] = VisitedCorner
				guard.y += 1
				guard.dir = Down
			} else {
				// CONTINUE
				if level[guard.y][guard.x+1] == VisitedVertical {
					level[guard.y][guard.x+1] = VisitedCorner
				} else {
					level[guard.y][guard.x+1] = VisitedHorizontal
				}
				guard.x += 1
			}
		}
	}
}

func Problem2(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	level := &GameMap{}
	var guard *Guard = nil
	cursorX := 0
	cursorY := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := []Cell{}

		for _, elem := range line {
			if elem == '.' {
				row = append(row, Empty)
			} else if elem == '#' {
				row = append(row, Obstacle)
			} else if elem == 'O' {
				row = append(row, NewObstacle)
			} else if elem == '^' {
				row = append(row, GuardUp)
				guard = &Guard{x: cursorX, y: cursorY, dir: Up}
			} else if elem == 'v' {
				row = append(row, GuardDown)
				guard = &Guard{x: cursorX, y: cursorY, dir: Down}
			} else if elem == '<' {
				row = append(row, GuardLeft)
				guard = &Guard{x: cursorX, y: cursorY, dir: Left}
			} else if elem == '>' {
				row = append(row, GuardRight)
				guard = &Guard{x: cursorX, y: cursorY, dir: Right}
			} else {
				fmt.Println("error")
				os.Exit(1)
			}
			if guard == nil {
				cursorX += 1
			}
		}

		cursorX = 0
		cursorY += 1

		*level = append(*level, row)
	}

	fmt.Println("=== START ===")
	printLevel(level, guard, false)

	levelMapSolved, guardSolved, _ := tracePath(level, *guard)

	fmt.Println("=== ORIGINAL MAP SOLVED ===")
	printLevel(levelMapSolved, guardSolved, false)

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no path")
		os.Exit(1)
	}
	path := os.Args[1]
	Problem2(path)
}
