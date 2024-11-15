package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"referal/pkg/db"
	"referal/pkg/server"
	"strings"
	"time"
	"unicode"
)

func MainPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../web/html/main.html")
}

func SignInPage(w http.ResponseWriter, r *http.Request) {	
	http.ServeFile(w, r, "../../web/html/signin.html")
}

func SignUpPage(w http.ResponseWriter, r *http.Request) {	
	http.ServeFile(w, r, "../../web/html/signup.html")
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pass := r.FormValue("password")
	
	out := make(chan string)
	go db.DB.Query("check", out, email)

	hashPass, err := validPass(pass)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	correctPass := <-out	

	if correctPass == "" {
		http.Error(w, pass, http.StatusBadRequest)
		return
	}

	if correctPass == hashPass {
		err = server.MakeToken(email, w)
		
		if err != nil {
			http.Error(w, "Ошибка создания токена", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else {
		http.Error(w, "Неправильный email или пароль", http.StatusUnauthorized)
	}
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pass := r.FormValue("password")

	hashPass, err := validPass(pass)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out := make(chan string)
	e := make(chan bool)

	go func() {
		defer close(e)

		err := <-out
		if err != "success" {
			http.Error(w, err, http.StatusBadRequest)
			e<- true
		} else {
			e<- false
		}
	}()

	go db.DB.Exec("create", out, email, hashPass)

	if <-e {
		w.Write([]byte("123123"))
		return
	}

	err = server.MakeToken(email, w)
		
	if err != nil {
		http.Error(w, "Ошибка создания токена", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("JWT")

	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	http.Redirect(w, r, "/signin", http.StatusFound)
}

func validPass(pass string) (string, error) {
	if len(pass) < 3 || len(pass) > 30 {
		return "", fmt.Errorf("длина пароля должна быть от 3 до 30 символов включительно")
	}

	hasUpper, hasDigit := false, false
	
	for _, char := range pass {
		switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsDigit(char):
				hasDigit = true
			case !unicode.IsLetter(char) && !unicode.IsDigit(char) && !strings.Contains("_!@#&*-", string(char)):
				return "", fmt.Errorf("пароль должен состоять из букв латнского алфавита, цифр и сиволов _!@#&*-")
			}
		}

	if !hasDigit || !hasUpper {
		return "", fmt.Errorf("пароль должен содержать хотя бы 1 заглвную букву и 1 цифру")
	}

	hash := genHash(pass)
	return hash, nil
}

func genHash(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash)
}