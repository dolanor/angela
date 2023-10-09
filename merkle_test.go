package merkle_test

import (
	"testing"

	"github.com/dolanor/merkle"
)

func TestFromContentSlice(t *testing.T) {
	tree := merkle.FromContentSlice([][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte("yoyo"),
		[]byte("zozo"),
		[]byte("magot"),
		[]byte("hello"),
		[]byte("world"),
		[]byte("yoyo"),
		[]byte("magot"),
	})

	_ = tree

	t.Logf("\n%s", tree)
	b := make([]byte, 64)
	// sha3.ShakeSum256(b, []byte("yoyo"))
	// t.Fatal(tree.Belongs(b))
	_ = b
}
