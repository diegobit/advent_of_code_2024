package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Matrix [][]string

func (m *Matrix) String() string {
	s := strings.Builder{}
	for _, row := range *m {
		for _, cell := range row {
			s.WriteString(fmt.Sprintf("%s ", cell))
		}
		s.WriteString("\n")
	}
	return s.String()
}

type Point struct {
	x     int
	y     int
	steps int
}


func solve(path string, to_read int) (int, *Matrix) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// Find gridSize
	scanner := bufio.NewScanner(file)
	gridSize := 0
	for scanner.Scan() {
		line := scanner.Text()
		chunks := strings.Split(line, ",")
		x, _ := strconv.Atoi(chunks[0])
		y, _ := strconv.Atoi(chunks[1])
		x++
		y++
		if x > gridSize {
			gridSize = x
		}
		if y > gridSize {
			gridSize = y
		}
	}

	// Allocate
	memory := &Matrix{}
	for range gridSize {
		row := make([]string, gridSize)
		for j := range gridSize {
			row[j] = "."
		}
		*memory = append(*memory, row)
	}
	// fmt.Printf("Empty memory:\n%v\n", memory)

	// Reset scanner
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
	scanner = bufio.NewScanner(file)

	// Read corrupted byte positions
	corrupted := make(map[[2]int]bool)
	i := 0
	for scanner.Scan() {
		if i >= to_read {
			break
		}
		line := scanner.Text()
		chunks := strings.Split(line, ",")
		x, _ := strconv.Atoi(chunks[0])
		y, _ := strconv.Atoi(chunks[1])
		corrupted[[2]int{x, y}] = true
		(*memory)[y][x] = "#"
		// fmt.Printf("x: %d y: %d\n%v\n", x, y, memory)
		i++
	}

	// Find shortest path
	// start := Point{0, 0, 0}
	// exit := Point{6, 6, int(^uint32(0) >> 1)}

	queue := []Point{{0, 0, 0}}
	visited := make(map[[2]int]bool)
	visited[[2]int{0, 0}] = true

	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.x == gridSize-1 && current.y == gridSize-1 {
			return current.steps, memory
		}

		for _, dir := range directions {
			nx, ny := current.x+dir[0], current.y+dir[1]
			pos := [2]int{nx, ny}

			if nx >= 0 && nx < gridSize && ny >= 0 && ny < gridSize && !corrupted[pos] && !visited[pos] {
				visited[pos] = true
				(*memory)[ny][nx] = "O"
				queue = append(queue, Point{nx, ny, current.steps + 1})
			}
		}
	}

	return -1, memory
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no path")
		os.Exit(1)
	}
	path := os.Args[1]
	result, memory := solve(path, 1024)
	fmt.Printf("Memory:\n%v\n", memory)
	fmt.Printf("Result: %d\n", result)

}
