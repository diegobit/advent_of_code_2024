package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

type Set map[any]struct{}

func (s Set) Add(el any) {
	s[el] = struct{}{}
}

func (s Set) Contains(el any) bool {
	_, ok := s[el]
	return ok
}

type Point struct {
	y int
	x int
}

type Dir struct {
	y int
	x int
}

type State struct {
	y   int
	x   int
	dir Dir
}

type Matrix [][]string

func (m *Matrix) String() string {
	s := strings.Builder{}
	for _, row := range *m {
		for _, cell := range row {
			s.WriteString(fmt.Sprintf("%s ", cell))
		}
		s.WriteString("\n")
	}
	return s.String()
}

type Maze struct {
	mat   *Matrix
	start State
	end   Point
}

type Move struct {
	state State
	cost  int
}

func (m *Maze) String() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("Maze(start=%v, end=%v\nmat=\n%v)\n", m.start, m.end, m.mat))
	return s.String()
}

func (m *Maze) Clone() *Maze {
	cloned := &Maze{}
	cloned.start = m.start
	cloned.end = m.end
	cloned.mat = &Matrix{}
	*cloned.mat = make(Matrix, len(*m.mat))
	for y := range *m.mat {
		(*cloned.mat)[y] = make([]string, len((*m.mat)[0]))
		copy((*cloned.mat)[y], (*m.mat)[y])
	}
	return cloned
}

func readFile(path string) (maze *Maze) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	maze = &Maze{}
	maze.start = State{}
	maze.mat = &Matrix{}

	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		x := 0
		line := scanner.Text()
		row := make([]string, 0)
		for _, cell := range line {
			row = append(row, string(cell))
			switch cell {
			case 'S':
				maze.start = State{y, x, Dir{0, 1}}
			case 'E':
				maze.end = Point{y, x}
			}
			x++
		}
		*maze.mat = append(*maze.mat, row)
		y++
	}

	return
}

func getNeighbors(s State, visited Set, m *Maze) []Move {
	neighbors := make([]Move, 0)
	candidates := make([]Move, 0)

	if (*m.mat)[s.y+s.dir.y][s.x+s.dir.x] != "#" {
		candidates = append(candidates, Move{State{s.y + s.dir.y, s.x + s.dir.x, s.dir}, 1}) // forward
	}

	if s.dir.y != 0 {
		candidates = append(candidates, Move{State{s.y, s.x, Dir{0, -1}}, 1000}) // face up
		candidates = append(candidates, Move{State{s.y, s.x, Dir{0, 1}}, 1000})  // face down
	}
	if s.dir.x != 0 {
		candidates = append(candidates, Move{State{s.y, s.x, Dir{-1, 0}}, 1000}) // face left
		candidates = append(candidates, Move{State{s.y, s.x, Dir{1, 0}}, 1000})  // face right
	}

	for _, curr := range candidates {
		if !visited.Contains(curr) {
			neighbors = append(neighbors, curr)
			visited.Add(curr)
		}
	}

	return neighbors
}

func dijkstra(maze *Maze) (dists map[State]int, prevs map[State]([]*State), finalState State) {
	startItem := &Item{maze.start, 0, 0}
	pq := NewPriorityQueue()
	pq.Add(startItem)
	s2i := make(map[State]*Item, 0)
	s2i[maze.start] = startItem
	dists = make(map[State]int)
	prevs = make(map[State]([]*State))
	visited := Set{}
	visited.Add(Move{maze.start, 0})

	idx := 1
	for y := 0; y < len(*maze.mat); y++ {
		for x := 0; x < len((*maze.mat)[0]); x++ {
			if (*maze.mat)[y][x] != "#" {
				for _, dir := range []Dir{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
					curr := State{y, x, dir}
					dists[curr] = math.MaxInt/2 + idx
					prevs[curr] = make([]*State, 0)
					currItem := &Item{curr, math.MaxInt/2 + idx, idx}
					pq.Add(currItem)
					s2i[curr] = currItem
					idx++
				}
			}
		}
	}
	dists[maze.start] = 0 // Starting point of the search

	for pq.Len() > 0 {
		item := pq.PopMin()
		curr := item.value
		neighbors := getNeighbors(curr, visited, maze)
		for _, neighbor := range neighbors {
			neighborDist := dists[curr] + neighbor.cost
			if neighborDist < dists[neighbor.state] {
				dists[neighbor.state] = neighborDist
				prevs[neighbor.state] = []*State{&curr}
				pq.updatePriority(s2i[neighbor.state], neighborDist)
			} else if neighborDist == dists[neighbor.state] {
				dists[neighbor.state] = neighborDist
				// prevs[neighbor.state] = &curr
				prevs[neighbor.state] = append(prevs[neighbor.state], &curr)
				pq.updatePriority(s2i[neighbor.state], neighborDist)
			}
		}
	}

	// Find the final state
	finalState = State{maze.end.y, maze.end.x, Dir{0, 1}}
	otherCandidates := []State{
		{maze.end.y, maze.end.x, Dir{0, -1}},
		{maze.end.y, maze.end.x, Dir{-1, 0}},
		{maze.end.y, maze.end.x, Dir{+1, 0}},
	}
	for _, state := range otherCandidates {
		if dists[state] < dists[finalState] {
			finalState = state
		}
	}

	return dists, prevs, finalState
}

func drawPath(maze *Maze, prevs map[State]([]*State), finalState State) (mazeTraced *Maze, nVisitedPts int) {
	mazeTraced = maze.Clone()
	visitedStates := Set{}
	visitedPts := Set{}
	currs := prevs[finalState]
	for currs != nil && len(currs) > 0 {
		news := make([]*State, 0)
		for _, curr := range currs {
			if !visitedStates.Contains(*curr) {
				if (*mazeTraced.mat)[curr.y][curr.x] == "." {
					(*mazeTraced.mat)[curr.y][curr.x] = "O"
				}
				visitedStates.Add(*curr)
				visitedPts.Add(Point{curr.x, curr.y})
				news = append(news, prevs[*curr]...)
			}
		}
		currs = news
	}
	nVisitedPts = 1 + len(visitedPts) // Adding E manually
	return
}

func simulate(maze *Maze) (nVisitedPts int) {
	dists, prevs, finalState := dijkstra(maze)
	fmt.Printf(("finalState: %v\n"), finalState)

	mazeTraced, nVisitedPts := drawPath(maze, prevs, finalState)
	fmt.Printf("Traced: %v", mazeTraced)

	minPathCost := dists[finalState]
	fmt.Printf("minPathCost: %d\n", minPathCost)

	return
}

func Problem(path string) (score int) {
	data := readFile(path)
	fmt.Printf("%v\n", data)

	score = simulate(data)
	fmt.Printf("Score: %d\n", score)
	return
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no path")
		os.Exit(1)
	}
	path := os.Args[1]
	Problem(path)
}
