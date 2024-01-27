package bst

import (
	"cmp"
)

type color int

const (
	black color = 0
	red   color = 1
)

type node[T cmp.Ordered] struct {
	left   *node[T]
	right  *node[T]
	parent *node[T]
	clr    color
	value  T
}

func (nd *node[T]) color() color {
	if nd == nil {
		// nil pointers are considered black
		return black
	}
	return nd.clr
}

// BST is implemented using Red-Black Tree
type RBTree[T cmp.Ordered] struct {
	root *node[T]
	len  int
}

func NewRBTree[T cmp.Ordered]() *RBTree[T] {
	return &RBTree[T]{}
}

func (rb *RBTree[T]) Len() int {
	return rb.len
}

// Insert a new node in the tree with the given value. Inserts even if the value already exists.
func (rb *RBTree[T]) Insert(val T) {
	// insert new node as usual

	nd := rb.root
	var p *node[T] = nil

	for nd != nil {
		p = nd
		if val <= nd.value {
			nd = nd.left
		} else {
			nd = nd.right
		}
	}

	rb.len++
	newNd := &node[T]{
		value: val,
	}

	if p == nil {
		newNd.clr = black
		rb.root = newNd
		return
	}

	newNd.clr = red
	newNd.parent = p

	if val <= p.value {
		p.left = newNd
	} else {
		p.right = newNd
	}

	// At this point all properties of red-black trees are satisfied, except p may be also be red.
	rb.fixInsert(newNd)
}

// Returns true if there exists a node having the given value in the tree.
// Returns false otherwise.
func (rb *RBTree[T]) Exists(val T) bool {
	nd := rb.root

	for nd != nil {
		if val == nd.value {
			return true
		} else if val > nd.value {
			nd = nd.right
		} else {
			nd = nd.left
		}
	}

	return false
}

// Returns the values of nodes in ascending order.
func (rb *RBTree[T]) GetValues() []T {
	values := make([]T, rb.Len())
	i := 0

	// morris inorder traversal
	nd := rb.root

	for nd != nil {
		if nd.left != nil {
			// find the node before nd in inorder traversal
			pre := nd.left
			for pre.right != nil && pre.right != nd {
				pre = pre.right
			}

			if pre.right == nil {
				// add backlink the ancestor nd
				pre.right = nd
				// advance to left subtree
				nd = nd.left
			} else {
				// pre.right == nd
				// backlink already exists
				// this means entire left subtree of nd has been visited
				// so just visit nd and advance to right subtree
				values[i] = nd.value
				i++
				nd = nd.right
				// remove the backlink
				pre.right = nil
			}
		} else {
			// visit the nd and advance to right subtree
			values[i] = nd.value
			i++
			nd = nd.right
		}
	}

	return values
}

// Deletes a node in the tree with the given value.
// If there are multiple such nodes, any one of them might be deleted.
// Non-nill error is returned if no such node is found. Otherwise, nil is returned.
func (rb *RBTree[T]) Delete(val T) error {
	return nil
}

// Left rotates the the node to balance the tree.
// The node will become the left child of it's current right child.
func (rb *RBTree[T]) rotateLeft(nd *node[T]) {
	if nd == nil || nd.right == nil {
		return
	}

	p := nd.parent
	r := nd.right

	r.parent = p
	if p == nil {
		rb.root = r
		rb.root.clr = black
	} else if nd == p.left {
		p.left = r
	} else {
		p.right = r
	}

	nd.parent = r
	nd.right = r.left
	if nd.right != nil {
		nd.right.parent = nd
	}
	r.left = nd
}

// Right rotates the the node to balance the tree.
// The node will become the right child of it's current left child.
func (rb *RBTree[T]) rotateRight(nd *node[T]) {
	if nd == nil || nd.left == nil {
		return
	}

	p := nd.parent
	l := nd.left

	l.parent = p
	if p == nil {
		rb.root = l
		rb.root.clr = black
	} else if nd == p.left {
		p.left = l
	} else {
		p.right = l
	}

	nd.parent = l
	nd.left = l.right
	if nd.left != nil {
		nd.left.parent = nd
	}
	l.right = nd
}

// Newly inserted non-root nodes are red by default.
// If parent of this new node is also red, then we need to fix this.
// A red node should have both its children black.
func (rb *RBTree[T]) fixInsert(nd *node[T]) {
	if nd.color() != red {
		return
	}

	for nd.parent.color() == red {
		p := nd.parent

		if p == p.parent.left {
			// If p is red, it cannot be the root node.
			// So p has non-nil parent, and it also has a sibling (might be nil).
			psib := p.parent.right

			if psib.color() == red {
				// Since it is red, psib cannot be nil.
				psib.clr = black
				p.clr = black
				p.parent.clr = red

				nd = p.parent
			} else {
				if nd == p.right {
					nd = p
					rb.rotateLeft(nd)
				}

				// At this point, nd and its parent are red, but parent's sibling is black.
				// This implies that parent's parent is also black.

				nd.parent.clr = black
				nd.parent.parent.clr = red

				rb.rotateRight(nd.parent.parent)
			}
		} else {
			// If p is red, it cannot be the root node.
			// So p has non-nil parent, and it also has a sibling (might be nil).
			psib := p.parent.left

			if psib.color() == red {
				// Since it is red, psib cannot be nil.
				psib.clr = black
				p.clr = black
				p.parent.clr = red

				nd = p.parent
			} else {
				if nd == p.left {
					nd = p
					rb.rotateRight(nd)
				}

				// At this point, nd and its parent are red, but parent's sibling is black.
				// This implies that parent's parent is also black.

				nd.parent.clr = black
				nd.parent.parent.clr = red

				rb.rotateLeft(nd.parent.parent)
			}
		}
	}

	// If nd is not nil, root is non-nill.
	rb.root.clr = black
}
