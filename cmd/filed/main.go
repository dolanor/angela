package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	s := server{
		fileServer: FileServer{
			buckets: map[string]Bucket{},
		},
	}

	m := mux.NewRouter()
	m.HandleFunc("/files", s.handleCreateFiles).Methods(http.MethodPost)
	m.HandleFunc("/files/{bucket_name}/{file_number}", s.handleGetFile).Methods(http.MethodGet)
	err := http.ListenAndServe(":7777", m)
	if err != nil {
		panic(err)
	}
}
