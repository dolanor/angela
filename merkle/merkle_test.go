package merkle_test

import (
	"testing"

	"github.com/dolanor/angela/merkle"
)

func TestFromContentSlice(t *testing.T) {
	tree := merkle.FromContentSlice([]merkle.Content{
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
	_ = b
}
