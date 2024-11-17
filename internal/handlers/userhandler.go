package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	dh "referal/internal/dbhandlers"
	"referal/pkg/db"
	"referal/pkg/server"
	"strings"
	"unicode"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		fmt.Println(err)
		server.AnswerHandler(w, 400, "Неверный формат данных")
		return
	}

	out := make(chan string)
	go db.DB.Query("check", out, dh.DBData{
		UserEmail: user.Email})

	hashPass, err := validPass(user.Pass)

	if err != nil {
		server.AnswerHandler(w, 400, err.Error())
		return
	}

	correctPass := <-out	

	if correctPass == "" {
		server.AnswerHandler(w, 400, "Пользователя с такой почтой не существует")
		return
	}

	if correctPass == hashPass {
		token, err := server.MakeToken(user.Email)
		
		if err != nil {
			server.AnswerHandler(w, 500, "Ошибка создания токена")
			return
		}

		server.AnswerHandler(w, 200, map[string]string{
			"Info": "Вход",
			"Code": "Bearer " + token,
		})	

	} else {
		server.AnswerHandler(w, 401,"Неправильный пароль")
	}
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		server.AnswerHandler(w, 400, "Неверный формат данных")
		return
	}

	hashPass, err := validPass(user.Pass)

	if err != nil {
		server.AnswerHandler(w, 400, err.Error())
		return
	}

	out := make(chan string)
	e := make(chan string)

	go func() {
		err := <-out
		if err != "success" {
			e<- err
		}

		e<- ""
	}()

	go db.DB.Query("create", out, dh.DBData{
		UserEmail: user.Email,
		UserPass: hashPass,
		RefString: user.Ref,
	})

	errS := <-e
	if errS != "" {
		server.AnswerHandler(w, 400, errS)
		return
	}
	server.AnswerHandler(w, 200, "Аккаунт зарегестрирован")
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