package dbhandlers

import (
	"database/sql"
	"fmt"

	"referal/pkg/log"
)

// setCode устанавливает реферальный код для пользователя
func setCode(database *sql.DB, out chan string, dataI interface{}) {
    defer close(out)
    data := dataI.(DBData)

    var codeID int
    err := database.QueryRow(` INSERT INTO code (code_string, expires) VALUES ($1, $2) RETURNING code_id `, data.RefString, data.DayExpires).Scan(&codeID)
    
    if err != nil {
        out<- err.Error()
        log.Logger.Error(fmt.Sprintf("Ошибка при вставке кода: %v", err))
        return
    }

    _, err = database.Exec(` UPDATE users SET "ref_code"=$1 WHERE "email"=$2; `, codeID, data.UserEmail)
    if err != nil {
        out<- err.Error()
        log.Logger.Error(fmt.Sprintf("Ошибка при обновлении пользователя: %v", err))
        return
    }

    out<- "success"
    log.Logger.Info("Код успешно установлен")
}

// getCode получает реферальный код для пользователя
func getCode(database *sql.DB, out chan string, dataI interface{}) {
    defer close(out)
    data := dataI.(DBData)

    var code string
    err := database.QueryRow(` SELECT "code_string" FROM code WHERE "code_id"=(SELECT "ref_code" FROM users WHERE "email"=$1); `, data.UserEmail).Scan(&code)
    if err != nil {
        out<- ""
        log.Logger.Error(fmt.Sprintf("Ошибка при получении кода: %v", err))
        return
    }

    out<- code
    log.Logger.Info("Код успешно получен")
}

// delCode удаляет реферальный код для пользователя
func delCode(database *sql.DB, out chan string, dataI interface{}) {
    defer close(out)
    data := dataI.(DBData)

    query := make(chan string)

    go getCode(database, query, data)

    exist := <-query 
    if exist == "" {
        out<- "400 Кода не существует"
        log.Logger.Warn("Код не существует")
        return
    }

    tx, err := database.Begin()
    if err != nil {
        out<- "500 Не получилось сформировать транзакцию"
        log.Logger.Error(fmt.Sprintf("Ошибка при создании транзакции: %v", err))
        return
    }

    defer func() {
        if err != nil {
            tx.Rollback()
            out<- "500 Ошибка удаления"
            log.Logger.Error(fmt.Sprintf("Ошибка при удалении кода: %v", err))
        } else {
            tx.Commit()
            out<- "200 Код удален"
            log.Logger.Info("Код успешно удален")
        }
    }()

    _, err = tx.Exec(` DELETE FROM code WHERE "code_id"=(SELECT "ref_code" FROM users WHERE "email"=$1) `, data.UserEmail)
    if err != nil {
        out<- err.Error()
        log.Logger.Error(fmt.Sprintf("Ошибка при удалении кода: %v", err))
        return
    }
        
    _, err = tx.Exec(` UPDATE users SET "ref_code"=NULL WHERE "email"=$1 `, data.UserEmail)
    if err != nil {
        out<- err.Error()
        log.Logger.Error(fmt.Sprintf("Ошибка при обновлении пользователя: %v", err))
        return
    }

    _, err = tx.Exec(` UPDATE users SET "ref_id"=NULL WHERE "ref_id"=(SELECT "user_id" FROM users WHERE "email"=$1) `, data.UserEmail)
    if err != nil {
        out<- err.Error()
        log.Logger.Error(fmt.Sprintf("Ошибка при обновлении ref_id у пользователя: %v", err))
        return
    }
}

// checkCode проверяет реферальный код пользователя
func checkCode(database *sql.DB, out chan string, dataI interface{}) {
    defer close(out)
    data := dataI.(DBData)

    var user_id string
    err := database.QueryRow(` SELECT "user_id" FROM users WHERE "ref_code"=( SELECT "code_id" FROM code WHERE "code_string"=$1); `, data.RefString).Scan(&user_id)
    if err != nil {
        out<- ""
        log.Logger.Error(fmt.Sprintf("Ошибка при проверке кода: %v", err))
        return
    }

    out<- user_id
    log.Logger.Info("Проверка кода завершена успешно")
}

// getRefBD получает реферальные данные для пользователя
func getRefBD(database *sql.DB, out chan string, dataI interface{}) {
    defer close(out)
    data := dataI.(DBData)

    rows, err := database.Query(` SELECT "user_id" FROM users WHERE "ref_id"=$1 `, data.UserId)
    if err != nil {
        out<- err.Error()
        log.Logger.Error(fmt.Sprintf("Ошибка при получении реферальных данных: %v", err))
        return
    }

    for rows.Next() {
        var id string
        err = rows.Scan(&id)
        if err != nil {
            out<- err.Error()
            log.Logger.Error(fmt.Sprintf("Ошибка при чтении строки: %v", err))
            return
        }

        out<- id
    }
    log.Logger.Info("Получение реферальных данных завершено успешно")
}