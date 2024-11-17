package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	dh "referal/internal/dbhandlers"
	"referal/pkg/db"
	"referal/pkg/log"
	"referal/pkg/server"
)

// GenRef генерирует реферальный код для пользователя
func GenRef(w http.ResponseWriter, r *http.Request) {
    tokenS := r.Header.Get("Authorization")
    tokenS = strings.Split(tokenS, " ")[1]

    email, err := server.DecodeJWT(tokenS)  // Декодируем JWT для получения email
    if err != nil {
        server.AnswerHandler(w, 500, err.Error())
        log.Logger.Error(fmt.Sprintf("Ошибка декодирования JWT: %v", err))
        return
    }

    code := genHash(email + time.Now().String())  // Генерируем реферальный код
    out := make(chan string)
    go db.DB.Query("get", out, dh.DBData{
        UserEmail: email,
    })

    dayS := r.FormValue("day")
    day, err := strconv.Atoi(dayS)  // Конвертируем строку в число

    if err != nil {
        server.AnswerHandler(w, 400, "Неверный формат данных")
        log.Logger.Error(fmt.Sprintf("Неверный формат данных: %v", err))
        return
    }

    exist := <-out
    if exist != "" {
        server.AnswerHandler(w, 400, "Перед генерацией нового кода необходимо удалить старый")
        log.Logger.Info("Перед генерацией нового кода необходимо удалить старый")
        return
    }

    wg := &sync.WaitGroup{}
    var rout chan string

    wg.Add(1)
    go func() {
        defer wg.Done()
        res := <-out
        if res != "success" {
            server.AnswerHandler(w, 500, "Ошибка записи в базу данных")
            log.Logger.Error(fmt.Sprintf("Ошибка записи в базу данных: %v", err))
            return
        }

        ctx := context.Background()
        rout = make(chan string)
        go db.NewKey(ctx, email, code, rout)  // Сохраняем код в Redis
    }()

    out = make(chan string)
    go db.DB.Query("generate", out, dh.DBData{
        RefString: code,
        DayExpires: time.Now().AddDate(0, 0, day),
        UserEmail: email})

    wg.Wait()
    res := <-rout

    if res != "success" {
        server.AnswerHandler(w, 500, res)
        log.Logger.Error(fmt.Sprintf("Ошибка генерации ключа: %v", res))
    }

    server.AnswerHandler(w, 200, code)
    log.Logger.Info(fmt.Sprintf("Реферальный код успешно сгенерирован: %s", code))
}

// DelRef удаляет реферальный код для пользователя
func DelRef(w http.ResponseWriter, r *http.Request) {
    tokenS := r.Header.Get("Authorization")
    tokenS = strings.Split(tokenS, " ")[1]

    email, err := server.DecodeJWT(tokenS)  // Декодируем JWT для получения email
    if err != nil {
        server.AnswerHandler(w, 500, err.Error())
        log.Logger.Error(fmt.Sprintf("Ошибка декодирования JWT: %v", err))
        return
    }

    out := make(chan string)
    go db.DB.Query("delete", out, dh.DBData{
        UserEmail: email,
    })

    ctx := context.Background()
    rout := make(chan string)
    go db.DelKey(ctx, email, rout)  // Удаляем код из Redis

    ans := <-out
    sliseAns := strings.Split(ans, " ")
    code, _ := strconv.Atoi(sliseAns[0])
    value := strings.Join(sliseAns[1:], " ")

    res := <-rout
    if res != "success" {
        server.AnswerHandler(w, 500, res)
        log.Logger.Error(fmt.Sprintf("Ошибка удаления ключа: %v", res))
        return
    }

    server.AnswerHandler(w, code, value)
    log.Logger.Info(fmt.Sprintf("Реферальный код успешно удален для пользователя: %s", email))
}

// GetCode получает реферальный код пользователя
func GetCode(w http.ResponseWriter, r *http.Request) {
    email := r.URL.Query().Get("email")

    ctx := context.Background()
    rout := make(chan string)
    go db.GetKey(ctx, email, rout)  // Получаем код из Redis

    out := make(chan string)
    check := make(chan string)

    code := <-rout
    if code != "" {
        server.AnswerHandler(w, 200, code)
        log.Logger.Info(fmt.Sprintf("Реферальный код успешно получен: %s", code))
        return
    }

    go db.DB.Query("get", out, dh.DBData{
        UserEmail: email,
    })
    go db.DB.Query("check", check, dh.DBData{
        UserEmail: email,
    })

    var ref bool = true
    select {
    case code = <-out:
        if code == "" {
            ref = false
            exist := <-check
            if exist == "" {
                server.AnswerHandler(w, 400, "Пользователя с такой почтой не существует")
                log.Logger.Info(fmt.Sprintf("Пользователь с почтой %s не существует", email))
                return
            }
        }

    case exist := <-check:
        if exist == "" {
            server.AnswerHandler(w, 400, "Пользователя с такой почтой не существует")
            log.Logger.Info(fmt.Sprintf("Пользователь с почтой %s не существует", email))
            return
        }

        code = <-out
        if code == "" {
            ref = false
        }
    }

    if !ref {
        server.AnswerHandler(w, 400, "У этого пользователя нет реферального кода")
        log.Logger.Info(fmt.Sprintf("У пользователя %s нет реферального кода", email))
        return
    }

    server.AnswerHandler(w, 200, code)
    log.Logger.Info(fmt.Sprintf("Реферальный код успешно получен: %s", code))
}

// GetRefs получает рефералов пользователя по ID
func GetRefs(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")

    ctx := context.Background()
    rout := make(chan string)
    go db.GetKey(ctx, id, rout)  // Получаем реферальные данные из Redis

    check := make(chan string)
    code := <-rout

    if code != "" {
        server.AnswerHandler(w, 200, code)
        log.Logger.Info(fmt.Sprintf("Реферальные данные успешно получены: %s", code))
        return
    }

    go db.DB.Query("checkid", check, dh.DBData{
        UserId: id,
    })

    out := make(chan string)
    wg := &sync.WaitGroup{}
    wg.Add(1)

    pass := <-check
    if pass == "" {
        server.AnswerHandler(w, 400, "Пользователя с таким id не существует")
        log.Logger.Info(fmt.Sprintf("Пользователь с ID %s не существует", id))
        return
    }

    go func() {
        defer wg.Done()
        users := []string{}

        for ref := range out {
            users = append(users, ref)  // Собираем всех рефералов пользователя
        }

        if len(users) == 0 {
            server.AnswerHandler(w, 200, "У данного пользователя нет рефералов")
            log.Logger.Info(fmt.Sprintf("У пользователя с ID %s нет рефералов", id))
            return
        }

        server.AnswerHandler(w, 200, users)
        log.Logger.Info(fmt.Sprintf("Рефералы пользователя с ID %s успешно получены", id))

        rout = make(chan string)
        go db.NewKey(ctx, id, strings.Join(users, " "), rout)  // Сохраняем реферальные данные в Redis
    }()

    go db.DB.Query("referals", out, dh.DBData{
        UserId: id,
    })

    wg.Wait()

    res := <-rout
    if res != "success" {
        log.Logger.Error(fmt.Sprintf("Ошибка при сохранении реферальных данных: %v", res))
        return
    }
}
