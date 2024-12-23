package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type LL struct {
	value int64
	prev  *LL
	next  *LL
}

func (ll *LL) Add(value int64) *LL {
	new := &LL{value, ll, ll.next}
	new.prev.next = new
	if ll.next != nil {
		ll.next.prev = new
		return ll.next
	} else {
		return new
	}
}

func (ll *LL) Len() int {
	len := 1
	if ll.next != nil {
		curr := ll
		for curr := curr.next; curr != nil; {
			len++
			curr = curr.next
		}
	}
	if ll.prev != nil {
		curr := ll
		for curr := curr.prev; curr != nil; {
			len++
			curr = curr.prev
		}
	}
	return len
}

func (ll *LL) String() string {
	s := strings.Builder{}
	curr := ll
	s.WriteString("LL{ ")
	for curr != nil {
		s.WriteString(fmt.Sprintf("%d ", curr.value))
		curr = curr.next
	}
	s.WriteString("}")
	return s.String()
}

func readFile(path string) (ll *LL, endLL *LL) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stones := strings.Split(line, " ")
		for i, stone := range stones {
			stoneInt, err := strconv.ParseInt(stone, 10, 64)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if i == 0 {
				ll = &LL{stoneInt, nil, nil}
				endLL = ll
			} else {
				endLL = endLL.Add(stoneInt)
			}
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

func blink(startLL *LL, nBlinks int) int {
	for i := 0; i < nBlinks; i++ {
		fmt.Println(i)
		curr := startLL
		for curr != nil {
			next := curr.next
			switch v := curr.value; {
			case v == 0:
				curr.value = 1
			case evenDigits(v):
				lh, rh := splitDigits(v)
				curr.value = lh
				curr.Add(rh)
			default:
				curr.value *= 2024
			}
			curr = next
		}
	}

	return startLL.Len()
}

func Problem(path string) {
	startLL, _ := readFile(path)
	fmt.Println(startLL)
	nStones := blink(startLL, 75)
	// fmt.Println(startLL)
	fmt.Println(nStones)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("ERROR: no path")
		os.Exit(1)
	}
	path := os.Args[1]
	// path := "input_sample_1_55312.txt"
	Problem(path)
}
