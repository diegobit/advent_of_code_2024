package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Set map[string]struct{}

func (s Set) Add(key string) {
	s[key] = struct{}{}
}

func (s Set) Contains(key string) bool {
	_, ok := s[key]
	return ok
}

func trail2Id(t Trail) string {
	return fmt.Sprintf("%d-%d-%d-%d", t.startY, t.startX, t.endY, t.endX)
}

func coord2Id(x int, y int) string {
	return fmt.Sprintf("%d-%d", y, x)
}

type Trail struct {
	startX int
	startY int
	endX   int
	endY   int
}

func readFile(path string) [][]int8 {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	terrain := make([][]int8, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		chars := strings.Split(line, "")
		terrainRow := make([]int8, 0)
		for _, char := range chars {
			if char == "." {
				terrainRow = append(terrainRow, int8(-1))
				continue
			}
			char_i, err := strconv.ParseInt(char, 10, 8)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			terrainRow = append(terrainRow, int8(char_i))
		}
		terrain = append(terrain, terrainRow)
	}

	return terrain
}

func scanMap(terrain [][]int8) (trailheads map[string]int, visitedTrails Set) {
	trailheads = make(map[string]int)
	visitedTrails = make(Set)
	for y, row := range terrain {
		for x, altitude := range row {
			// fmt.Printf("%d %d %d\n", x, y, char)
			if altitude == 0 {
				var trails []Trail
				trails = getTrails(terrain, x, y, x, y, -1, -1, 0)
				for _, trail := range trails {
					trailheads[coord2Id(trail.startX, trail.startY)]++
					visitedTrails.Add(trail2Id(trail))
				}
			}
		}
	}
	return
}


func getTrails(terrain [][]int8, sX int, sY int, x int, y int, prevX int, prevY int, expectedAlt int8) []Trail {
	// fmt.Println(terrain)
	// fmt.Printf("sX: %d, xY: %d, x: %d, y: %d, prevX: %d, currX: %d, alt: %d\n", sX, sY, x, y, prevX, prevY, alt)
	if expectedAlt == 9 && terrain[y][x] == 9 {
		return []Trail{{startX: sX, startY: sY, endX: x, endY: y}}
	}
	if expectedAlt > 9 {
		return []Trail{}
	}

	trails := []Trail{}
	if expectedAlt == terrain[y][x] {
		if x > 0 && prevX != x-1 && prevY != y {
			res := getTrails(terrain, sX, sY, x-1, y, prevX, prevY, expectedAlt+1)
			if res != nil {
				trails = append(trails, res...)
			}
		}
		if y > 0 && prevX != x && prevY != y-1 {
			res := getTrails(terrain, sX, sY, x, y-1, prevX, prevY, expectedAlt+1)
			if res != nil {
				trails = append(trails, res...)
			}
		}
		if x < len(terrain[0])-1 && prevX != x+1 && prevY != y {
			res := getTrails(terrain, sX, sY, x+1, y, prevX, prevY, expectedAlt+1)
			if res != nil {
				trails = append(trails, res...)
			}
		}
		if y < len(terrain)-1 && prevX != x && prevY != y+1 {
			res := getTrails(terrain, sX, sY, x, y+1, prevX, prevY, expectedAlt+1)
			if res != nil {
				trails = append(trails, res...)
			}
		}
	}

	return trails
}

func getSumScores(trailheads map[string]int) int {
	sum := 0
	for _, num := range trailheads {
		sum += num
	}
	return sum
}

func Problem(path string) {
	terrain := readFile(path)
	trailheads, visitedTrails := scanMap(terrain)
	fmt.Println(trailheads)
	fmt.Println(visitedTrails)
	fmt.Printf("Answer: %d\n", getSumScores(trailheads))

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no path")
		os.Exit(1)
	}
	path := os.Args[1]
	Problem(path)
}
