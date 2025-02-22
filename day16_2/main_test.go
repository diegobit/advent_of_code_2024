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

	if score != 3 {
		t.Errorf("Expected 3, got %d", score)
	}
}

func TestSimulate2(t *testing.T) {
	maze := readFile("input_test_2.txt")
	t.Log(maze)

	score := simulate(maze)
	t.Logf("Score: %d\n", score)
	t.Logf("Final maze:\n%v\n", maze)

	if score != 4 {
		t.Errorf("Expected 4, got %d", score)
	}
}

func TestInputSample1(t *testing.T) {
	score := Problem("input_sample_1_7036.txt")
	if score != 45 {
		t.Errorf("Expected 45, got %d", score)
	}
}

func TestInputSample2(t *testing.T) {
	score := Problem("input_sample_2_11048.txt")
	if score != 64 {
		t.Errorf("Expected 64, got %d", score)
	}
}
