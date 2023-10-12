package main

import (
	"encoding/json"
	"net/http"

	"github.com/dolanor/angela/merkle"
)

func (s *server) handleCreateFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "data should be sent in POST", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Files      []merkle.Content `json:"files"`
		BucketName string           `json:"bucket_name"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "couldn't decode file list: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = s.fileServer.StoreFiles(request.BucketName, request.Files)
	if err != nil {
		http.Error(w, "couldn't decode file list: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *server) handleGetFile(w http.ResponseWriter, r *http.Request) {
}
