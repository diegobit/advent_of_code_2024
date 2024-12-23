package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readFile(path string) (stones []int64) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stoneStrings := strings.Split(line, " ")
		for _, stone := range stoneStrings {
			stoneInt, err := strconv.ParseInt(stone, 10, 64)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			stones = append(stones, stoneInt)
		}
	}

	return
}

func evenDigits(v int64) bool {
	s := strconv.FormatInt(v, 10)
	return len(s)%2 == 0
}

func splitDigits(v int64) (int64, int64) {
	s := strconv.FormatInt(v, 10)
	sleft, sright := s[0:len(s)/2], s[len(s)/2:]
	left, _ := strconv.ParseInt(sleft, 10, 64)
	right, _ := strconv.ParseInt(sright, 10, 64)
	return left, right
}

func blink(stones *[]int64, nBlinks int, nBlinksPerStep int, cache *map[int64](*[]int64)) (nStones int64) {
	if nBlinks <= 0 {
		nStones = int64(len(*stones))
		return
	}
	// fmt.Printf("nBlinks: %d, nBlinksPerStep: %d, depth: %d\n", nBlinks, nBlinksPerStep, depth)

	for _, s := range *stones {
		newStones, ok := (*cache)[s]
		if !ok {
			newStones = blinkStep(s, nBlinksPerStep)
			(*cache)[s] = newStones
		}
		nb := blink(newStones, nBlinks-nBlinksPerStep, nBlinksPerStep, cache)
		nStones += nb
		// if nBlinks == 75 {
		// 	fmt.Printf("Stone: %-7d, newStones: %-10d, TOTAL: %-15d\n", s, nb, nStones)
		// }
	}

	return
}

func blinkStep(stone int64, nBlinks int) *[]int64 {
	currStack := &[]int64{stone}
	nextStack := &[]int64{}
	for i := 0; i < nBlinks; i++ {
		for _, s := range *currStack {
			if s == 0 {
				(*nextStack) = append(*nextStack, 1)
			} else if evenDigits(s) {
				lh, rh := splitDigits(s)
				(*nextStack) = append(*nextStack, lh)
				(*nextStack) = append(*nextStack, rh)
			} else {
				(*nextStack) = append(*nextStack, s*2024)
			}
		}
		currStack = nextStack
		nextStack = &[]int64{}
	}

	return currStack
}

func Problem(path string) {
	stones := readFile(path)
	fmt.Println(stones)

	cache := make(map[int64](*[]int64))
	nStones := blink(&stones, 75, 25, &cache)

	fmt.Println(nStones)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no path")
		os.Exit(1)
	}
	path := os.Args[1]
	Problem(path)
}
