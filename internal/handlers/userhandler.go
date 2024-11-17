package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	dh "referal/internal/dbhandlers"
	"referal/pkg/db"
	"referal/pkg/log"
	"referal/pkg/server"
)

// SignIn обрабатывает запрос на вход пользователя
func SignIn(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)  // Декодируем тело запроса в структуру User

    if err != nil {
        server.AnswerHandler(w, 400, "Неверный формат данных")
        log.Logger.Error(fmt.Sprintf("Ошибка декодирования данных: %v", err))
        return
    }

    
    log.Logger.Info("Проверка валидности email")
    if !isValidEmail(user.Email) {
        server.AnswerHandler(w, 400, "Неверная форма email")
        return
    }

    out := make(chan string)
    go db.DB.Query("check", out, dh.DBData{  // Проверяем наличие пользователя в базе данных
        UserEmail: user.Email})

    hashPass, err := validPass(user.Pass)  // Проверяем и хэшируем пароль

    if err != nil {
        server.AnswerHandler(w, 400, err.Error())
        log.Logger.Info(fmt.Sprintf("Неверный пароль: %v", err))
        return
    }

    correctPass := <-out  // Получаем хэш пароля из базы данных

    if correctPass == "" {
        server.AnswerHandler(w, 400, "Пользователя с такой почтой не существует")
        log.Logger.Warn(fmt.Sprintf("Пользователь с почтой %s не существует", user.Email))
        return
    }

    if correctPass == hashPass {
        token, err := server.MakeToken(user.Email)  // Генерируем токен для авторизации

        if err != nil {
            server.AnswerHandler(w, 500, "Ошибка создания токена")
            log.Logger.Error(fmt.Sprintf("Ошибка создания токена для пользователя %s: %v", user.Email, err))
            return
        }

        server.AnswerHandler(w, 200, map[string]string{
            "Info": "Вход",
            "Code": "Bearer " + token,
        })
        log.Logger.Info(fmt.Sprintf("Пользователь %s успешно авторизован", user.Email))

    } else {
        server.AnswerHandler(w, 401,"Неправильный пароль")
        log.Logger.Info(fmt.Sprintf("Неправильный пароль для пользователя %s", user.Email))
    }
}

// SignUp обрабатывает запрос на регистрацию пользователя
func SignUp(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)  // Декодируем тело запроса в структуру User

    if err != nil {
        server.AnswerHandler(w, 400, "Неверный формат данных")
        log.Logger.Error(fmt.Sprintf("Ошибка декодирования данных: %v", err))
        return
    }

    hashPass, err := validPass(user.Pass)  // Проверяем и хэшируем пароль

    if err != nil {
        server.AnswerHandler(w, 400, err.Error())
        log.Logger.Info(fmt.Sprintf("Неверный пароль: %v", err))
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

    go db.DB.Query("create", out, dh.DBData{  // Создаем пользователя в базе данных
        UserEmail: user.Email,
        UserPass: hashPass,
        RefString: user.Ref,
    })

    errS := <-e
    if errS != "" {
        server.AnswerHandler(w, 400, errS)
        log.Logger.Error(fmt.Sprintf("Ошибка регистрации пользователя %s: %v", user.Email, errS))
        return
    }
    server.AnswerHandler(w, 200, "Аккаунт зарегистрирован")
    log.Logger.Info(fmt.Sprintf("Пользователь %s успешно зарегистрирован", user.Email))
}

// validPass проверяет и хэширует пароль
func validPass(pass string) (string, error) {
    if len(pass) < 3 || len(pass) > 30 {
        return "", fmt.Errorf("длина пароля должна быть от 3 до 30 символов включительно")
    }

    hasUpper, hasDigit := false, false

    // Проверяем каждый символ пароля
    for _, char := range pass {
        switch {
            case unicode.IsUpper(char):
                hasUpper = true  // Проверяем наличие заглавной буквы
            case unicode.IsDigit(char):
                hasDigit = true  // Проверяем наличие цифры
            case !unicode.IsLetter(char) && !unicode.IsDigit(char) && !strings.Contains("_!@#&*-", string(char)):
                return "", fmt.Errorf("пароль должен состоять из букв латнского алфавита, цифр и символов _!@#&*-")
            }
        }

    if !hasDigit || !hasUpper {
        return "", fmt.Errorf("пароль должен содержать хотя бы 1 заглавную букву и 1 цифру")
    }

    hash := genHash(pass)  // Хэшируем пароль
    return hash, nil
}

// genHash создает хэш пароля
func genHash(str string) string {
    hasher := sha256.New()
    hasher.Write([]byte(str))  // Преобразуем строку в хэш
    hash := hasher.Sum(nil)

    return hex.EncodeToString(hash)  // Возвращаем хэш в виде шестнадцатеричной строки
}

// isValidEmail проверяет, является ли строка допустимым адресом электронной почты
func isValidEmail(email string) bool {
    // Регулярное выражение для проверки адреса электронной почты
    var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    
    return emailRegex.MatchString(email)
}