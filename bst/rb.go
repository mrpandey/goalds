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

// BST is implemented using Red-Black Tree.
// An RBTree has following properties
//  1. All nodes are either red or black.
//  2. Root node is black.
//  3. Leaf (nil) nodes are considered black.
//  4. A red node does not have a red child.
//  5. In any subtree, all simple paths from root of the subtree to leaves (nil nodes) contain the same number of black nodes.
//  6. Corollary: Color of a single child must be red. If it were black, then property 5 would be violated.
//     This means that a non-nil black node always has a non-nil sibling.
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

	// At this point all properties of red-black trees are satisfied, except parent may be also be red.
	rb.fixInsert(newNd)
}

// Returns true if there exists a node having the given value in the tree.
// Returns false otherwise.
func (rb *RBTree[T]) Exists(val T) bool {
	nd := rb.findNode(val)
	return nd != nil
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
// Working is similar to deletion of node in a normal BST. Only addition is the fixing part.
// TODO: Understand how the fuck fixing works.
func (rb *RBTree[T]) Delete(val T) error {
	nd := rb.findNode(val)
	if nd == nil {
		return ErrValueDoesNotExist
	}

	rb.len--

	ogColor := nd.clr
	var ndToFix *node[T] = nil

	if nd.left == nil {
		ndToFix = nd.right
		rb.replace(nd, ndToFix)
	} else if nd.right == nil {
		ndToFix = nd.left
		rb.replace(nd, ndToFix)
	} else {
		// substitute for nd
		sub := nd.right.getMin()
		ogColor = sub.clr
		ndToFix = sub.right

		if sub.parent != nd {
			// first replace substitute by its right child
			// this is easy since sub.left == nil
			rb.replace(sub, sub.right)

			// update right child of sub
			sub.right = nd.right
			sub.right.parent = sub
		}

		rb.replace(nd, sub)
		sub.left = nd.left
		if sub.left != nil {
			sub.left.parent = sub
		}
		sub.clr = nd.clr
	}

	if ogColor == black {
		rb.fixDelete(ndToFix)
	}

	return nil
}

// Returns non-nil pointer to the first node found with the given value.
func (rb *RBTree[T]) findNode(val T) *node[T] {
	nd := rb.root

	for nd != nil {
		if nd.value == val {
			return nd
		} else if nd.value < val {
			nd = nd.right
		} else {
			nd = nd.left
		}
	}

	return nil
}

// Find the node with minimum value in subtree rooted at nd.
func (nd *node[T]) getMin() *node[T] {
	for nd.left != nil {
		nd = nd.left
	}
	return nd
}

// Replace a node with its substitute in the tree without affecting their children.
// Substitute can be nil, but not the node.
func (rb *RBTree[T]) replace(nd, sub *node[T]) {
	p := nd.parent

	if p == nil {
		rb.root = sub
	} else if nd == p.left {
		p.left = sub
	} else {
		p.right = sub
	}

	if sub != nil {
		sub.parent = p
	}
}

// Left rotates the the node to balance the tree.
// The node will become the left child of it's current right child.
func (rb *RBTree[T]) rotateLeft(nd *node[T]) {
	if nd == nil || nd.right == nil {
		return
	}

	r := nd.right
	rb.replace(nd, r)

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

	l := nd.left
	rb.replace(nd, l)

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

// Copied brainlessly from CLRS.
// TODO: Understand how it works.
func (rb *RBTree[T]) fixDelete(nd *node[T]) {
	for nd != rb.root && nd.color() == black {
		// parent is non-nill since nd != root
		if nd == nd.parent.left {
			sib := nd.parent.right

			if sib.color() == red {
				// sib is non-nill since color is red
				sib.clr = black
				nd.parent.clr = red
				rb.rotateLeft(nd.parent)
				// sib will change after rotation
				sib = nd.parent.right
			}

			// what if sib == nil?? CLRS doesn't cover this so I will pretend to be blind.

			if sib.left.color() == black && sib.right.color() == black {
				sib.clr = red
				nd = nd.parent
			} else {
				if sib.right.color() == black {
					sib.left.clr = black
					sib.clr = red
					rb.rotateRight(sib)
					// sib will change after rotation
					sib = nd.parent.right
				}

				sib.clr = nd.parent.clr
				nd.parent.clr = black
				sib.right.clr = black
				rb.rotateLeft(nd.parent)
				nd = rb.root
			}
		} else {
			sib := nd.parent.left

			if sib.color() == red {
				// sib is non-nill since color is red
				sib.clr = black
				nd.parent.clr = red
				rb.rotateRight(nd.parent)
				// sib will change after rotation
				sib = nd.parent.left
			}

			// what if sib == nil?? CLRS doesn't cover this so I will pretend to be blind.

			if sib.left.color() == black && sib.right.color() == black {
				sib.clr = red
				nd = nd.parent
			} else {
				if sib.left.color() == black {
					sib.right.clr = black
					sib.clr = red
					rb.rotateRight(sib)
					// sib will change after rotation
					sib = nd.parent.left
				}

				sib.clr = nd.parent.clr
				nd.parent.clr = black
				sib.left.clr = black
				rb.rotateRight(nd.parent)
				nd = rb.root
			}
		}
	}

	rb.root.clr = black
}
