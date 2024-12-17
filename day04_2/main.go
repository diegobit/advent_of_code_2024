package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Matrix = [][]string

type Dir = []int

func search(i int, j int, data Matrix) int {
	if data[i][j] != "A" {
		return 0
	}

	directions := []Dir{
		{1, 1},
		{1, -1},
		{-1, 1},
		{-1, -1},
	}

	nRows := len(data)
	nCols := len(data[0])
	min_i := i-1
	max_i := i+1
	min_j := j-1
	max_j := j+1

	for _, dir := range directions {
		dv := dir[0]
		dh := dir[1]
		if min_i >= 0 && max_i < nRows && min_j >= 0 && max_j < nCols {
			// fmt.Printf("i=%d, j=%d, dv=%d, dh=%d\n", i, j, dv, dh)
			// fmt.Printf("i3dv=%d, j3dh=%d\n", i+3*dv, j+3*dh)
			// fmt.Printf("%s%s%s\n", data[i+1*dv][j+1*dh], data[i+2*dv][j+2*dh], data[i+3*dv][j+3*dh])
			if data[i + 1*dv][j + 1*dh] == "M" {
				if data[i - 1*dv][j - 1*dh] == "S" {

					if data[i + 1*dv][j - 1*dh] == "M" {
						if data[i - 1*dv][j + 1*dh] == "S" {
							fmt.Printf("+--+ MATCH: %d, %d\n", i, j)
							return 1
						}
					} else if data[i - 1*dv][j + 1*dh] == "M" {
						if data[i + 1*dv][j - 1*dh] == "S" {
							fmt.Printf("-++- MATCH: %d, %d\n", i, j)
							return 1
						}
					}

				}
			}
		}
	}

	return 0
}

func Xmas() {
	path := "input.txt"
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	nRows := 0
	nCols := 0
	data := make(Matrix, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > nCols {
			nCols = len(line)
		}
		nRows += 1
		data = append(data, strings.Split(line, ""))
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(data[0][0])
	fmt.Println(data[0:2])

	nMatches := 0

	for i := 0; i < nRows; i++ {
		for j := 0; j < nCols; j++ {
			nMatches += search(i, j, data)
		}
	}

	fmt.Println(nMatches)


}

func main() {
	Xmas()
}

