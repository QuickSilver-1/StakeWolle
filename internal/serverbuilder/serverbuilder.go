package serverbuilder

import (
	"net/http"
	"referal/internal/handlers"
	"referal/pkg/server"

	"github.com/gorilla/mux"
)

func MakeServer(port string, readWait, writeWait int) *http.Server {
	mux := mux.NewRouter()

	mux.Use(server.Middleware)
	mux.HandleFunc("/signup", handlers.SignUp).Methods("POST")
	mux.HandleFunc("/signin", handlers.SignIn).Methods("POST")

	afterAuth := mux.PathPrefix("/").Subrouter()
	afterAuth.Use(server.CheckJWT)

	afterAuth.HandleFunc("/generate", handlers.GenRef).Methods("GET")
	afterAuth.HandleFunc("/delete", handlers.DelRef).Methods("GET")

	afterAuth.HandleFunc("/code", handlers.GetCode).Methods("GET")
	afterAuth.HandleFunc("/ref", handlers.GetRefs).Methods("GET")

	return server.NewServer(port, mux, readWait, writeWait)
}