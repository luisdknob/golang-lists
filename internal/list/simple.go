package list

// SimpleList definition with Mutex
type SimpleList struct {
	head *Node
}

// Add is
func (l *SimpleList) Add(item int) bool {
	var pred, curr *Node

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
func (l *SimpleList) Contains(item int) bool {
	var pred, curr *Node

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
func (l *SimpleList) Remove(item int) bool {
	var pred, curr *Node

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
func (l *SimpleList) Count() (num int) {
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

// NewSimpleList is
func NewSimpleList() (l *SimpleList) {
	head := Node{
		item: 0,
		next: &Node{
			item: int(^uint(0) >> 1),
		},
	}
	l = &SimpleList{
		head: &head,
	}
	return
}
