package main

import (
	"container/heap"
	"fmt"
)

// Define a custom type for the heap elements
type Item struct {
	value    string
	priority int
	index    int
}

// Define a priority queue type
type PriorityQueue []*Item

// Implement the heap.Interface methods for the PriorityQueue type
func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func main() {
	// Create an empty priority queue
	pq := make(PriorityQueue, 0)

	// Push items into the priority queue
	heap.Push(&pq, &Item{value: "item1", priority: 3})
	heap.Push(&pq, &Item{value: "item2", priority: 1})
	heap.Push(&pq, &Item{value: "item3", priority: 2})

	// Pop items from the priority queue
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("Value: %s, Priority: %d\n", item.value, item.priority)
	}
}
