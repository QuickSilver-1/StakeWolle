package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func NewServer(port string, mux *mux.Router, r, w int) *http.Server {

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