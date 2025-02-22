package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Robot struct {
	x  int
	y  int
	vx int
	vy int
}

type Restroom struct {
	robots []*Robot
	lenX   int
	lenY   int
}

func (r Robot) String() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("p=%d,%d v=%d,%d", r.x, r.y, r.vx, r.vy))
	return s.String()
}

func (r Restroom) String() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("Restroom(lenX=%d, lenY=%d, robots=%v)\n", r.lenX, r.lenY, r.robots))
	return s.String()
}

func readFile(path string) (restroom *Restroom) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	re := regexp.MustCompile(`p\=(-?\d+),(-?\d+) v\=(-?\d+),(-?\d+)`)
	restroom = &Restroom{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) < 4 {
			fmt.Printf("ERROR: Not enough matches in '%s'. Found matches: %v\n", line, matches)
			os.Exit(1)
		}
		xStr, yStr, vxStr, vyStr := matches[1], matches[2], matches[3], matches[4]
		var x, y, vx, vy int
		var err error
		if x, err = strconv.Atoi(xStr); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if y, err = strconv.Atoi(yStr); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if vx, err = strconv.Atoi(vxStr); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if vy, err = strconv.Atoi(vyStr); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		restroom.robots = append(restroom.robots, &Robot{x, y, vx, vy})
		if x > restroom.lenX {
			restroom.lenX = x
		}
		if y > restroom.lenY {
			restroom.lenY = y
		}
	}
	// Make a length
	restroom.lenX += 1
	restroom.lenY += 1

	return
}

func simulate(restroom *Restroom) (newRestroom *Restroom) {
	for step := 0; step < 100; step++ {
		for _, robot := range restroom.robots {
			robot.x += robot.vx
			robot.y += robot.vy
		}
	}
	// Compute modulo only once
	for _, robot := range restroom.robots {
		robot.x = robot.x % restroom.lenX
		if robot.x < 0 {
			robot.x += restroom.lenX
		}
		robot.y = robot.y % restroom.lenY
		if robot.y < 0 {
			robot.y += restroom.lenY
		}
	}
	return restroom
}

func computeSafetyFactor(restroom *Restroom) (factor int) {
	sumRobotsTopLeft := 0
	sumRobotsTopRight := 0
	sumRobotsBottomLeft := 0
	sumRobotsBottomRight := 0
	middleX := (restroom.lenX-1)/2
	middleY := (restroom.lenY-1)/2
	for _, robot := range restroom.robots {
		if robot.x < middleX && robot.y < middleY {
			sumRobotsTopLeft++
		} else if robot.x > middleX && robot.y < middleY {
			sumRobotsTopRight++
		} else if robot.x < middleX && robot.y > middleY {
			sumRobotsBottomLeft++
		} else if robot.x > middleX && robot.y > middleY {
			sumRobotsBottomRight++
		}
	}
	fmt.Printf("sumRobotsTopLeft: %d, sumRobotsTopRight: %d, sumRobotsBottomLeft: %d, sumRobotsBottomRight: %d\n", sumRobotsTopLeft, sumRobotsTopRight, sumRobotsBottomLeft, sumRobotsBottomRight)

	factor = sumRobotsTopLeft * sumRobotsTopRight * sumRobotsBottomLeft * sumRobotsBottomRight

	return
}

func Problem(path string) {
	restroom := readFile(path)
	fmt.Println(restroom)
	fmt.Println()

	newRestroom := simulate(restroom)
	fmt.Println(newRestroom)
	fmt.Println()

	fmt.Printf("TOTAL: %d\n", computeSafetyFactor(newRestroom))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no path")
		os.Exit(1)
	}
	path := os.Args[1]
	Problem(path)
}
