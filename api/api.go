package api

import "github.com/dolanor/angela/merkle"

type CreateFilesRequest struct {
	Files      []merkle.Content `json:"files"`
	BucketName string           `json:"bucket_name"`
}

type GetFileResponse struct {
	Content []byte
	Proof   []merkle.ProofStep
}
