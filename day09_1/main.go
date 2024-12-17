package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Disk []int

func (d Disk) String() string {
	s := strings.Builder{}
	for _, elem := range d {
		if elem == -1 {
			s.WriteString(".")
			continue
		}
		s.WriteString(fmt.Sprintf("%d", elem))
	}
	return s.String()
}

func nextFreeIdx(disk *Disk, i int) int {
	for (*disk)[i] != -1 && i < len(*disk) {
		i += 1
	}
	return i
}

func move(disk *Disk, dst int, src int) {
	(*disk)[dst] = (*disk)[src]
	(*disk)[src] = -1
}

func compactDisk(disk *Disk) {
	i := nextFreeIdx(disk, 0)
	for j := len(*disk)-1; j > 0; j-- { // j, cursor of next block to move
		if i >= j {
			break
		}
		move(disk, i, j)
		i = nextFreeIdx(disk, i)
	}
}

func computeChecksum(disk *Disk) (cksum int) {
	cksum = 0
	for i, id := range *disk {
		if id == -1 {
			continue
		}
		cksum += i*id
	}
	return
}

func Problem(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	disk := Disk{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		diskMap := scanner.Text()
		currId := 0
		for i, elem := range strings.Split(diskMap, "") {
			val, err := strconv.Atoi(elem)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if i % 2 == 0 {
				// block
				for j := 0; j < val; j++ {
					disk = append(disk, currId)
				}
				currId += 1
			} else {
				// free space
				for j := 0; j < val; j++ {
					disk = append(disk, -1)
				}
			}
		}
	}

	fmt.Printf("DISK BEF: %v\n", disk)

	compactDisk(&disk)

	fmt.Printf("DISK AFT: %v\n", disk)
	fmt.Printf("cksum: %d\n", computeChecksum(&disk))
}


func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no path")
		os.Exit(1)
	}
	path := os.Args[1]
	Problem(path)
}
