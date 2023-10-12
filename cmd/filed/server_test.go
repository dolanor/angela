package main

import (
	"testing"

	"github.com/dolanor/angela/merkle"
)

func TestFileServer_RequestFile(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		content := []merkle.Content{
			[]byte("hello"),
			[]byte("world"),
			[]byte("data"),
		}
		tree := merkle.FromContentSlice(content)
		rootHash := tree.Root.Hash

		fs := FileServer{
			buckets: map[string]Bucket{},
		}
		bucketName := "test bucket"
		err := fs.StoreFiles(bucketName, content)
		if err != nil {
			t.Fatal(err)
		}

		gotf, gotp, err := fs.RequestFile(bucketName, 1)
		if err != nil {
			t.Fatal(err)
		}

		if string(gotf) != string(content[1]) {
			t.Fatalf("got %q, expected %q", gotf, content[1])
		}

		err = merkle.Verify(rootHash, gotp, content[1])
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("modify file requested on server", func(t *testing.T) {
		content := []merkle.Content{
			[]byte("hello"),
			[]byte("world"),
			[]byte("data"),
		}

		tree := merkle.FromContentSlice(content)
		rootHash := tree.Root.Hash

		fs := FileServer{
			buckets: map[string]Bucket{},
		}
		bucketName := "test bucket"
		err := fs.StoreFiles(bucketName, content)
		if err != nil {
			t.Fatal(err)
		}

		fs.buckets[bucketName].files[1] = []byte("here")

		gotf, gotp, err := fs.RequestFile(bucketName, 1)
		if err != nil {
			t.Fatal(err)
		}

		if string(gotf) != string(content[1]) {
			t.Fatalf("got %q, expected %q", gotf, content[1])
		}

		err = merkle.Verify(rootHash, gotp, content[1])
		if err == nil {
			t.Fatal("the verification should fail as the file requested as been tampered on the server")
		}
	})
}
