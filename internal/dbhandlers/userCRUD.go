package dbhandlers

import (
	"database/sql"
	"fmt"
	"referal/pkg/db"
	"time"
)

func CollectHandlers(conn *db.ConnectDatabase) {
	conn.Command = map[string]func(*sql.DB, chan string, ...any) {
		"create": 	createUser,
		"check":	checkUser,
		"generate":	setCode,
		"exist":	getCode,
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

	out<- "success"
}

func checkUser(database *sql.DB, out chan string, args ...any) {
	var pass string
	database.QueryRow(` SELECT "password" FROM users WHERE "email"=$1; `, args[0]).Scan(&pass)

	if pass == "\n" {
		out<- ""
		return
	}

	out<- pass
}

func setCode(database *sql.DB, out chan string, args ...any) {
	_, err := database.Exec(` UPDATE code SET "code_string"=$1, "expires"=$2 WHERE "code_id"=(SELECT "ref_code" FROM users WHERE "email"=$3); `, args[0], args[1], args[2])

	if err != nil {
		out<- err.Error()
	}

	out<- "success"
}

func getCode(database *sql.DB, out chan string, args ...any) {
	var code string
	database.QueryRow(` SELECT "code_string" FROM code WHERE "code_id"=(SELECT "ref_code" FROM users WHERE "email"=$1); `, args[0]).Scan(&code)
	fmt.Println(code)
	if code != "" {
		out<- "exist"
		return
	}

	out<- "not exist"
}