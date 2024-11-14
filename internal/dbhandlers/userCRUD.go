package dbhandlers

import (
	"database/sql"
	"referal/pkg/db"
	"time"
)

func CollectHandlers(conn *db.ConnectDatabase) {
	conn.Command = map[string]func(*sql.DB, chan string, ...any) {
		"create": 	createUser,
		"check":	checkUser,
	}
}

func createUser(database *sql.DB, out chan string, args ...any) {
	defer close(out)
	tx, err := database.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var codeID int
	err = tx.QueryRow(` INSERT INTO code (expires) VALUES ($1) RETURNING code_id `, time.Now().AddDate(0, 1, 0)).Scan(&codeID)

	if err != nil {
		out<- err.Error()
		return
	}
	_, err = tx.Exec(` INSERT INTO users (email, password, ref_code) VALUES ($1, $2, $3) `, args[0], args[1], codeID)

	if err != nil {
		out<- err.Error()
		return
	}

	out<- ""
}

func checkUser(database *sql.DB, out chan string, args ...any) {
	var pass string
	err := database.QueryRow(` SELECT password FROM user WHERE email=$1 `, args[0]).Scan(&pass)
	
	if err != nil {
		out<- err.Error()
		return
	}

	if len(pass) == 0 {
		out<- ""
		return
	}

	out<- pass
}

