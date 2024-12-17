package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// Assumptions:
// - no conflicting rules

const (
	ReadRules = iota
	ReadManual
)

func Manual() {
	path := "input.txt"
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	rules := map[string][]string{}
	// sortedManuals := []string{}
	sum := 0

	mode := ReadRules

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		switch mode {
		case ReadRules:
			if line == "" {
				mode = ReadManual
				continue
			}
			values := strings.Split(line, "|")
			_, ok := rules[values[0]]
			if !ok {
				rules[values[0]] = []string{}
			}
			rules[values[0]] = append(rules[values[0]], values[1])

		case ReadManual:
			pages := strings.Split(line, ",")
			correct := true
			// fmt.Println(pages)
			sort.Slice(pages, func(i, j int) bool {
				a := pages[i]
				b := pages[j]
				rule := rules[b]
				if slices.Contains(rule, a) {
					return false
				}
				correct = false
				return true
			})
			// fmt.Println(pages)
			// fmt.Println()
			// sortedManuals = append(sortedManuals, strings.Join(pages, ","))
			val, err := strconv.Atoi(pages[(len(pages)-1)/2])
			// fmt.Println(pages)
			// fmt.Printf("val: %d, page: %s\n", val, pages[(len(pages)-1)/2])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			// fmt.Println(val)
			if correct {
				sum += val
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(sum)
}

func main() {
	Manual()
}
