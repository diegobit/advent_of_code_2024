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
	pts []Point
	pmt []Point
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

func addFencesTo(region *Region, pMap Matrix[string], p Point, nRows int, nCols int) {
	if region.pmt == nil {
		region.pmt = make([]Point, 0)
	}
	plant := pMap[p.y][p.x]
	if p.y-1 < 0 || pMap[p.y-1][p.x] != plant {
		region.pmt = append(region.pmt, Point{p.y-1, p.x})
	}
	if p.y+1 >= nRows || pMap[p.y+1][p.x] != plant {
		region.pmt = append(region.pmt, Point{p.y+1, p.x})
	}
	if p.x-1 < 0 || pMap[p.y][p.x-1] != plant {
		region.pmt = append(region.pmt, Point{p.y, p.x-1})
	}
	if p.x+1 >= nCols || pMap[p.y][p.x+1] != plant {
		region.pmt = append(region.pmt, Point{p.y, p.x+1})
	}
}

func getNeighborId(plant string, pt Point, plantMap Matrix[string], pt2id map[Point]int, nRows int, nCols int) int {
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
			pt := Point{y,x}
			plant := plantMap[y][x]
			var region *Region = nil
			topPt := Point{pt.y-1, pt.x}
			leftPt := Point{pt.y, pt.x-1}
			neighborTopId := getNeighborId(plant, topPt, plantMap, pt2id, nRows, nCols)
			neighborLeftId := getNeighborId(plant, leftPt, plantMap, pt2id, nRows, nCols)

			addToTop := func() {
				region = regions[neighborTopId]
				region.pts = append(region.pts, pt)
				pt2id[pt] = neighborTopId
				addFencesTo(region, plantMap, pt, nRows, nCols)
			}

			addToLeft := func() {
				region = regions[neighborLeftId]
				region.pts = append(region.pts, pt)
				pt2id[pt] = neighborLeftId
				addFencesTo(region, plantMap, pt, nRows, nCols)
			}

			createNew := func() {
				regions[nextFreeId] = &Region{plant, []Point{pt}, []Point{}}
				pt2id[pt] = nextFreeId
				addFencesTo(regions[nextFreeId], plantMap, pt, nRows, nCols)
				nextFreeId++
			}

			mergeTopWithLeft := func() {
				region = regions[neighborTopId]
				region.pts = append(region.pts, regions[neighborLeftId].pts...)
				region.pmt = append(region.pmt, regions[neighborLeftId].pmt...)
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
		perimeter := len(region.pmt)
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
