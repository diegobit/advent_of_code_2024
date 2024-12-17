package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func loadTxt(path string) ([][]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file: %w", err)
	}
	defer file.Close()

	lines := make([][]int, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Split(scanner.Text(), " ")
		if len(words) < 2 {
			fmt.Printf("Line does not contain two ints: %v\n", words)
			continue
		}
		words_i := make([]int, 0)
		for _, word := range words {
			word_i, err := strconv.Atoi(word)
			if err != nil {
				fmt.Printf("word is not an int, %s\n", word)
				os.Exit(1)
			}
			words_i = append(words_i, word_i)
		}
		lines = append(lines, words_i)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Scanner error: %w", err)
	}

	return lines, nil
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func CheckLine(line []int) bool {
	isAsc := func(l1 int, l2 int) bool {
		if l2 > l1 {
			return true
		} else {
			return false
		}
	}

	isSafe := func(l1 int, l2 int) bool {
		if Abs(l2-l1) >= 1 && Abs(l2-l1) <= 3 {
			return true
		}
		return false
	}

	prev_l := line[0]
	curr_l := line[1]
	dir := isAsc(prev_l, curr_l)
	if !isSafe(prev_l, curr_l) {
		return false
	}
	// fmt.Printf("%d, %d, %t, %t\n", prev_l, curr_l, isAsc(prev_l, curr_l), isSafe(prev_l, curr_l))

	for _, level := range line[2:] {
		prev_l = curr_l
		curr_l = level
		// fmt.Printf("%d, %d, %t, %t\n", prev_l, curr_l, isAsc(prev_l, curr_l), isSafe(prev_l, curr_l))

		if isAsc(prev_l, curr_l) != dir {
			return false
		}

		if !isSafe(prev_l, curr_l) {
			return false
		}
	}
	// fmt.Println("---")

	return true
}

func main() {
	data, err := loadTxt("./data/input_2.txt")
	if err != nil {
		fmt.Printf("Error loading file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("lines: %d, first: %+v\n", len(data), data[0])
	fmt.Println("---")

	n_safe := 0
	for _, line := range data {
		// fmt.Println(line)
		if CheckLine(line) {
			// fmt.Println("true")
			n_safe += 1
		} else {
			found := false
			// fmt.Println(line[1:])
			if CheckLine(line[1:]) {
				// fmt.Println("true")
				n_safe += 1
				found = true
			}
			if !found {
				for i := 1; i < len(line); i++ {
					lcopy := make([]int, len(line))
					copy(lcopy, line)
					newline := append(lcopy[0:i], lcopy[i+1:]...)
					// fmt.Println(newline)
					if CheckLine(newline) {
						// fmt.Println("true")
						n_safe += 1
						found = true
						break
					}
				}
			}
		}
	}
	fmt.Println(n_safe)
}
