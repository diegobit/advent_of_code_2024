package main

import "testing"

func BenchmarkDay31(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		b.StartTimer()
		// sorting.BubbleSort(numbers)
		MyMain()
	}
}
