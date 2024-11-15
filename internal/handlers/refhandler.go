package handlers

import (
	"fmt"
	"net/http"
	"referal/pkg/db"
	"referal/pkg/server"
	"strconv"
	"sync"
	"time"
)

type answer struct {
	statusCode	int
	value		interface{}
}

func GenRef(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Cookie("JWT")
	email, err := server.DecodeJWT(session.Value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	code := genHash(email + time.Now().String())
	out := make(chan string)

	go db.DB.Query("exist", out, email)

	dayS := r.FormValue("day")
	day, err := strconv.Atoi(dayS)

	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	exist := <-out

	if exist == "exist" {
		http.Error(w, "Перед генерацией нового кода необходимо удалить старый", http.StatusBadRequest)
		return
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer close(out)
		defer wg.Done()
		res := <-out

		if res != "success" {
			fmt.Print(res)
			http.Error(w, "Ошибка записи в базу данных", http.StatusInternalServerError)
			return
		}

	}()
	go db.DB.Exec("generate", out, code, time.Now().AddDate(0, 0, day), email)

	wg.Wait()
	w.Write([]byte(code))
}