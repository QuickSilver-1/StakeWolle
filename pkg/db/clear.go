package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func ClearDB(every time.Duration, db *sql.DB, table string, close chan interface{}) {
	ticker := time.NewTicker(every)

	for {
		select {
		case <-ticker.C:
			deleteExpired(db, table)
		case <-close:
			ticker.Stop()
			return
		}
	}
}

func deleteExpired(db *sql.DB, table string) {
	query := fmt.Sprintf("DELETE FROM %s WHERE code_expires < NOW();", table)
	_, err := db.Exec(query)

	if err != nil {
		log.Printf("Ошибка удаления истекших кодов: %v", err)

	} else {
		log.Println("Истекшие коды успешно удалены")
	}
}