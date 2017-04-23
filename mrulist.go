// Copyright 2017 by Michael Haeuptle
// Free to use for any purpose
//
// This package implements a most recently used double-linked list that can be used
// to create a bounded cache that evicts the least recently used entry.
// The mru list implements a simple strategy that moves an entry to the top (beginning) of the list
// with each access thus making it unlikely for eviction which is done from the bottom of the list.


package mrulist

import "fmt"

type Node struct {
	Prev *Node
	Next *Node
	Data interface{}
}

type MruList struct {
	size int   // max size
	count int  // current size
	root *Node
	last *Node
}

func NewMruList(size int) *MruList {
	ml := new(MruList)
	ml.size = size

	return ml
}

func (ml *MruList) GetRoot() *Node {
	return ml.root
}

func (ml *MruList) GetLast() *Node {
	return ml.last
}

// Add a data item and return a new node and optionally the old node that got removed
// The old node can be used to evict (remove) the corresponding entry from a map
func (ml *MruList) Add(data interface{}) (*Node, *Node) {

	var removeNode *Node

	if ml.count == ml.size {
		removeNode = ml.RemoveLast()
	}

	newNode := new(Node)
	newNode.Data = data
	ml.count++

	if ml.root == nil { // first node
		ml.root = newNode
		ml.last = newNode
		return newNode, removeNode
	}

	// append to end of list
	newNode.Prev = ml.last
	newNode.Next = nil

	ml.last.Next = newNode
	ml.last = newNode

	return newNode, removeNode
}

func (ml *MruList) RemoveLast() *Node {
	if ml.root == nil {
		return nil
	}

	ml.count--
	n := ml.last

	// only one element
	if ml.root == ml.last {
		ml.root = nil
		ml.last = nil
		return n
	}

	ml.last.Prev.Next = nil
	ml.last = ml.last.Prev

	return n
}

func (ml *MruList) Dump() {
	n := ml.root

	for n != nil {
		fmt.Printf("%p", n)
		fmt.Println(": ",n)
		n = n.Next
	}
	fmt.Println("----")
}

func (ml *MruList) DumpR() {
	n := ml.last

	for n != nil {
		fmt.Println(n)
		n = n.Prev
	}
	fmt.Println("----")
}

// Moves Node n up one position, closer to the beginning of the list
func (ml *MruList) MoveUp(n *Node) {
	// first element case
	if n == nil || n == ml.root {
		return
	}

	nPrev := n.Prev
	nPrevPrev := n.Prev.Prev
	nNext := n.Next

	n.Prev = nPrevPrev
	n.Next = nPrev
	if nPrevPrev != nil {
		nPrevPrev.Next = n
	}
	nPrev.Next = nNext
	nPrev.Prev = n
	if nNext != nil {
		nNext.Prev = nPrev
	}

	if nNext == nil {
		ml.last = nPrev
	}

	if nPrevPrev == nil {
		ml.root = n
	}
}

func (ml *MruList) Get(data interface{}) *Node {
	n := ml.root

	for n != nil {
		if n.Data == data {
			ml.MoveUp(n)
			return n
		}
		n = n.Next
	}

	return nil
}


