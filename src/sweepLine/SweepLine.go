package sweepLine

// This package is a adaption of the btree from "github.com/ross-oreto/go-tree"
// It shares the same balacing features and basic access methods, but it uses a ccv ordering to insert values
// The nodes are also restricted to contain lines.

import (
"fmt"
	. "geometry"
)

// Btree represents an AVL tree
type SweepLine struct {
	Root *Node
	len  int
}


// Node represents a node in the tree with a value, left and right children, and a height/balance of the node.
type Node struct {
	Value       Line
	left, right, parent *Node
	height      int8
}

// New returns a new btree
func NewSweepLine() *SweepLine { return new(SweepLine).Init() }

// Init initializes all values/clears the tree and returns the tree pointer
func (t *SweepLine) Init() *SweepLine {
	t.Root = nil
	t.len = 0
	return t
}

// Empty returns true if the tree is empty
func (t *SweepLine) Empty() bool {
	return t.Root == nil
}

// NotEmpty returns true if the tree is not empty
func (t *SweepLine) NotEmpty() bool {
	return t.Root != nil
}

func (t *SweepLine) balance() int8 {
	if t.Root != nil {
		return balance(t.Root)
	}
	return 0
}

// Insert inserts a new value into the tree and returns the tree pointer
func (t *SweepLine) Insert(value Line) *Node {
	if value.Start.X != value.End.X {
		panic("Vertical Lines / Points are not supported by the Sweep line")
	}
	var insertedNode Node
	t.Root = t.Root.insert(nil, value, &insertedNode)
	if insertedNode.parent != nil {
		t.len++
	}
	return &insertedNode
}

func (n *Node) insert(parent *Node, value Line, insertedNode *Node) *Node {
	if n == nil {
		// If this is a empty leaf insert the line here
		newNode := Node{Value: value}
		newNode.Init().setParent(parent)
		insertedNode = &newNode
		return (&Node{Value: value}).Init().setParent(parent)
	}
	ccw := Ccw(n.Value, value.Start)
	if ccw > 0 {
		n.right = n.right.insert(n, value, insertedNode)
	} else {
		// Points with overlap or to the left of the line are inserted to its left
		n.left = n.left.insert(n, value, insertedNode)
	}

	n.height = n.maxHeight() + 1
	currentBalance := balance(n)

	if currentBalance > 1 {
		ccw = Ccw(n.left.Value, value.Start)
		if ccw <= 0 {
			return n.rotateRight()
		} else {
			n.left = n.left.rotateLeft()
			return n.rotateRight()
		}
	} else if currentBalance < -1 {
		ccw = Ccw(n.right.Value, value.Start)
		if ccw > 0 {
			return n.rotateLeft()
		} else {
			n.right = n.right.rotateRight()
			return n.rotateLeft()
		}
	}
	return n
}

// Len return the number of nodes in the tree
func (t *SweepLine) Len() int {
	return t.len
}

// Head returns the first value in the tree
func (t *SweepLine) Head() *Line {
	if t.Root == nil {
		return nil
	}
	var beginning = t.Root
	for beginning.left != nil {
		beginning = beginning.left
	}
	if beginning == nil {
		for beginning.right != nil {
			beginning = beginning.right
		}
	}
	if beginning != nil {
		return &beginning.Value
	}
	return nil
}

// Tail returns the last value in the tree
func (t *SweepLine) Tail() *Line {
	if t.Root == nil {
		return nil
	}
	var beginning = t.Root
	for beginning.right != nil {
		beginning = beginning.right
	}
	if beginning == nil {
		for beginning.left != nil {
			beginning = beginning.left
		}
	}
	if beginning != nil {
		return &beginning.Value
	}
	return nil
}

// Len return the number of nodes in the tree
func (n *Node) Left() *Node {
	if n.left != nil {
		return n.left.max()
	}

	// I am a left node with no children, search a parent which is left of me
	currentParent := n.parent
	currentParrentChild := n
	for currentParent != nil {
		if currentParrentChild == currentParent.right {
			// We found a path where the tree we came from is on the right, therfore the node is to the left
			return currentParent
		}
		currentParrentChild = currentParent
		currentParent = currentParent.parent
	}

	// If nothing matched for now, this is the last node in the tree
	return nil
}

// Len return the number of nodes in the tree
func (n *Node) Right() *Node {
	if n.right != nil {
		return n.right.min()
	}

	// I am a left node with no children, search a parent which is left of me
	currentParent := n.parent
	currentParrentChild := n
	for currentParent != nil {
		if currentParrentChild == currentParent.left {
			// We found a path where the tree we came from is on the left, therfore the node is to the right
			return currentParent
		}
		currentParrentChild = currentParent
		currentParent = currentParent.parent
	}

	// If nothing matched for now, this is the last node in the tree
	return nil
}

func deleteSweepNode(n *Node, value Line, deleted *bool) *Node {
	if n == nil {
		return n
	}

	c := 1

	if c < 0 {
		n.left = deleteSweepNode(n.left, value, deleted)
	} else if c > 0 {
		n.right = deleteSweepNode(n.right, value, deleted)
	} else {
		if n.left == nil {
			// Replace myself with my right node
			t := n.right
			t.parent = n.parent
			n.Init()
			return t
		} else if n.right == nil {
			// Replace myself with my left node
			t := n.left
			t.parent = n.parent
			n.Init()
			return t
		}
		t := n.right.min()
		n.Value = t.Value
		n.right = deleteSweepNode(n.right, t.Value, deleted)
		*deleted = true
	}

	//re-balance
	if n == nil {
		return n
	}
	n.height = n.maxHeight() + 1
	bal := balance(n)
	if bal > 1 {
		if balance(n.left) >= 0 {
			return n.rotateRight()
		}
		n.left = n.left.rotateLeft()
		return n.rotateRight()
	} else if bal < -1 {
		if balance(n.right) <= 0 {
			return n.rotateLeft()
		}
		n.right = n.right.rotateRight()
		return n.rotateLeft()
	}

	return n
}


// Init initializes the values of the node or clears the node and returns the node pointer
func (n *Node) Init() *Node {
	n.height = 1
	n.left = nil
	n.right = nil
	return n
}

func (n *Node) setParent(parent *Node) *Node {
	n.parent = parent
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

	l.parent = n.parent
	n.parent = l

	// update heights
	n.height = n.maxHeight() + 1
	l.height = l.maxHeight() + 1

	return l
}

func (n *Node) rotateLeft() *Node {
	r := n.right
	// Rotation
	r.left, n.right = n, r.left

	r.parent = n.parent
	n.parent = r

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

func (n *Node) max() *Node {
	current := n
	for current.right != nil {
		current = current.right
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