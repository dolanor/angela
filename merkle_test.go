package merkle_test

import (
	"testing"

	"github.com/dolanor/merkle"
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
	// sha3.ShakeSum256(b, []byte("yoyo"))
	// t.Fatal(tree.Belongs(b))
	_ = b
}

func TestSendFiles(t *testing.T) {
	rootHash, err := merkle.SendFiles([]merkle.Content{
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
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("zozo")

	ok, err := merkle.Verify(rootHash, merkle.Content(data))
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("content not valid: not ok")
	}
}
