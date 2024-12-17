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
			s.WriteString(fmt.Sprintf("%-2s", "•"))
			continue
		}
		s.WriteString(fmt.Sprintf("%-2d", elem))
	}
	return s.String()
}

func nextFreeIdx(disk *Disk, i int) int {
	for (*disk)[i] != -1 && i < len(*disk) {
		i += 1
	}
	return i
}

func nextBlockToMove(disk *Disk, cursor int) (start int, blockSize int) {
	val := (*disk)[cursor]
	for j := cursor - 1; j > 0; j-- {
		if (*disk)[j] != val {
			start = j + 1
			blockSize = cursor - start + 1
			return
		}
	}
	start = 0
	blockSize = 0
	return
}

func nextFreeBlock(disk *Disk, cursor int, maxCursor int, size int) (start int) {
	for i := cursor; i < maxCursor; i++ {
		isValid := nextFreeSpaceHelper(disk, i, size)
		if isValid {
			start = i
			return
		}
	}
	return -1
}

func nextFreeSpaceHelper(disk *Disk, cursor int, size int) bool {
	if size == 0 {
		return true
	}
	if (*disk)[cursor] == -1 {
		return nextFreeSpaceHelper(disk, cursor+1, size-1)
	}
	return false
}

func move(disk *Disk, dst int, src int) {
	(*disk)[dst] = (*disk)[src]
	(*disk)[src] = -1
}

func moveBlock(disk *Disk, dstStart int, srcStart int, blockSize int) {
	for i := 0; i < blockSize; i++ {
		// fmt.Printf("moving %d from %d to %d\n", (*disk)[srcStart+k], srcStart+k, dstStart+k)
		move(disk, dstStart+i, srcStart+i)
	}
}

func compactDisk(disk *Disk) {
	i := nextFreeIdx(disk, 0)
	j := len(*disk) - 1

	for {
		if (*disk)[j] == -1 {
			j--
			continue
		}
		if i >= j {
			break
		}
		srcStart, blockSize := nextBlockToMove(disk, j)
		dstStart := nextFreeBlock(disk, i, srcStart, blockSize)
		if dstStart != -1 && srcStart != -1 {
			moveBlock(disk, dstStart, srcStart, blockSize)
			// fmt.Printf("DISK AFT: %v\n", disk)
			i = nextFreeIdx(disk, i)
		}
		j -= blockSize
	}
}

func computeChecksum(disk *Disk) (cksum int64) {
	cksum = 0
	for i, id := range *disk {
		if id == -1 {
			continue
		}
		cksum += int64(i) * int64(id)
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
			if i%2 == 0 {
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

	if len(disk) < 1000 {
		fmt.Printf("DISK BEF: %v\n", disk)
	} else {
		fmt.Printf("DISK BEF: %v[...]%v\n", disk[0:20], disk[len(disk)-20:])
	}

	compactDisk(&disk)

	if len(disk) < 1000 {
		fmt.Printf("DISK END: %v\n", disk)
	} else {
		fmt.Printf("DISK END: %v[...]%v\n", disk[0:20], disk[len(disk)-20:])
	}
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
