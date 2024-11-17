package db

import (
	"database/sql"
	"fmt"
	"time"

	"referal/pkg/log"
)

// ClearDB запускает процесс очистки базы данных через заданные интервалы времени
func ClearDB(every time.Duration, db *sql.DB, table string, close chan interface{}) {
    ticker := time.NewTicker(every)

    for {
        select {
        case <-ticker.C:
            log.Logger.Info("Начало очистки базы данных")
            DeleteExpired(db, table)
        case <-close:
            ticker.Stop()
            return
        }
    }
}

// deleteExpired удаляет записи, срок действия которых истек, из указанной таблицы
func DeleteExpired(db *sql.DB, table string) {
    _, err := db.Exec(fmt.Sprintf(` DELETE FROM %s WHERE "expires" < NOW(); `, table))

    if err != nil {
        log.Logger.Info(fmt.Sprintf("Ошибка удаления истекших кодов: %v", err))
    } else {
        log.Logger.Info("Истекшие коды успешно удалены")
    }
}
