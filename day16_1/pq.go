package main

import (
	"container/heap"
)

var _ heap.Interface = (*PriorityQueue)(nil)

type Item struct {
	value    State
	priority int
	index    int // The index is needed by update and is maintained by the heap.Interface methods.
}

type PriorityQueue []*Item

/////
// heap.Interface methods
/////

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority // min-priority queue
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	item := x.(*Item)
	n := len(*pq)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

/////
// Other methods
/////

func NewPriorityQueue() *PriorityQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	return &pq
}

func (pq *PriorityQueue) Add(item *Item) {
	heap.Push(pq, item)
}

func (pq *PriorityQueue) PopMin() *Item {
	return heap.Pop(pq).(*Item)
}

func (pq *PriorityQueue) updateValue(item *Item, value State) {
	item.value = value
	heap.Fix(pq, item.index)
}

func (pq *PriorityQueue) updatePriority(item *Item, priority int) {
	item.priority = priority
	heap.Fix(pq, item.index)
}
