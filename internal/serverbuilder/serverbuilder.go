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
	mux.HandleFunc("/signup", handlers.SignUpPage).Methods("GET")
	mux.HandleFunc("/signup", handlers.SignUp).Methods("POST")

	mux.HandleFunc("/signin", handlers.SignInPage).Methods("GET")
	mux.HandleFunc("/signin", handlers.SignIn).Methods("POST")

	afterAuth := mux.PathPrefix("/").Subrouter()
	afterAuth.Use(server.CheckJWT)
	afterAuth.HandleFunc("/", handlers.MainPage)
	afterAuth.HandleFunc("/signout", handlers.SignOut).Methods("GET")
	afterAuth.HandleFunc("/generate", handlers.GenRef)

	return server.NewServer(port, mux, readWait, writeWait)
}