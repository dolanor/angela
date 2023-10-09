package merkle_test

import (
	"testing"

	"github.com/dolanor/merkle"
	"golang.org/x/crypto/sha3"
)

func TestFromContentSlice(t *testing.T) {
	tree := merkle.FromContentSlice([][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte("yoyo"),
	})

	_ = tree

	b := make([]byte, 64)
	sha3.ShakeSum256(b, []byte("yoyo"))
	t.Fatal(tree.Belongs(b))
}
