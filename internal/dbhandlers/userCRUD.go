package dbhandlers

import (
	"database/sql"
	"fmt"
	"time"

	"referal/pkg/log"
)

type DBData struct {
    UserEmail   string
    UserPass    string
    UserId      string
    DayExpires  time.Time
    RefString   string
}

// createUser создает нового пользователя в базе данных
func createUser(database *sql.DB, out chan string, dataI interface{}) {
    defer close(out)
    data := dataI.(DBData)

    check := make(chan string)
    go checkCode(database, check, data)

    if data.RefString == "" {
        _, err := database.Exec(` INSERT INTO users (email, password) VALUES ($1, $2); `, data.UserEmail, data.UserPass)
        
        if err != nil {
            if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
                out<- "Такой пользователь уже существует"
                log.Logger.Warn(fmt.Sprintf("Пользователь с email: %s уже существует", data.UserEmail))
                return
            }

            out<- err.Error()
            log.Logger.Error(fmt.Sprintf("Ошибка при создании пользователя с email: %s, ошибка: %v", data.UserEmail, err))
            return
        }

        log.Logger.Info(fmt.Sprintf("Пользователь с email: %s успешно создан без реферального кода", data.UserEmail))
    } else {
        ref := <-check
        if ref == "" {
            out<- "Несуществующий реферальный код"
            log.Logger.Warn(fmt.Sprintf("Несуществующий реферальный код для пользователя с email: %s", data.UserEmail))
            return
        }
        
        _, err := database.Exec(` INSERT INTO users (email, password, ref_id) VALUES ($1, $2, $3); `, data.UserEmail, data.UserPass, ref)
    
        if err != nil {
            if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
                out<- "Такой пользователь уже существует"
                log.Logger.Warn(fmt.Sprintf("Пользователь с email: %s уже существует", data.UserEmail))
                return
            }

            out<- err.Error()
            log.Logger.Error(fmt.Sprintf("Ошибка при создании пользователя с email: %s с реферальным кодом, ошибка: %v", data.UserEmail, err))
            return
        }

        log.Logger.Info(fmt.Sprintf("Пользователь с email: %s успешно создан с реферальным кодом", data.UserEmail))
    }

    out<- "success"
}

// checkUser проверяет существование пользователя по email
func checkUser(database *sql.DB, out chan string, dataI interface{}) {
    defer close(out)
    data := dataI.(DBData)

    var pass string
    err := database.QueryRow(` SELECT "password" FROM users WHERE "email"=$1; `, data.UserEmail).Scan(&pass)
    if err != nil {
        out<- ""
        log.Logger.Error(fmt.Sprintf("Ошибка при проверке пользователя с email: %s, ошибка: %v", data.UserEmail, err))
        return
    }
    
    if pass == "\n" {
        out<- ""
        log.Logger.Info(fmt.Sprintf("Пользователь с email: %s не найден", data.UserEmail))
        return
    }

    out<- pass
    log.Logger.Info(fmt.Sprintf("Пользователь с email: %s успешно проверен", data.UserEmail))
}