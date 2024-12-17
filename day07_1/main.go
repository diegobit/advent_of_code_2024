package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Op int

const (
	None = iota
	OpAdd
	OpMul
)

func (o Op) String() string {
	names := map[Op]string{
		None:  ".",
		OpAdd: "+",
		OpMul: "*",
	}

	if name, ok := names[o]; ok {
		return name
	}
	return fmt.Sprintf("Unknown(%d)", o)
}

var ops = []Op{OpAdd, OpMul}

func getPermutationHelper(result *[][]Op, nOps int, currOps []Op) {
	if nOps == 0 {
		if currOps != nil && len(currOps) > 0 {
			// fmt.Printf("currOps: %+v\n", currOps)
			*result = append(*result, currOps)
		}
		return
	}
	if currOps == nil {
		currOps = []Op{}
	}
	for _, op := range ops {
		// fmt.Printf("BEGIN: nOps: %d, currOps: %+v, op: %s\n", nOps, currOps, op)
		newOps := make([]Op, len(currOps)+1)
		if len(currOps) > 0 {
			copy(newOps[0:len(currOps)], currOps)
		}
		// fmt.Printf("AFTCP: currOps: %+v, newOps: %+v\n", currOps, newOps)
		newOps[len(currOps)] = op
		// fmt.Printf("AFTAS: currOps: %+v, newOps: %+v\n", currOps, newOps)
		getPermutationHelper(result, nOps-1, newOps)
	}
}

func getPermutations(nOps int) *[][]Op {
	if nOps < 1 {
		return nil
	}

	result := make([][]Op, nOps)
	getPermutationHelper(&result, nOps, nil)

	return &result
}

func Problem(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	total := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		equation := strings.Split(line, ":")
		lh, _ := strconv.Atoi(equation[0])
		rhToks := strings.Split(strings.TrimSpace(equation[1]), " ")
		rh := []int{}
		for _, w := range rhToks {
			i, err := strconv.Atoi(w)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			rh = append(rh, i)
		}

		nOps := len(rh) - 1

		permutations := getPermutations(nOps)
		fmt.Printf("%d: %s\n", lh, rhToks)
		// fmt.Printf("Num p: %d\n", nOps)

		// for _, permutation := range *permutations {
		// 	fmt.Println(permutation)
		// }
		// fmt.Println("")

		for _, permutation := range *permutations {
			if permutation == nil {
				continue
			}

			curr := rh[0]
			for i := 0; i < len(rh)-1; i++ {
				op := permutation[i]
				if op == OpAdd {
					curr = curr + rh[i+1]
				} else if op == OpMul {
					curr = curr * rh[i+1]
				}
			}

			if lh == curr {
				fmt.Printf("👉 Found! Total: %d -> %d\n", total, total+curr)
				total += curr
				break
			}
		}
	}

	fmt.Printf("Total: %d\n", total)

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no path")
		os.Exit(1)
	}
	path := os.Args[1]
	Problem(path)
}
