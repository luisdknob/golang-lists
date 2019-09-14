package list

import (
	"fmt"
	"sync"
)

// CoarseList definition with Mutex
type CoarseList struct {
	head *Node
	mux  sync.Mutex
}

// Print is
func (l *CoarseList) Print() {
	if l.head.next != nil {
		tmp := l.head
		for {
			fmt.Printf("Valor: %d\n", tmp.item)
			if tmp.next != nil {
				tmp = tmp.next
			} else {
				break
			}
		}
	}
}

// Add is
func (l *CoarseList) Add(item int) bool {
	var pred, curr *Node
	l.mux.Lock()
	defer l.mux.Unlock()

	pred = l.head
	curr = pred.next
	for curr.item < item {
		pred = curr
		curr = curr.next
	}
	if item == curr.item {
		return false
	}

	node := &Node{
		item: item,
	}
	node.next = curr
	pred.next = node

	return true
}

// Contains is
func (l *CoarseList) Contains(item int) bool {
	var pred, curr *Node
	l.mux.Lock()
	defer l.mux.Unlock()

	pred = l.head
	curr = pred.next
	for curr.item < item {
		pred = curr
		curr = curr.next
	}
	if item == curr.item {
		return true
	}

	return false
}

// Remove is
func (l *CoarseList) Remove(item int) bool {
	var pred, curr *Node
	l.mux.Lock()
	defer l.mux.Unlock()

	pred = l.head
	curr = pred.next
	for curr.item < item {
		pred = curr
		curr = curr.next
	}
	if item == curr.item {
		pred.next = curr.next
		return true
	}

	return false
}

//Count is
func (l *CoarseList) Count() (num int) {
	num = 0
	last := l.head
	for {
		if last.next == nil {
			break
		} else {
			num++
			last = last.next
		}
	}
	num--
	return
}

// NewCoarseList is
func NewCoarseList() (l *CoarseList) {
	head := Node{
		item: 0,
		next: &Node{
			item: int(^uint(0) >> 1),
		},
	}
	l = &CoarseList{
		head: &head,
	}
	return
}
