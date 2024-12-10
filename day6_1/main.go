package main

import (
	"bufio"
	"fmt"
	"os"
)

// Assumptions:
// - 1 guard
// - no conflicts, eg. corners in which tht guard has to turn twice

type GameMap = [][]Cell
type VisitedMap = [][]bool

type Cell = int

const (
	Empty Cell = iota
	Obstacle
	GuardUp
	GuardDown
	GuardLeft
	GuardRight
)

type Dir = int

const (
	None = iota
	Up
	Down
	Left
	Right
)

type Guard struct {
	x   int
	y   int
	dir Dir
}

func printLevel(level GameMap, guard Guard, visitedMap VisitedMap) {
	fmt.Printf("Empty: %v, Obstacle: %v, GuardUp: %v, GuardDown: %v, GuardLeft: %v, GuardRight: %v\n", Empty, Obstacle, GuardUp, GuardDown, GuardLeft, GuardRight)
	for _, row := range level {
		fmt.Println(row)
	}
	fmt.Printf("Guard: x=%d, y=%d, dir=%d\n", guard.x, guard.y, guard.dir)
	steps := 0
	for _, row := range visitedMap {
		for _, val := range row {
			if val {
				steps += 1
			}
		}
	}
	fmt.Printf("visitedMap: %d\n", steps)
}

func tracePath(level GameMap, visitedMap VisitedMap, guard Guard) (GameMap, Guard, VisitedMap) {
	isOnBorder := func(guard Guard, level GameMap) bool {
		return (guard.y == 0 && guard.dir == Up) ||
			(guard.y == len(level[0])-1 && guard.dir == Down) ||
			(guard.x == 0 && guard.dir == Left) ||
			(guard.x == len(level)-1 && guard.dir == Right)
	}

	for {
		if isOnBorder(guard, level) {
			return level, guard, visitedMap
		}

		switch guard.dir {
		case Up:
			if level[guard.y-1][guard.x] == Empty {
				level[guard.y-1][guard.x] = GuardUp
				level[guard.y][guard.x] = Empty
				guard.y -= 1
			} else if level[guard.y-1][guard.x] == Obstacle {
				level[guard.y][guard.x+1] = GuardRight
				level[guard.y][guard.x] = Empty
				guard.x += 1
				guard.dir = Right
			}
			visitedMap[guard.y][guard.x] = true
		case Down:
			if level[guard.y+1][guard.x] == Empty {
				level[guard.y+1][guard.x] = GuardDown
				level[guard.y][guard.x] = Empty
				guard.y += 1
			} else if level[guard.y+1][guard.x] == Obstacle {
				level[guard.y][guard.x-1] = GuardLeft
				level[guard.y][guard.x] = Empty
				guard.x -= 1
				guard.dir = Left
			}
			visitedMap[guard.y][guard.x] = true
		case Left:
			if level[guard.y][guard.x-1] == Empty {
				level[guard.y][guard.x-1] = GuardLeft
				level[guard.y][guard.x] = Empty
				guard.x -= 1
			} else if level[guard.y][guard.x-1] == Obstacle {
				level[guard.y-1][guard.x] = GuardUp
				level[guard.y][guard.x] = Empty
				guard.y -= 1
				guard.dir = Up
			}
			visitedMap[guard.y][guard.x] = true
		case Right:
			if level[guard.y][guard.x+1] == Empty {
				level[guard.y][guard.x+1] = GuardRight
				level[guard.y][guard.x] = Empty
				guard.x += 1
			} else if level[guard.y][guard.x+1] == Obstacle {
				level[guard.y+1][guard.x] = GuardDown
				level[guard.y][guard.x] = Empty
				guard.y += 1
				guard.dir = Down
			}
			visitedMap[guard.y][guard.x] = true
		}
	}
}

func PredictPath() {
	path := "input.txt"
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	levelMap := make(GameMap, Empty)
	visitedMap := make(VisitedMap, 0)
	guard := Guard{x: 0, y: 0, dir: None}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := []Cell{}

		for _, elem := range line {
			visitedMap = append(visitedMap, make([]bool, len(line)))
			if elem == '.' {
				row = append(row, Empty)
				if guard.dir == None {
					guard.x += 1
				}
			} else if elem == '#' {
				row = append(row, Obstacle)
				if guard.dir == None {
					guard.x += 1
				}
			} else if elem == '^' {
				row = append(row, GuardUp)
				guard.dir = Up
			} else if elem == 'v' {
				row = append(row, GuardDown)
				guard.dir = Down
			} else if elem == '<' {
				row = append(row, GuardLeft)
				guard.dir = Left
			} else if elem == '>' {
				row = append(row, GuardRight)
				guard.dir = Right
			} else {
				fmt.Println("error")
				os.Exit(1)
			}
		}

		if guard.dir == None {
			guard.x = 0
			guard.y += 1
		} else {
			visitedMap[guard.y][guard.x] = true
		}

		levelMap = append(levelMap, row)
	}

	printLevel(levelMap, guard, visitedMap)

	levelMap, guard, visitedMap = tracePath(levelMap, visitedMap, guard)

	printLevel(levelMap, guard, visitedMap)
	// fmt.Println(visitedMap)
}

func main() {
	PredictPath()
}
