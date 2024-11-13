package pkg

import (
	"net/http"
	"time"
)

func NewServer(port string, r, w int) *http.Server {

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  time.Duration(r) * time.Second,
		WriteTimeout: time.Duration(w) * time.Second,
	}

	return server
}

func StartServer(server *http.Server) error {
	err := server.ListenAndServe()

	return err
}