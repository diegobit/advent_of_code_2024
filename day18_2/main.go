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

func walk(corrupted map[[2]int]bool, gridSize int) (int) {
	queue := []Point{{0, 0, 0}}
	visited := make(map[[2]int]bool)
	visited[[2]int{0, 0}] = true

	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.x == gridSize-1 && current.y == gridSize-1 {
			return current.steps
		}

		for _, dir := range directions {
			nx, ny := current.x+dir[0], current.y+dir[1]
			pos := [2]int{nx, ny}

			if nx >= 0 && nx < gridSize && ny >= 0 && ny < gridSize && !corrupted[pos] && !visited[pos] {
				visited[pos] = true
				queue = append(queue, Point{nx, ny, current.steps + 1})
			}
		}
	}

	return -1
}

func solve(path string, to_read int) (string) {
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
	next_corr := make([][2]int, 0)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		chunks := strings.Split(line, ",")
		x, _ := strconv.Atoi(chunks[0])
		y, _ := strconv.Atoi(chunks[1])
		if i < to_read {
			corrupted[[2]int{x, y}] = true
			(*memory)[y][x] = "#"
		} else {
			next_corr = append(next_corr, [2]int{x, y})
		}
		// fmt.Printf("x: %d y: %d\n%v\n", x, y, memory)
		i++
	}

	// Find byte that blocks the path
	for i := range next_corr {
		next_byte := next_corr[i]
		corrupted[next_byte] = true
		min_steps := walk(corrupted, gridSize)
		fmt.Printf("steps to reach the exit: %d, next_byte: %v\n", min_steps, next_byte)
		if min_steps == -1 {
			x := strconv.Itoa(next_byte[0])
			y := strconv.Itoa(next_byte[1])
			str_byte := strings.Join([]string{x, y}, ",")
			return str_byte
		}
	}

	return "-1,-1"

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no path")
		os.Exit(1)
	}
	path := os.Args[1]
	result := solve(path, 1024)
	fmt.Printf("Result: %s\n", result)

}
