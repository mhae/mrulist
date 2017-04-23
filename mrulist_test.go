package mrulist_test

import (
	"testing"
	"mrulist"
	"fmt"
)

func TestAdd(t *testing.T) {
	ml := mrulist.NewMruList(5)
	for i:=0; i<10; i++ {
		ml.Add(i);
	}

	if ml.GetRoot().Data != 0 {
		t.Error("First element not 0")
	}

	if ml.GetLast().Data != 9 {
		t.Error("Last element not 9 (got: ", ml.GetLast().Data,")")
	}
}

const (
	FORWARD = 1
	BACKWARD = 2
)

func checkTraversal(t *testing.T, ml *mrulist.MruList, expected []int, direction int) {
	i := 0
	var n *mrulist.Node
	if direction == FORWARD {
		n = ml.GetRoot()
	} else {
		n = ml.GetLast()
	}
	for n != nil {
		if n.Data != expected[i] {
			t.Error("Traversal, expected", expected[i], "found: ", n.Data)
			break
		}
		i++
		if direction == FORWARD {
			n = n.Next
		} else {
			n = n.Prev
		}

	}
	if i != len(expected) {
		t.Error("Traversal error")
	}
}

func TestMoveUpWithTwoElements(t *testing.T) {
	ml := mrulist.NewMruList(5)

	ml.Add(1)
	ml.Add(2)
	ml.Get(1)
	if ml.GetRoot().Data != 1 {
		t.Error("First element is not 1")
	}
	if ml.GetLast().Data != 2 {
		t.Error("Last element is not 2")
	}

	ml.Get(2)
	if ml.GetRoot().Data != 2 {
		t.Error("First element is not 2")
	}
	if ml.GetLast().Data != 1 {
		t.Error("Last element is not 1")
	}

	// check traversals
	checkTraversal(t, ml, []int { 2, 1}, FORWARD);
	checkTraversal(t, ml, []int { 1, 2}, BACKWARD);
}

func TestMoveUpWithThreeElements(t *testing.T) {
	ml := mrulist.NewMruList(5)

	ml.Add(1)
	ml.Add(2)
	ml.Add(3)

	// list is now { 1, 2, 3 }

	// move up 3 one position
	ml.Get(3)
	checkTraversal(t, ml, []int { 1, 3, 2}, FORWARD);
	checkTraversal(t, ml, []int { 2, 3, 1}, BACKWARD);

	// move up 3 one position
	ml.Get(3)
	checkTraversal(t, ml, []int { 3, 1, 2}, FORWARD);
	checkTraversal(t, ml, []int { 2, 1, 3}, BACKWARD);

	// move up 2 one position
	ml.Get(2)
	checkTraversal(t, ml, []int { 3, 2, 1}, FORWARD);
	checkTraversal(t, ml, []int { 1, 2, 3}, BACKWARD);

	// move up 2 one position
	ml.Get(2)
	checkTraversal(t, ml, []int { 2, 3, 1}, FORWARD);
	checkTraversal(t, ml, []int { 1, 3, 2}, BACKWARD);

	// move up 3 one position
	ml.Get(3)
	checkTraversal(t, ml, []int { 3, 2, 1}, FORWARD);
	checkTraversal(t, ml, []int { 1, 2, 3}, BACKWARD);
}


// Example of a bounded MRU cache
type entry struct { // the cache entry
	value string // the cache value
	mruNode *mrulist.Node // ref to the mru node
}
type cache struct {
	mru *mrulist.MruList // the mru list
	data map[string]*entry
}

func (c *cache) add(k string, v string) {
	newNode, removedNode := c.mru.Add(k)
	if removedNode != nil { // need to remove the element that was returned since it is the least recently used
		delete(c.data, removedNode.Data.(string))
	}

	// create a new entry and add; keep track for the new mru node which is we need to update during access
	entry := new(entry)
	entry.value = v
	entry.mruNode = newNode
	c.data[k] = entry
}

func (c *cache) get(k string) (string, bool) {
	entry, ok := c.data[k]

	if !ok {
		return entry.value, false
	}

	c.mru.MoveUp(entry.mruNode); // increase this entry's usage by moving it up one position
	return entry.value, ok
}

func (c *cache) dump() {
	for key, value := range c.data {
		fmt.Println("Key:", key, "Value:", value)
	}
}

func TestCache(t *testing.T) {
	var myCache cache
	myCache.mru = mrulist.NewMruList(4)
	myCache.data = make(map[string]*entry)

	myCache.add("1", "one")
	myCache.add("2", "two")
	myCache.add("3", "three")
	myCache.add("4", "four")
	myCache.add("5", "five")
	var s, ok = myCache.get("5")
	if !ok || s != "five" {
		t.Error("Expected value 'five'")
	}
	myCache.add("6", "six")
	if len(myCache.data) != 4 {
		t.Error("Expected cache size is 4")
	}
}