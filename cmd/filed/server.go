package main

import (
	"errors"
	"log"
	"sync"

	"github.com/dolanor/angela/merkle"
	"golang.org/x/crypto/sha3"
)

type server struct {
	fileServer FileServer
}

type FileServer struct {
	mu      sync.Mutex
	buckets map[string]Bucket
}

type Bucket struct {
	name       string
	files      []merkle.Content
	merkleTree merkle.Tree
}

type File struct {
	Data []byte
}

func (fs *FileServer) StoreFiles(bucketName string, content []merkle.Content) error {
	t := merkle.FromContentSlice(content)
	log.Printf("added files:\n%s", t)

	bucket := Bucket{
		name:       bucketName,
		files:      content,
		merkleTree: t,
	}

	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.buckets[bucket.name] = bucket

	return nil
}

func (fs *FileServer) RequestFile(bucketName string, fi uint) (merkle.Content, []merkle.ProofStep, error) {
	bucket, ok := fs.buckets[bucketName]
	if !ok {
		return merkle.Content{}, nil, errors.New("unknown bucket")
	}

	if int(fi) >= len(bucket.files) {
		return merkle.Content{}, nil, errors.New("fi is too big")
	}

	content := bucket.files[fi]
	hash := make([]byte, 64)
	sha3.ShakeSum256(hash, content)
	merkleProof := bucket.merkleTree.GenerateProof(hash)

	return content, merkleProof, nil
}
