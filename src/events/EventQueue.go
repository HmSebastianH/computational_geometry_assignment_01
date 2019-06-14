package events

// This package is a adaption of the btree from "github.com/ross-oreto/go-tree"
// It shares the same balacing features and basic access methods, but it uses a ccv ordering to insert values
// The nodes are also restricted to contain lines.

import (
	"fmt"
)

// EventQueue represents an AVL tree
type EventQueue struct {
	Root *Node
}

// Node represents a node in the tree with a value, left and right children, and a height/balance of the node.
type Node struct {
	Value       SweepEvent
	left, right *Node
	height      int8
}

// New returns a new btree
func NewEventQueue() *EventQueue { return new(EventQueue).Init() }

// Init initializes all values/clears the tree and returns the tree pointer
func (t *EventQueue) Init() *EventQueue {
	t.Root = nil
	return t
}

func (t *EventQueue) balance() int8 {
	if t.Root != nil {
		return balance(t.Root)
	}
	return 0
}

func (t *EventQueue) PrintOut() {
	fmt.Print("Events: [")
	t.Root.printOut()
	fmt.Println("]")

}

func (n *Node) printOut() {
	if n == nil {
		return
	}
	n.left.printOut()
	fmt.Print(n.Value, " ")
	n.right.printOut()
}

func (t *EventQueue) AssertOrder() bool {
	return t.Root.assertOrder()
}

func (n *Node) assertOrder() bool {
	if n == nil {
		return true
	}
	n.left.assertOrder()
	n.right.assertOrder()
	if n.left != nil && n.Value.CompareTo(n.left.Value) != 1 {
		return false
	}
	if n.right != nil && n.Value.CompareTo(n.right.Value) != -1 {
		return false
	}
	return true
}

// Insert inserts a new value into the tree and returns the tree pointer
func (t *EventQueue) Insert(value SweepEvent) *EventQueue {
	added := false
	t.Root = insert(t.Root, value, &added)
	t.AssertOrder()
	return t
}

func insert(n *Node, value SweepEvent, added *bool) *Node {
	if n == nil {
		// If this is a empty leaf insert the line here
		*added = true
		return (&Node{Value: value}).Init()
	}
	comp := value.CompareTo(n.Value)
	if value.CompareTo(n.Value)*(-1) != n.Value.CompareTo(value) {
		panic("!!")
		//aRes := value.CompareTo(n.Value)
		//bRes := n.Value.CompareTo(value)
		//_, _ = aRes, bRes
	}
	if comp > 0 {
		n.right = insert(n.right, value, added)
	} else if comp < 0 {
		// Points with overlap or to the left of the line are inserted to its left
		n.left = insert(n.left, value, added)
	} else {
		// TODO: This is because intersections can be detected multiple times
		fmt.Println("Warning: Duplicate Event", value.GetX(), value.String()) // This is not a warning, we could just replace it
	}

	n.height = n.maxHeight() + 1
	/*
	// TODO: Replace this with a non tree datastructure
	currentBalance := balance(n)

	if currentBalance > 1 {
		comp := value.CompareTo(n.left.Value)
		if comp < 0 {
			return n.rotateRight()
		} else if comp > 0 {
			n.left = n.left.rotateLeft()
			return n.rotateRight()
		}
	} else if currentBalance < -1 {
		comp = value.CompareTo(n.right.Value)
		if comp > 0 {
			return n.rotateLeft()
		} else {
			n.right = n.right.rotateRight()
			return n.rotateLeft()
		}
	}*/

	return n
}

// Head returns the first value in the tree
func (t *EventQueue) Head() *Node {
	if t.Root == nil {
		return nil
	}
	n := t.Root
	for n.left != nil {
		n = n.left
	}
	return n
}

func (t *EventQueue) Pop() SweepEvent {
	if t.Root == nil {
		return nil
	}
	if t.Root.left == nil {
		// Pop the Root itself
		headNode := t.Root
		t.Root = headNode.right
		return headNode.Value
	}
	// Find "first" event
	parentNode := t.Root
	headNode := t.Root.left
	for headNode.left != nil {
		parentNode = headNode
		headNode = headNode.left
	}
	parentNode.left = headNode.right

	t.balance() // TODO: test performance without balancing
	// Note: A tree is probably not the bestt structure for almost always accessing the first member
	return headNode.Value
}

// Init initializes the values of the node or clears the node and returns the node pointer
func (n *Node) Init() *Node {
	n.height = 1
	n.left = nil
	n.right = nil
	return n
}

// String returns a string representing the node
func (n *Node) String() string {
	return fmt.Sprint(n.Value)
}

// Debug prints out useful debug information about the tree node for debugging purposes
func (n *Node) Debug() {
	var children string
	if n.left == nil && n.right == nil {
		children = "no children |"
	} else if n.left != nil && n.right != nil {
		children = fmt.Sprint("left child:", n.left.String(), " right child:", n.right.String())
	} else if n.right != nil {
		children = fmt.Sprint("right child:", n.right.String())
	} else {
		children = fmt.Sprint("left child:", n.left.String())
	}

	fmt.Println(n.String(), "|", "height", n.height, "|", "balance", balance(n), "|", children)
}

func height(n *Node) int8 {
	if n != nil {
		return n.height
	}
	return 0
}

func balance(n *Node) int8 {
	if n == nil {
		return 0
	}
	return height(n.left) - height(n.right)
}

func (n *Node) rotateRight() *Node {
	l := n.left
	// Rotation
	l.right, n.left = n, l.right

	// update heights
	n.height = n.maxHeight() + 1
	l.height = l.maxHeight() + 1

	return l
}

func (n *Node) rotateLeft() *Node {
	r := n.right
	// Rotation
	r.left, n.right = n, r.left

	// update heights
	n.height = n.maxHeight() + 1
	r.height = r.maxHeight() + 1

	return r
}

func (n *Node) min() *Node {
	current := n
	for current.left != nil {
		current = current.left
	}
	return current
}

func (n *Node) maxHeight() int8 {
	rh := height(n.right)
	lh := height(n.left)
	if rh > lh {
		return rh
	}
	return lh
}
