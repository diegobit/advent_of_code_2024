package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Assumptions:
// - only antinodes outsides
// - only one pair of antinodes
// - count antinodes even if under antenna

type Coord struct {
	x int
	y int
}

type Pair [2]Coord

type CoordId string

func (c Coord) c2id() CoordId {
	var cid CoordId = CoordId(fmt.Sprintf("%d,%d", c.y, c.x))
	return cid
}

func (c CoordId) id2c() Coord {
	coordStrings := strings.Split(string(c), ",")
	y, _ := strconv.Atoi(coordStrings[0])
	x, _ := strconv.Atoi(coordStrings[1])
	return Coord{y: y, x: x}
}

func getPairs(positions []Coord) []Pair {
	permutations := []Pair{}
	for _, c1 := range positions {
		for _, c2 := range positions[1:] {
			if c1 == c2 {
				continue
			}
			pair := Pair{c1, c2}
			permutations = append(permutations, pair)
		}
	}
	return permutations
}

func printAntinodeMap(antinodeMap map[CoordId]bool, maxX int, maxY int) {
	m := [][]string{}
	for y := 0; y < maxY; y++ {
		row := []string{}
		for x := 0; x < maxX; x++ {
			cid := Coord{x: x, y: y}.c2id()
			_, ok := antinodeMap[cid]
			if ok {
				row = append(row, "#")
			} else {
				row = append(row, ".")
			}
		}
		m = append(m, row)
	}
	for _, row := range m {
		fmt.Println(row)
	}
}

func updateAntinodeMap(antinodeMap map[CoordId]bool, an Coord) {
	anId := an.c2id()
	_, ok := antinodeMap[anId]
	if !ok {
		antinodeMap[anId] = true
	}
}

func countAntinodes(antinodeMap map[CoordId]bool) int {
	return len(antinodeMap)
}

func getAntinodePos(p Pair, maxX int, maxY int) []Coord {
	isInside := func(c Coord) bool {
		return c.x >= 0 && c.y >= 0 && c.x < maxX && c.y < maxY
	}

	diffX := int(math.Abs(float64(p[0].x - p[1].x)))
	diffY := int(math.Abs(float64(p[0].y - p[1].y)))

	antinodes := []Coord{}
	antinodes = append(antinodes, p[0])
	antinodes = append(antinodes, p[1])

	if p[0].x < p[1].x && p[0].y < p[1].y {
		newX := p[0].x - diffX
		newY := p[0].y - diffY
		for isInside(Coord{x: newX, y: newY}) {
			antinodes = append(antinodes, Coord{x: newX, y: newY})
			newX = newX - diffX
			newY = newY - diffY
		}
		newX = p[1].x + diffX
		newY = p[1].y + diffY
		for isInside(Coord{x: newX, y: newY}) {
			antinodes = append(antinodes, Coord{x: newX, y: newY})
			newX = newX + diffX
			newY = newY + diffY
		}
	} else if p[0].x < p[1].x && p[0].y > p[1].y {
		newX := p[0].x - diffX
		newY := p[0].y + diffY
		for isInside(Coord{x: newX, y: newY}) {
			antinodes = append(antinodes, Coord{x: newX, y: newY})
			newX = newX - diffX
			newY = newY + diffY
		}
		newX = p[1].x + diffX
		newY = p[1].y - diffY
		for isInside(Coord{x: newX, y: newY}) {
			antinodes = append(antinodes, Coord{x: newX, y: newY})
			newX = newX + diffX
			newY = newY - diffY
		}
	} else if p[0].x > p[1].x && p[0].y > p[1].y {
		newX := p[0].x + diffX
		newY := p[0].y + diffY
		for isInside(Coord{x: newX, y: newY}) {
			antinodes = append(antinodes, Coord{x: newX, y: newY})
			newX = newX + diffX
			newY = newY + diffY
		}
		newX = p[1].x - diffX
		newY = p[1].y - diffY
		for isInside(Coord{x: newX, y: newY}) {
			antinodes = append(antinodes, Coord{x: newX, y: newY})
			newX = newX - diffX
			newY = newY - diffY
		}
	} else { //if c[0].x > c[1].x && c[0].y < c[1].y {
		newX := p[0].x + diffX
		newY := p[0].y - diffY
		for isInside(Coord{x: newX, y: newY}) {
			antinodes = append(antinodes, Coord{x: newX, y: newY})
			newX = newX + diffX
			newY = newY - diffY
		}
		newX = p[1].x - diffX
		newY = p[1].y + diffY
		for isInside(Coord{x: newX, y: newY}) {
			antinodes = append(antinodes, Coord{x: newX, y: newY})
			newX = newX - diffX
			newY = newY + diffY
		}
	}

	return antinodes
}

func Problem2(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	antennas := map[rune][]Coord{}
	antinodeMap := map[CoordId]bool{}
	cursorY := 0
	maxX := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		maxX = len(line)
		for cursorX, elem := range line {
			if elem != '.' {
				_, ok := antennas[elem]
				if !ok {
					antennas[elem] = []Coord{}
				}
				antennas[elem] = append(antennas[elem], Coord{x: cursorX, y: cursorY})
			}
		}
		cursorY++
	}

	for _, positions := range antennas {
		fmt.Printf("Positions: %+v\n", positions)
		permutations := getPairs(positions)
		for _, pair := range permutations {
			fmt.Printf("Pair: %+v\n", pair)
			antinodes := getAntinodePos(pair, maxX, cursorY)
			for _, an := range antinodes {
				fmt.Printf("Antinode: %+v\n", an)
				updateAntinodeMap(antinodeMap, an)
			}
		}
	}

	nAntinodes := countAntinodes(antinodeMap)

	printAntinodeMap(antinodeMap, maxX, cursorY)
	fmt.Printf("Antinodes: %d\n", nAntinodes)
}


func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no path")
		os.Exit(1)
	}
	path := os.Args[1]
	Problem2(path)
}
