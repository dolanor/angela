package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	cfg, err := loadConfig(os.Args[1:])
	if err != nil {
		panic(err)
	}

	logger, err := getLogger(cfg.logFormat, cfg.logLevel)
	if err != nil {
		panic(err)
	}

	s := server{
		fileServer: FileServer{
			logger:  logger,
			buckets: map[string]Bucket{},
		},

		logger: logger,
	}

	r := mux.NewRouter()
	r.HandleFunc("/files", s.handleCreateFiles).Methods(http.MethodPost)
	r.HandleFunc("/files/{bucket_name}/{file_number}", s.handleGetFile).Methods(http.MethodGet)

	hostPort := fmt.Sprintf(":%d", cfg.port)
	logger.Info("listening", "host_port", hostPort)

	err = http.ListenAndServe(hostPort, r)
	if err != nil {
		panic(err)
	}
}
