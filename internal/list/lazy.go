package list

// LazyList definition with Mutex
type LazyList struct {
	head *MarkedNode
}

func (l *LazyList) validate(pred, curr *MarkedNode) bool {
	return !pred.marked && !curr.marked && pred.next == curr
}

// Add is
func (l *LazyList) Add(item int) bool {
	var pred, curr *MarkedNode

	for {
		pred = l.head
		curr = pred.next
		for curr.item < item {
			pred = curr
			curr = curr.next
		}
		pred.mux.Lock()
		curr.mux.Lock()
		if l.validate(pred, curr) {
			if item == curr.item {
				pred.mux.Unlock()
				curr.mux.Unlock()
				return false
			}
			node := &MarkedNode{
				item: item,
			}
			node.next = curr
			pred.next = node
			pred.mux.Unlock()
			curr.mux.Unlock()
			return true
		}
		pred.mux.Unlock()
		curr.mux.Unlock()
	}
}

// Contains is
func (l *LazyList) Contains(item int) bool {
	for {
		curr := l.head

		for curr.item < item {
			curr = curr.next
		}

		return curr.item == item && !curr.marked
	}
}

// Remove is
func (l *LazyList) Remove(item int) bool {
	var pred, curr *MarkedNode

	for {
		pred = l.head
		curr = pred.next
		for curr.item < item {
			pred = curr
			curr = curr.next
		}
		pred.mux.Lock()
		curr.mux.Lock()
		if l.validate(pred, curr) {
			if item != curr.item {
				pred.mux.Unlock()
				curr.mux.Unlock()
				return false
			}
			curr.marked = true
			pred.next = curr.next
			pred.mux.Unlock()
			curr.mux.Unlock()
			return true
		}
		pred.mux.Unlock()
		curr.mux.Unlock()
	}
}

// NewLazyList is
func NewLazyList() (l *LazyList) {
	head := MarkedNode{
		item: 0,
		next: &MarkedNode{
			item: int(^uint(0) >> 1),
		},
	}
	l = &LazyList{
		head: &head,
	}
	return
}
