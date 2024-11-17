package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Answer struct {
	StatusCode int
	Value      interface{}
}

func AnswerHandler(w http.ResponseWriter, code int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")

	answer := Answer{
		StatusCode: code,
		Value: value,
	}

	err := json.NewEncoder(w).Encode(answer)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

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