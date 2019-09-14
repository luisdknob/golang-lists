package list

// FineList is
type FineList struct {
	head *MuxNode
}

// Add is
func (l *FineList) Add(item int) bool {
	l.head.mux.Lock()
	pred := l.head
	curr := pred.next
	curr.mux.Lock()

	for curr.item < item {
		pred.mux.Unlock()
		pred = curr
		curr = curr.next
		curr.mux.Lock()
	}
	if item == curr.item {
		pred.mux.Unlock()
		curr.mux.Unlock()
		return false
	}

	node := &MuxNode{
		item: item,
	}
	node.next = curr
	pred.next = node

	pred.mux.Unlock()
	curr.mux.Unlock()

	return true
}

//Remove is
func (l *FineList) Remove(item int) bool {

	l.head.mux.Lock()
	pred := l.head
	curr := pred.next
	curr.mux.Lock()

	for curr.item < item {
		pred.mux.Unlock()
		pred = curr
		curr = curr.next
		curr.mux.Lock()
	}
	if item == curr.item {
		pred.next = curr.next
		pred.mux.Unlock()
		curr.mux.Unlock()
		return true
	}

	pred.mux.Unlock()
	curr.mux.Unlock()
	return false
}

// Contains is
func (l *FineList) Contains(item int) bool {
	l.head.mux.Lock()
	pred := l.head
	curr := pred.next
	curr.mux.Lock()

	for curr.item < item {
		pred.mux.Unlock()
		pred = curr
		curr = curr.next
		curr.mux.Lock()

	}
	if item == curr.item {
		pred.mux.Unlock()
		curr.mux.Unlock()
		return true
	}

	pred.mux.Unlock()
	curr.mux.Unlock()
	return false
}

// Count is
func (l *FineList) Count() (num int) {
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

// NewFineList is
func NewFineList() (l *FineList) {
	head := MuxNode{
		item: 0,
		next: &MuxNode{
			item: int(^uint(0) >> 1),
		},
	}
	l = &FineList{
		head: &head,
	}
	return
}
