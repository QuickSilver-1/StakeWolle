package serverbuilder

import (
	"net/http"

	"referal/internal/handlers"
	"referal/pkg/server"

	"github.com/gorilla/mux"
)

// MakeServer создает и настраивает HTTP сервер с маршрутизацией
func MakeServer(port string, readWait, writeWait int) *http.Server {
    mux := mux.NewRouter()

    // Устанавливаем middleware для логирования запросов и защиты от brute force
    mux.Use(server.Middleware)
    mux.Use(server.LimitMiddleware)
    
    // Маршруты для регистрации и авторизации
    mux.HandleFunc("/signup", handlers.SignUp).Methods("POST")
    mux.HandleFunc("/signin", handlers.SignIn).Methods("POST")

    // Подмаршруты, требующие авторизации
    afterAuth := mux.PathPrefix("/").Subrouter()
    afterAuth.Use(server.CheckJWT)

    // Маршрут для генерации реферального кода
    afterAuth.HandleFunc("/generate", handlers.GenRef).Methods("GET")
    
    // Маршрут для удаления реферального кода
    afterAuth.HandleFunc("/delete", handlers.DelRef).Methods("GET")

    // Маршрут для получения реферального кода по email
    afterAuth.HandleFunc("/code", handlers.GetCode).Methods("GET")
    
    // Маршрут для получения рефералов по ID пользователя
    afterAuth.HandleFunc("/ref", handlers.GetRefs).Methods("GET")

    return server.NewServer(port, mux, readWait, writeWait)
}
