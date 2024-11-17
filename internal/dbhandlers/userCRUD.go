package dbhandlers

import (
	"database/sql"
	"time"
)

type DBData struct {
	UserEmail	string
	UserPass	string
	UserId		string
	DayExpires	time.Time
	RefString	string
}

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
				return
			}

			out<- err.Error()
			return
		}

	} else {
		ref := <-check
		if ref == "" {
			out<- "Несуществующий реферальный код"
			return
		}
		
		_, err := database.Exec(` INSERT INTO users (email, password, ref_id) VALUES ($1, $2, $3); `, data.UserEmail, data.UserPass, ref)
	
		if err != nil {
			if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
				out<- "Такой пользователь уже существует"
				return
			}

			out<- err.Error()
			return
		}
	}

	out<- "success"
}

func checkUser(database *sql.DB, out chan string, dataI interface{}) {
	defer close(out)
	data := dataI.(DBData)

	var pass string
	database.QueryRow(` SELECT "password" FROM users WHERE "email"=$1; `, data.UserEmail).Scan(&pass)
	
	if pass == "\n" {
		out<- ""
		return
	}

	out<- pass
}

func checkUserID(database *sql.DB, out chan string, dataI interface{}) {
	defer close(out)
	data := dataI.(DBData)

	var pass string
	database.QueryRow(` SELECT "password" FROM users WHERE "user_id"=$1; `, data.UserId).Scan(&pass)

	if pass == "\n" {
		out<- ""
		return
	}

	out<- pass
}