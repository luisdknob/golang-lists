package list

// OptmisticList definition with Mutex
type OptmisticList struct {
	head *MuxNode
}

func (l *OptmisticList) validate(pred *MuxNode, curr *MuxNode) bool {
	node := l.head
	for node.item <= pred.item {
		if node == pred {
			return pred.next == curr
		}
		node = node.next
	}
	return false
}

// Add is
func (l *OptmisticList) Add(item int) bool {
	var pred, curr *MuxNode

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
			node := &MuxNode{
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
func (l *OptmisticList) Contains(item int) bool {
	pred := l.head
	curr := pred.next

	for {
		for curr.item < item {
			pred = curr
			curr = curr.next
		}
		pred.mux.Lock()
		curr.mux.Lock()
		if l.validate(pred, curr) {
			pred.mux.Unlock()
			curr.mux.Unlock()
			return (curr.item == item)
		}
		pred.mux.Unlock()
		curr.mux.Unlock()
	}
}

// Remove is
func (l *OptmisticList) Remove(item int) bool {
	var pred, curr *MuxNode

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
			pred.next = curr.next
			pred.mux.Unlock()
			curr.mux.Unlock()
			return true
		}
		pred.mux.Unlock()
		curr.mux.Unlock()
	}
}

// NewOptimisticList is
func NewOptimisticList() (l *OptmisticList) {
	head := MuxNode{
		item: 0,
		next: &MuxNode{
			item: int(^uint(0) >> 1),
		},
	}
	l = &OptmisticList{
		head: &head,
	}
	return
}
