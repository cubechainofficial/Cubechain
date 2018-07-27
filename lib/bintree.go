package lib

import "fmt"

type Node struct {
	value interface{}
	left  *Node
	right *Node
}

type BinaryTree struct {
	root  *Node
	size  int 
	order int 
}

func (n *Node) Left() *Node {
	return n.left
}

func (n *Node) Right() *Node {
	return n.right
}

func (n *Node) Value() interface{} {
	return n.value
}

func (n *Node) SetLeft(left *Node) {
	n.left = left
}

func (n *Node) SetRight(right *Node) {
	n.right = right
}

func (n *Node) SetValue(val interface{}) {
	n.value = val
}

func (b *BinaryTree) Size() int {
	return b.size
}

func (b *BinaryTree) Order() int {
	return b.order
}

func (b *BinaryTree) uniqueValue(val interface{}) bool {
	if b.root == nil {
		return true
	}
	var q []*Node
	q = append(q, b.root)

	for len(q) != 0 {
		n := q[0]
		q = q[1:]

		if n.Value() == val {
			return false
		}

		if n.Left() != nil {
			q = append(q, n.Left())
		}
		if n.Right() != nil {
			q = append(q, n.Right())
		}
	}
	return true
}

func (b *BinaryTree) Insert(val interface{}) {
	switch {
	case b.root == nil:
		b.root = &Node{value: val}
		b.size++
		return
	case b.uniqueValue(val) == false:
		fmt.Println(val, "is already present in the Binary Tree")
	default:
		var q []*Node
		q = append(q, b.root)
		for len(q) != 0 {
			node := q[0]
			q = q[1:]
			if node.Left() == nil {
				node.SetLeft(&Node{value: val})
				b.size++
				return
			} else if node.Right() == nil {
				node.SetRight(&Node{value: val})
				b.size++
				return
			}
			if node.Left() != nil {
				q = append(q, node.Left())
			}
			if node.Right() != nil {
				q = append(q, node.Right())
			}
		}
	}
}

func (b *BinaryTree) Each(f func(val interface{})) {
	if b.root == nil {
		fmt.Println("Empty Tree")
	}
	var q []*Node
	q = append(q, b.root)
	for len(q) != 0 {
		n := q[0]
		q = q[1:]
		f(n.Value())
		if n.Left() != nil {
			q = append(q, n.Left())
		}
		if n.Right() != nil {
			q = append(q, n.Right())
		}
	}
}

func (b *BinaryTree) Exists(v interface{}) bool {
	if b.uniqueValue(v) {
		fmt.Printf("The value %v does not exists in the tree!\n", v)
		return false
	} else {
		fmt.Printf("The value %v exists in the tree!\n", v)
		return true
	}
}

func (b *BinaryTree) deleteDeepestRightMostNode(drn *Node) {
	var q []*Node
	q = append(q, b.root)

	for len(q) != 0 {
		n := q[0]
		q = q[1:]

		if n.Left() != nil {
			if n.Left() == drn {
				n.SetLeft(nil)
				return
			}
			q = append(q, n.Left())
		}
		if n.Right() != nil {
			if n.Right() == drn {
				n.SetRight(nil)
				return
			}
			q = append(q, n.Right())
		}
	}
}


func (b *BinaryTree) Delete(val interface{}) {
	if b.root == nil {
		fmt.Println("Empty Tree")
		return
	}
	var q []*Node
	var dtn *Node
	var n *Node
	q = append(q, b.root)

	for len(q) != 0 {
		n = q[0]
		q = q[1:]

		if n.Value() == val {
			dtn = n
		}
		if n.Left() != nil {
			q = append(q, n.Left())
		}
		if n.Right() != nil {
			q = append(q, n.Right())
		}
	}
	if dtn == nil {
		fmt.Println("The value does not exist in the Tree")
		return
	}
	b.deleteDeepestRightMostNode(n)
	dtn.SetValue(n.Value())
	b.size--
}

func (b *BinaryTree) Print() {
	if b.root == nil {
		fmt.Println("Empty Tree")
		return
	}
	var q []*Node
	q = append(q, b.root)
	for len(q) != 0 {
		n := q[0]
		q = q[1:]
		fmt.Println("Value:", n.Value())
		if n.Left() != nil {
			q = append(q, n.Left())
		}
		if n.Right() != nil {
			q = append(q, n.Right())
		}
	}
}

