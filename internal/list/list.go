package list

import "sync"

// List is the list Definition
type List interface {
	Add(item int) bool
	Remove(item int) bool
	Contains(item int) bool
}

// Node is
type Node struct {
	item int
	next *Node
}

// MuxNode is
type MuxNode struct {
	item int
	next *MuxNode
	mux  sync.Mutex
}
