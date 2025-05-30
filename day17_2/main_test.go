package main

import (
	"testing"
)

func TestSimulate1(t *testing.T) {
	maze := readFile("input_test_1.txt")
	t.Log(maze)

	score := simulate(maze)
	t.Logf("Score: %d\n", score)
	t.Logf("Final maze:\n%v\n", maze)

	if score != 2 {
		t.Errorf("Expected 2, got %d", score)
	}
}

func TestSimulate2(t *testing.T) {
	maze := readFile("input_test_2.txt")
	t.Log(maze)

	score := simulate(maze)
	t.Logf("Score: %d\n", score)
	t.Logf("Final maze:\n%v\n", maze)

	if score != 2003 {
		t.Errorf("Expected 2003, got %d", score)
	}
}

func TestInputSample1(t *testing.T) {
	score := Problem("input_sample_1_7036.txt")
	if score != 7036 {
		t.Errorf("Expected 7036, got %d", score)
	}
}

func TestInputSample2(t *testing.T) {
	score := Problem("input_sample_2_11048.txt")
	if score != 11048 {
		t.Errorf("Expected 11048, got %d", score)
	}
}
