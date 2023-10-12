package main

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/dolanor/angela/merkle"
)

var (
	ErrUnknownBucket       = errors.New("unknown bucket")
	ErrFileIndexOutOfRange = errors.New("file index out of range")
)

type server struct {
	fileServer FileServer

	logger *slog.Logger
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

// StoreFiles stores files in the storage system.
func (fs *FileServer) StoreFiles(bucketName string, content []merkle.Content) error {
	t := merkle.FromContentSlice(content)

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

// RequestFile extracts a files from the system and generates a merkle proof for this file.
func (fs *FileServer) RequestFile(bucketName string, fi uint) (merkle.Content, []merkle.ProofStep, error) {
	bucket, ok := fs.buckets[bucketName]
	if !ok {
		return merkle.Content{}, nil, fmt.Errorf("%q: %w", bucketName, ErrUnknownBucket)
	}

	if int(fi) >= len(bucket.files) {
		return merkle.Content{}, nil, fmt.Errorf("%q: %w", fi, ErrFileIndexOutOfRange)
	}

	content := bucket.files[fi]
	hash := content.Hash()
	merkleProof := bucket.merkleTree.GenerateProof(hash)

	return content, merkleProof, nil
}
