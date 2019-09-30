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
	object       unsafe.Pointer
}

//MarkablePointer is ..
type MarkablePointer struct {
	marked bool
	next   *LockFreeNode
}

// NewNonBlockingList is
func NewNonBlockingList() (l *NonBlockingList) {
	headValue := 0
	tailValue := int(^uint(0) >> 1)

	head := LockFreeNode{
		object: unsafe.Pointer(&headValue),
		markableNext: &MarkablePointer{
			marked: false,
			next: &LockFreeNode{
				object: unsafe.Pointer(&tailValue),
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
		fmt.Println(*(*int)(curr.object))
		if curr.markableNext != nil {
			curr = curr.markableNext.next
		} else {
			break
		}
	}

}

func (l *NonBlockingList) deleteObject(object unsafe.Pointer) bool {
	var previous *LockFreeNode
	cursor := l.head
	for {
		if cursor == nil {
			break
		}
		if *(*int)(cursor.object) == *(*int)(object) {
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
				l.deleteObject(object)
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
				l.deleteObject(object)
			}
			break
		} else if *(*int)(cursor.object) > *(*int)(object) {
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

func (l *NonBlockingList) insertObject(object unsafe.Pointer) bool {
	cursor := l.head
	for {
		if *(*int)(object) < *(*int)(cursor.markableNext.next.object) {
			currentNext := cursor.markableNext
			if currentNext.marked {
				continue
			}
			if *(*int)(cursor.object) == *(*int)(object) {
				return false
			}
			newNode := LockFreeNode{
				object: object,
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
				l.insertObject(object)
				return true
			}
			break
		}
		cursor = cursor.markableNext.next
	}
	return true
}

// Add is
func (l *NonBlockingList) Add(item int) bool {
	object := unsafe.Pointer(&item)
	return l.insertObject(object)
}

// Contains is
func (l *NonBlockingList) Contains(item int) bool {
	curr := l.head
	marked := false
	for *(*int)(curr.object) < item {
		marked = curr.markableNext.marked
		curr = curr.markableNext.next
	}
	return *(*int)(curr.object) == item && !marked
}

// Remove is
func (l *NonBlockingList) Remove(item int) bool {
	object := unsafe.Pointer(&item)
	return l.deleteObject(object)
}
