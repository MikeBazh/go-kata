package tree

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Node struct {
	Key    int
	Height int
	Left   *Node
	Right  *Node
}

type AVLTree struct {
	Root *Node
}

func NewNode(key int) *Node {
	return &Node{Key: key, Height: 1}
}

func (t *AVLTree) Insert(key int) {
	t.Root = insert(t.Root, key)
}

func (t *AVLTree) ToMermaid() string {
	var builder strings.Builder

	builder.WriteString("graph TD\n")

	if t.Root != nil {
		generateMermaid(t.Root, &builder)
	}

	return builder.String()
}

func generateMermaid(node *Node, builder *strings.Builder) {
	if node != nil {
		builder.WriteString(fmt.Sprintf("  %d\n", node.Key))
		if node.Left != nil {
			builder.WriteString(fmt.Sprintf("  %d --> %d\n", node.Key, node.Left.Key))
			generateMermaid(node.Left, builder)
		}
		if node.Right != nil {
			builder.WriteString(fmt.Sprintf("  %d --> %d\n", node.Key, node.Right.Key))
			generateMermaid(node.Right, builder)
		}
	}
}

func height(node *Node) int {
	if node == nil {
		return 0
	}
	return node.Height
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func updateHeight(node *Node) {
	node.Height = 1 + max(height(node.Left), height(node.Right))
}

func getBalance(node *Node) int {
	if node == nil {
		return 0
	}
	return height(node.Left) - height(node.Right)
}

func leftRotate(x *Node) *Node {
	y := x.Right
	x.Right = y.Left
	y.Left = x

	updateHeight(x)
	updateHeight(y)

	return y
}

func rightRotate(y *Node) *Node {
	x := y.Left
	y.Left = x.Right
	x.Right = y

	updateHeight(y)
	updateHeight(x)

	return x
}

func insert(node *Node, key int) *Node {
	if node == nil {
		return NewNode(key)
	}

	if key < node.Key {
		node.Left = insert(node.Left, key)
	} else if key > node.Key {
		node.Right = insert(node.Right, key)
	} else {
		return node
	}

	updateHeight(node)

	balance := getBalance(node)

	// Left Left Case
	if balance > 1 && key < node.Left.Key {
		return rightRotate(node)
	}
	// Right Right Case
	if balance < -1 && key > node.Right.Key {
		return leftRotate(node)
	}
	// Left Right Case
	if balance > 1 && key > node.Left.Key {
		node.Left = leftRotate(node.Left)
		return rightRotate(node)
	}
	// Right Left Case
	if balance < -1 && key < node.Right.Key {
		node.Right = rightRotate(node.Right)
		return leftRotate(node)
	}

	return node
}

func GenerateTree(count int) *AVLTree {
	tree := AVLTree{}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < count; i++ {
		key := rand.Intn(1000)
		tree.Insert(key)
	}

	return &tree
}
