package merkle

import (
	"fmt"

	"github.com/dolanor/hashemo"
	"golang.org/x/crypto/sha3"
)

type Tree struct {
	Root  *Node
	Nodes []*Node
}

type Node struct {
	Parent *Node
	Hash   []byte
}

func FromContentSlice(content [][]byte) Tree {
	var nodes []*Node
	for _, data := range content {
		b := make([]byte, 64)
		sha3.ShakeSum256(b, data)

		fmt.Println("hash:", hashemo.FromBytes(b))
		n := Node{
			Hash: b,
		}

		nodes = append(nodes, &n)
	}

	if len(nodes)%2 != 0 {
		nodes = append(nodes, nodes[len(nodes)-1])
	}

	n := buildCousins(nodes)

	return Tree{
		Root:  n,
		Nodes: nodes,
	}
}

func buildCousins(nodes []*Node) *Node {
	var parents []*Node
	for i := 0; i < len(nodes); i += 2 {
		left, right := nodes[i], nodes[i+1]
		b := make([]byte, 64)
		bb := left.Hash
		bb = append(bb, right.Hash...)
		sha3.ShakeSum256(b, bb)
		parent := Node{
			Hash: b,
		}

		left.Parent = &parent
		right.Parent = &parent
		parents = append(parents, &parent)
	}

	if len(parents) == 1 {
		return parents[0]
	}

	n := buildCousins(parents)
	return n
}

func (t Tree) Belongs(hash []byte) bool {
	fmt.Println("======")
	fmt.Println(hashemo.FromBytes(hash))
	fmt.Println()
	for _, v := range t.Nodes {
		fmt.Println(hashemo.FromBytes(v.Hash))
		fmt.Println()
		if string(v.Hash) == string(hash) {
			return true
		}
	}
	return false
}

func (t Tree) GenerateProof(hash []byte) [][]byte {
	// FIXME
	return nil
}
