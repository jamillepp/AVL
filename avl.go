package avl

import (
	"fmt"
)

type node struct {
	value   int
	child   [2]*node
	height  [2]int
	balance int
	parent  *node
}

func NewAVL() node {
	return node{
		value: -1,
	}
}

func (tree *node) Search(v int) *node {
	if tree.value == -1 {
		fmt.Println("tree its <nil>")
		return nil
	}

	if tree.value == v {
		fmt.Println("The value was found")
		return tree
	} else {
		side := compare(v, tree.value)
		if tree.child[side] != nil {
			return (tree.child[side]).Search(v)
		} else {
			fmt.Println("The value was not found")
			return nil
		}
	}
}

func (tree *node) Insert(v int) {
	if tree.value == -1 {
		tree.value = v
		fmt.Printf("Value %v its now the root\n", v)
	} else {
		side := compare(v, tree.value)
		if tree.child[side] == nil {
			tree.child[side] = &node{
				value:  v,
				parent: tree,
			}
			if !tree.hasTwoChilds() {
				tree.setHeight(side, "I")
			} else {
				tree.height[side] += 1
			}
			fmt.Printf("Value %v inserted\n", v)
			tree.verifyBalance()
		} else {
			(tree.child[side]).Insert(v)
		}
	}
}

func compare(value, nodeValue int) int {
	if value < nodeValue {
		return 0
	} else if value > nodeValue {
		return 1
	} else {
		fmt.Println("Insert of a existing value")
		return -1
	}
}

func (tree *node) setHeight(side int, op string) {
	if op == "I" {
		tree.height[side] += 1
	} else {
		tree.height[side] -= 1
	}
	if tree.parent != nil {
		side = compare(tree.value, tree.parent.value)
		tree.parent.setHeight(side, op)
	}
}

func (tree *node) verifyBalance() {
	tree.balance = tree.height[1] - tree.height[0]
	if tree.balance > 1 || tree.balance < -1 {
		if tree.parent != nil {
			side := compare(tree.value, tree.parent.value)
			tree.parent.setHeight(side, "D")
		}
		fmt.Printf("Value %v its unbalanced, %v\n", tree.value, tree.balance)
		if tree.balance < -1 && tree.child[0].balance < 0 {
			fmt.Printf("Simple rotation [L-L]\n")
			tree.rightRotation()
			tree.child[1].verifyBalance()
		} else if tree.balance > 1 && tree.child[1].balance > 0 {
			fmt.Printf("Simple rotation [R-R]\n")
			tree.leftRotation()
			tree.child[0].verifyBalance()
		} else if tree.balance < -1 && tree.child[0].balance > 0 {
			fmt.Printf("Double rotation [L-R]\n")
			tree.child[0].leftRotation()
			tree.rightRotation()
			tree.child[0].verifyBalance()
			tree.child[1].verifyBalance()
		} else if tree.balance > 1 && tree.child[1].balance < 0 {
			fmt.Printf("Double rotation [R-L]\n")
			tree.child[1].rightRotation()
			tree.leftRotation()
			tree.child[0].verifyBalance()
			tree.child[1].verifyBalance()
		}
	} else if tree.parent != nil {
		tree.parent.verifyBalance()
	}
}

func (tree *node) rightRotation() {
	z := *tree
	*tree = *(tree.child[0])
	tree.parent = z.parent
	z.parent = tree
	z.child[0] = tree.child[1]
	z.height[0] = tree.height[1]
	tree.child[1] = &z
	if z.child[0] != nil {
		z.child[0].parent = &z
	}
	if z.child[1] != nil {
		z.child[1].parent = &z
	}
	if tree.child[1].height[0] > tree.child[1].height[1] {
		tree.height[1] = tree.child[1].height[0] + 1
	} else {
		tree.height[1] = tree.child[1].height[1] + 1
	}
}

func (tree *node) leftRotation() {
	z := *tree
	*tree = *(tree.child[1])
	tree.parent = z.parent
	z.parent = tree
	z.child[1] = tree.child[0]
	z.height[1] = tree.height[0]
	tree.child[0] = &z
	if z.child[0] != nil {
		z.child[0].parent = &z
	}
	if z.child[1] != nil {
		z.child[1].parent = &z
	}
	if tree.child[0].height[0] > tree.child[0].height[1] {
		tree.height[0] = tree.child[0].height[0] + 1
	} else {
		tree.height[0] = tree.child[0].height[1] + 1
	}
}

func (tree *node) Delete(v int) bool {
	if tree.value == -1 {
		fmt.Println("tree its <nil>")
		return false
	}

	// Se é a raiz
	if tree.value == v {
		if tree.itsaLeaf() {
			tree.value = -1
			tree.child = [2]*node{nil, nil}
			tree.parent = nil
		} else if tree.hasOneChild() {
			if tree.child[0] != nil {
				*tree = *tree.child[0]
			} else {
				*tree = *tree.child[1]
			}
			tree.parent = nil
		} else if tree.hasTwoChilds() {
			successor := tree.inorderSuccessor(nil)
			successor.child[0] = tree.child[0]
			successor.child[1] = tree.child[1]
			successor.parent = nil
			*tree = *successor
		}
		fmt.Printf("Value %v deleted\n", v)
		tree.verifyBalance()
		return true
	} else {
		if tree.itsaLeaf() {
			fmt.Println("The value was not found")
			return false
		}

		side := compare(v, tree.value)
		treeChild := tree.child[side]
		if treeChild == nil {
			fmt.Println("The value was not found")
			return false
		}

		if treeChild.value == v && treeChild.itsaLeaf() {
			tree.child[side] = nil
			tree.setHeight(side, "D")
			fmt.Printf("Value %v deleted\n", v)
			tree.verifyBalance()
			return true
		} else if treeChild.value == v && treeChild.hasOneChild() {
			if treeChild.child[0] != nil {
				treeChild.child[0].parent = treeChild.parent
				*treeChild = *treeChild.child[0]
			}
			if treeChild.child[1] != nil {
				treeChild.child[1].parent = treeChild.parent
				*treeChild = *treeChild.child[1]
			}
			tree.setHeight(side, "D")
			fmt.Printf("Value %v deleted\n", v)
			tree.verifyBalance()
			return true
		} else if treeChild.value == v && treeChild.hasTwoChilds() {
			successor := tree.inorderSuccessor(nil)
			successor.child[0] = treeChild.child[0]
			successor.child[1] = treeChild.child[1]
			successor.height[0] = treeChild.height[0]
			successor.height[1] = treeChild.height[1]
			successor.parent = treeChild.parent
			*treeChild = *successor
			tree.verifyBalance()
			fmt.Printf("Value %v deleted\n", v)
			return true
		} else if treeChild.value != v {
			return treeChild.Delete(v)
		}
	}
	fmt.Println("Something went wrong")
	return false
}

func (tree *node) itsaLeaf() bool {
	if tree.child[0] == nil && tree.child[1] == nil {
		return true
	} else {
		return false
	}
}

func (tree *node) hasOneChild() bool {
	if tree.child[0] != nil && tree.child[1] == nil {
		return true
	} else if tree.child[0] == nil && tree.child[1] != nil {
		return true
	} else {
		return false
	}
}

func (tree *node) hasTwoChilds() bool {
	if tree.child[0] != nil && tree.child[1] != nil {
		return true
	} else {
		return false
	}
}

func (tree *node) inorderSuccessor(successor *node) *node {
	if successor == nil {
		successor = tree
		if tree.child[1].itsaLeaf() {
			successor = tree.child[1]
			tree.child[1] = nil
			return successor
		}
		return tree.child[1].inorderSuccessor(successor)
	} else if tree.itsaLeaf() {
		successor = tree
		return successor
	} else {
		if tree.child[0].itsaLeaf() {
			if !tree.hasTwoChilds() {
				tree.setHeight(0, "D")
			} else {
				tree.height[0] -= 1
			}
			successor = tree.child[0]
			tree.child[0] = nil
			return successor
		}
		return tree.child[0].inorderSuccessor(successor)
	}
}

func (tree *node) Print() {
	fmt.Printf("{%v ", tree.value)
	if tree.child[0] != nil {
		fmt.Printf("[%v, ", tree.child[0].value)
		defer (tree.child[0]).Print()
	} else {
		fmt.Printf("[<>, ")
	}
	if tree.child[1] != nil {
		fmt.Printf("%v]", tree.child[1].value)
		defer (tree.child[1]).Print()
	} else {
		fmt.Printf("<>]")
	}
	fmt.Printf(" %v [%v %v], ", tree.balance, tree.height[0], tree.height[1])
	if tree.parent != nil {
		fmt.Printf("%v}\n", tree.parent.value)
	} else {
		fmt.Printf("<>}\n")
	}
}
