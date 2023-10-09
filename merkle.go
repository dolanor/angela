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

		//		fmt.Printf("hash: %c\n", []rune(hashemo.FromBytes(b))[0])
		n := Node{
			Hash: b,
		}

		nodes = append(nodes, &n)
	}

	n := buildCousins(nodes)

	return Tree{
		Root:  n,
		Nodes: nodes,
	}
}

func buildCousins(nodes []*Node) *Node {
	if len(nodes)%2 != 0 {
		nodes = append(nodes, nodes[len(nodes)-1])
	}

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

func (t Tree) String() string {
	return buildCousinLine(0, t.Nodes)
}

func buildCousinLine(depth int, nodes []*Node) string {
	prefix := ""
	middle := " "
	for i := 0; i < depth; i++ {
		prefix += "  "
		middle += "  "
	}

	var s string
	var parents []*Node

	for i, n := range nodes {
		if i == 0 {
			s += fmt.Sprintf("%s", prefix)
		}
		s += fmt.Sprintf("%s%s", hashemo.FromBytes(n.Hash[0:1]), middle)
		if i%2 == 0 {
			if n.Parent == nil {
				continue
			}
			parents = append(parents, n.Parent)
		}
	}
	s += fmt.Sprintln()
	if len(parents) == 0 {
		return s
	}

	depth++
	s += buildCousinLine(depth, parents)

	return s
}
