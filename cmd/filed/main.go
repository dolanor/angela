package main

import (
	"net/http"
)

func main() {
	s := server{
		fileServer: FileServer{
			buckets: map[string]Bucket{},
		},
	}

	http.HandleFunc("/files/", s.handleCreateFiles)
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		panic(err)
	}
}
