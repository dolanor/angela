package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dolanor/angela/merkle"
	"github.com/gorilla/mux"
)

func (s *server) handleCreateFiles(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Files      []merkle.Content `json:"files"`
		BucketName string           `json:"bucket_name"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "could not decode file list: "+err.Error(), http.StatusBadRequest)
		return
	}
	s.logger.Info("create files", "bucket_name", request.BucketName, "file_amount", len(request.Files))

	err = s.fileServer.StoreFiles(request.BucketName, request.Files)
	if err != nil {
		http.Error(w, "could not store file list: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *server) handleGetFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucket_name"]
	fileNumberStr := vars["file_number"]

	fileNumber, err := strconv.ParseUint(fileNumberStr, 10, 64)
	if err != nil {
		http.Error(w, "could not parse file number, it should be positive: "+err.Error(), http.StatusBadRequest)
		return
	}
	s.logger.Info("get file", "bucket_name", bucketName, "file_number", fileNumber)

	content, proof, err := s.fileServer.RequestFile(bucketName, uint(fileNumber))
	if err != nil {
		http.Error(w, "could not get the file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Content []byte
		Proof   []merkle.ProofStep
	}{
		Content: content,
		Proof:   proof,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "could not encode the file data: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
