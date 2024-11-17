package dbhandlers

import (
	"database/sql"
)

func setCode(database *sql.DB, out chan string, dataI interface{}) {
	defer close(out)
	data := dataI.(DBData)

	var codeID int
	err := database.QueryRow(` INSERT INTO code (code_string, expires) VALUES ($1, $2) RETURNING code_id `, data.RefString, data.DayExpires).Scan(&codeID)
	
	if err != nil {
		out<- err.Error()
		return
	}

	_, err = database.Exec(` UPDATE users SET "ref_code"=$1 WHERE "email"=$2; `, codeID, data.UserEmail)

	if err != nil {
		out<- err.Error()
		return
	}

	out<- "success"
}

func getCode(database *sql.DB, out chan string, dataI interface{}) {
	defer close(out)
	data := dataI.(DBData)

	var code string
	database.QueryRow(` SELECT "code_string" FROM code WHERE "code_id"=(SELECT "ref_code" FROM users WHERE "email"=$1); `, data.UserEmail).Scan(&code)

	out<- code
}

func delCode(database *sql.DB, out chan string, dataI interface{}) {
	defer close(out)
	data := dataI.(DBData)

	query := make(chan string)

	go getCode(database, query, data)

	exist := <-query 
	if exist == "" {
		out<- "400 Кода не существует"
		return
	}

	tx, err := database.Begin()

	if err != nil {
		out<- "500 Не получилось сформировать транзакцию"
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			out<- "500 Ошибка удаления"
		} else {
			tx.Commit()
			out<- "200 Код удален"
		}
	}()

	_, err = tx.Exec(` DELETE FROM code WHERE "code_id"=(SELECT "ref_code" FROM users WHERE "email"=$1) `, data.UserEmail)
	
	if err != nil {
		out<- err.Error()
		return
	}
		
	_, err = tx.Exec(` UPDATE users SET "ref_code"=NULL WHERE "email"=$1 `, data.UserEmail)

	if err != nil {
		out<- err.Error()
		return
	}

	_, err = tx.Exec(` UPDATE users SET "ref_id"=NULL WHERE "ref_id"=(SELECT "user_id" FROM users WHERE "email"=$1) `, data.UserEmail)
	
	if err != nil {
		out<- err.Error()
		return
	}
}

func checkCode(database *sql.DB, out chan string, dataI interface{}) {
	defer close(out)
	data := dataI.(DBData)

	var user_id string
	database.QueryRow(` SELECT "user_id" FROM users WHERE "ref_code"=( SELECT "code_id" FROM code WHERE "code_string"=$1); `, data.RefString).Scan(&user_id)

	out<- user_id
}

func getRefBD(database *sql.DB, out chan string, dataI interface{}) {
	defer close(out)
	data := dataI.(DBData)

	rows, err := database.Query(` SELECT "user_id" FROM users WHERE "ref_id"=$1 `, data.UserId)

	if err != nil {
		out<- err.Error()
		return
	}

	for rows.Next() {
		var id string
		err = rows.Scan(&id)

		if err != nil {
			out<- err.Error()
			return
		}

		out<- id
	}
}