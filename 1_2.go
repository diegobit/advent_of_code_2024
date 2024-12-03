package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Lists struct {
	left  []int
	right []int
}

func loadTxt(path string) (*Lists, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file: %w", err)
	}
	defer file.Close()

	lists := &Lists{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Split(scanner.Text(), "   ")
		if len(words) < 2 {
			fmt.Printf("Line does not contain two ints: %v\n", words)
			continue
		}
		left, err := strconv.Atoi(words[0])
		if err != nil {
			fmt.Printf("First word is not an int, %s\n", words[0])
			continue
		}
		right, err := strconv.Atoi(words[1])
		if err != nil {
			fmt.Printf("Second word is not an int, %s\n", words[1])
			continue
		}
		// lists = append(lists, Lists{left, right})
		lists.left = append(lists.left, left)
		lists.right = append(lists.right, right)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Scanner error: %w", err)
	}

	if len(lists.left) == 0 || len(lists.right) == 0 {
		return nil, fmt.Errorf("File empty?")
	}

	if len(lists.left) != len(lists.right) {
		return nil, fmt.Errorf("Lists have different length! %d, %d", len(lists.left), len(lists.right))
	}

	return lists, nil
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	data, err := loadTxt("./data/input_1.txt")
	if err != nil {
		fmt.Printf("Error loading file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", data.left[0:3])
	fmt.Printf("%+v\n", data.right[0:3])
	fmt.Println("---")

	sort.Ints(data.left)
	sort.Ints(data.right)
	fmt.Printf("%+v\n", data.left[0:3])
	fmt.Printf("%+v\n", data.right[0:3])
	fmt.Println("---")

	sim := 0
	for i := 0; i < len(data.left); i++ {
		for j := 0; j < len(data.right); j++ {
			if data.left[i] == data.right[j] {
				sim += data.left[i]
			}
		}
	}

	fmt.Println(sim)
}
