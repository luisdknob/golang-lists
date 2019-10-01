// http://scottlobdell.me/2016/09/thread-safe-non-blocking-linked-lists-golang/

package list

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// NonBlockingList definition with Lock Free Pointer
type NonBlockingList struct {
	head *LockFreeNode
}

//LockFreeNode is ...
type LockFreeNode struct {
	markableNext *MarkablePointer
	item         int
}

//MarkablePointer is ..
type MarkablePointer struct {
	marked bool
	next   *LockFreeNode
}

// NewNonBlockingList is
func NewNonBlockingList() (l *NonBlockingList) {

	head := LockFreeNode{
		item: 0,
		markableNext: &MarkablePointer{
			marked: false,
			next: &LockFreeNode{
				item: int(^uint(0) >> 1),
			},
		},
	}
	l = &NonBlockingList{
		head: &head,
	}
	return
}

//Print is
func (l *NonBlockingList) Print() {

	curr := l.head

	for {
		fmt.Println(curr.item)
		if curr.markableNext != nil {
			curr = curr.markableNext.next
		} else {
			break
		}
	}

}

// Remove is
func (l *NonBlockingList) Remove(item int) bool {
	var previous *LockFreeNode
	cursor := l.head
	for {
		if cursor == nil {
			break
		}
		if cursor.item == item {
			nextNode := cursor.markableNext.next
			newNext := MarkablePointer{
				marked: true,
				next:   nextNode,
			}
			operationSucceeded := atomic.CompareAndSwapPointer(
				(*unsafe.Pointer)(unsafe.Pointer(&(cursor.markableNext))),
				unsafe.Pointer(cursor.markableNext),
				unsafe.Pointer(&newNext),
			)
			if !operationSucceeded {
				l.Remove(item)
				return true
			}
			newNext = MarkablePointer{
				next: nextNode,
			}
			if previous != nil {
				operationSucceeded = atomic.CompareAndSwapPointer(
					(*unsafe.Pointer)(unsafe.Pointer(&(previous.markableNext))),
					unsafe.Pointer(previous.markableNext),
					unsafe.Pointer(&newNext),
				)
			}
			if !operationSucceeded {
				l.Remove(item)
			}
			break
		} else if cursor.item > item {
			return false
		}

		previous = cursor
		if cursor.markableNext == nil {
			return false
		}
		cursor = cursor.markableNext.next
	}
	return true
}

// Add is
func (l *NonBlockingList) Add(item int) bool {
	cursor := l.head
	for {
		if item < cursor.markableNext.next.item {
			currentNext := cursor.markableNext
			if currentNext.marked {
				continue
			}
			if cursor.item == item {
				return false
			}
			newNode := LockFreeNode{
				item: item,
				markableNext: &MarkablePointer{
					next: currentNext.next,
				},
			}
			newNext := MarkablePointer{
				next: &newNode,
			}
			operationSucceeded := atomic.CompareAndSwapPointer(
				(*unsafe.Pointer)(unsafe.Pointer(&(cursor.markableNext))),
				unsafe.Pointer(currentNext),
				unsafe.Pointer(&newNext),
			)
			if !operationSucceeded {
				l.Add(item)
				return true
			}
			break
		}
		cursor = cursor.markableNext.next
	}
	return true
}

// Contains is
func (l *NonBlockingList) Contains(item int) bool {
	curr := l.head
	marked := false
	for curr.item < item {
		marked = curr.markableNext.marked
		curr = curr.markableNext.next
	}
	return curr.item == item && !marked
}
