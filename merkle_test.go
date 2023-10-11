package merkle_test

import (
	"testing"

	"github.com/dolanor/hashemo"
	"github.com/dolanor/merkle"
	"golang.org/x/crypto/sha3"
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
	content := []merkle.Content{
		[]byte("hello"),
		[]byte("world"),
		[]byte("yoyo"),
		[]byte("zozo"),
		[]byte("magot"),
	}
	tree := merkle.FromContentSlice(content)
	rootHash, err := merkle.SendFiles(content)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("zozo")
	b := make([]byte, 64)
	sha3.ShakeSum256(b, data)

	t.Logf("\n%s", tree)
	proof := tree.GenerateProof(b)

	for _, p := range proof {
		t.Logf("proof: %s\n", hashemo.FromBytes(p.Hash)[:4])
	}

	err = merkle.Verify(rootHash, proof, merkle.Content(data))
	if err != nil {
		t.Fatal(err)
	}
}
