package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*********
** MATRIX
**********/
type Cell interface {
	string | int
}

type Matrix[C Cell] [][]C

/*********
** REGION
**********/
type Point struct {
	y int
	x int
}

type Region struct {
	plant string
	pts   []Point
	pmt   int
}

/*********
** FUNCS
**********/
func (m Matrix[C]) String() string {
	s := strings.Builder{}
	for _, row := range m {
		s.WriteString(fmt.Sprintf("%v\n", row))
	}
	return s.String()
}

func readFile(path string) (plantMap Matrix[string]) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := make([]string, 0)
		line := scanner.Text()
		plants := strings.Split(line, "")
		for _, plant := range plants {
			row = append(row, plant)
		}
		plantMap = append(plantMap, row)
	}

	return
}

func insideBorders(pt Point, nRows int, nCols int) bool {
	return pt.y >= 0 && pt.y < nRows && pt.x >= 0 && pt.x < nCols
}

func addFencesTo(region *Region, plantMap Matrix[string], pt Point) {
	plant := plantMap[pt.y][pt.x]
	nRows, nCols := len(plantMap), len(plantMap[0])

	simLeft := func() bool {
		ptOther := Point{pt.y, pt.x - 1}
		return insideBorders(ptOther, nRows, nCols) && plantMap[ptOther.y][ptOther.x] == plant
	}
	simTopLeft := func() bool {
		ptOther := Point{pt.y - 1, pt.x - 1}
		return insideBorders(ptOther, nRows, nCols) && plantMap[ptOther.y][ptOther.x] == plant
	}
	simTop := func() bool {
		ptOther := Point{pt.y - 1, pt.x}
		return insideBorders(ptOther, nRows, nCols) && plantMap[ptOther.y][ptOther.x] == plant
	}
	simTopRight := func() bool {
		ptOther := Point{pt.y - 1, pt.x + 1}
		return insideBorders(ptOther, nRows, nCols) && plantMap[ptOther.y][ptOther.x] == plant
	}

	// See addFencesTo_logic.jpg for documentation
	if !simLeft() && !simTopLeft() && !simTop() && !simTopRight() {
		region.pmt += 4
	} else if !simLeft() && simTopLeft() && !simTop() && !simTopRight() {
		region.pmt += 4
	} else if !simLeft() && !simTopLeft() && !simTop() && simTopRight() {
		region.pmt += 4
	} else if simLeft() && simTopLeft() && !simTop() && !simTopRight() {
		region.pmt += 2
	} else if simLeft() && !simTopLeft() && simTop() && !simTopRight() {
		region.pmt -= 2
	} else if !simLeft() && simTopLeft() && simTop() && !simTopRight() {
		region.pmt += 2
	} else if !simLeft() && simTopLeft() && !simTop() && simTopRight() {
		region.pmt += 4
	} else if !simLeft() && !simTopLeft() && simTop() && simTopRight() {
		region.pmt += 2
	} else if simLeft() && simTopLeft() && simTop() && !simTopRight() {
		region.pmt -= 2
	} else if simLeft() && simTopLeft() && !simTop() && simTopRight() {
		region.pmt += 2
	} else if !simLeft() && simTopLeft() && simTop() && simTopRight() {
		region.pmt += 4
	}
}

func getNeighborId(plant string, pt Point, plantMap Matrix[string], pt2id map[Point]int) int {
	nRows, nCols := len(plantMap), len(plantMap[0])
	if pt.y >= 0 && pt.y < nRows && pt.x >= 0 && pt.x < nCols && plantMap[pt.y][pt.x] == plant {
		return pt2id[pt]
	}
	return -1
}

func getRegions(plantMap Matrix[string]) (regions map[int]*Region) {
	nextFreeId := 0
	regions = make(map[int]*Region)
	pt2id := make(map[Point]int)

	nRows, nCols := len(plantMap), len(plantMap[0])
	for y := 0; y < nRows; y++ {
		for x := 0; x < nCols; x++ {
			pt := Point{y, x}
			plant := plantMap[y][x]
			var region *Region = nil
			topPt := Point{pt.y - 1, pt.x}
			leftPt := Point{pt.y, pt.x - 1}
			neighborTopId := getNeighborId(plant, topPt, plantMap, pt2id)
			neighborLeftId := getNeighborId(plant, leftPt, plantMap, pt2id)

			addToTop := func() {
				region = regions[neighborTopId]
				region.pts = append(region.pts, pt)
				pt2id[pt] = neighborTopId
				addFencesTo(region, plantMap, pt)
			}

			addToLeft := func() {
				region = regions[neighborLeftId]
				region.pts = append(region.pts, pt)
				pt2id[pt] = neighborLeftId
				addFencesTo(region, plantMap, pt)
			}

			createNew := func() {
				regions[nextFreeId] = &Region{plant, []Point{pt}, 0}
				pt2id[pt] = nextFreeId
				addFencesTo(regions[nextFreeId], plantMap, pt)
				nextFreeId++
			}

			mergeTopWithLeft := func() {
				region = regions[neighborTopId]
				region.pts = append(region.pts, regions[neighborLeftId].pts...)
				// region.pmt = append(region.pmt, regions[neighborLeftId].pmt...)
				region.pmt += regions[neighborLeftId].pmt
				for _, leftPt := range regions[neighborLeftId].pts {
					pt2id[leftPt] = neighborTopId
				}
				delete(regions, neighborLeftId)
			}

			if neighborTopId != -1 && neighborLeftId == -1 {
				addToTop()
			} else if neighborTopId == -1 && neighborLeftId != -1 {
				addToLeft()
			} else if neighborTopId != -1 && neighborLeftId != -1 {
				// Add to top, then merge left region with top region
				addToTop()
				if neighborTopId != neighborLeftId {
					mergeTopWithLeft()
				}
			} else {
				createNew()
			}
		}
	}
	return
}

func getPrice(regions map[int]*Region) (totalPrice int) {
	for id, region := range regions {
		area := len(region.pts)
		perimeter := region.pmt
		price := area * perimeter
		fmt.Printf("%d(%s):\tA * P = $\t%d * %d = %d\n", id, region.plant, area, perimeter, price)
		totalPrice += price
	}
	return
}

func Problem(path string) {
	plantMap := readFile(path)
	fmt.Println(plantMap)

	regions := getRegions(plantMap)
	fmt.Println("Regions:")
	fmt.Println(regions)
	fmt.Println()

	fmt.Println("Price calc:")
	fmt.Println(getPrice(regions))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no path")
		os.Exit(1)
	}
	path := os.Args[1]
	Problem(path)
}
