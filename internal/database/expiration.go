package database

import (
	"container/heap"
	"time"
)

type ExpirationHeap []*HeapItem

type HeapItem struct {
	Key    string
	Expiry time.Time
	Index  int
}

func NewExpirationHeap() *ExpirationHeap {
	h := &ExpirationHeap{}
	heap.Init(h)
	return h
}

func (h ExpirationHeap) Len() int { return len(h) }

func (h ExpirationHeap) Less(i, j int) bool {
	return h[i].Expiry.Before(h[j].Expiry)
}

func (h ExpirationHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

func (h *ExpirationHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*HeapItem)
	item.Index = n
	*h = append(*h, item)
}

func (h *ExpirationHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*h = old[0 : n-1]
	return item
}

func (h *ExpirationHeap) PushKey(key string, expiry time.Time) {
	heap.Push(h, &HeapItem{
		Key:    key,
		Expiry: expiry,
	})
}
