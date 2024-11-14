package db

import (
	"database/sql"
	"fmt"
	"time"

	"referal/internal/config"

	_ "github.com/lib/pq"
)

var(
	DB = NewDB(*config.AppConfig.PgHost, *config.AppConfig.PgUser, *config.AppConfig.PgPass, *config.AppConfig.PgName, *config.AppConfig.PgPort)
)


type ConnectDatabase struct {
	Quit 		chan interface{}
	Connection	*sql.DB
	Command 	map[string]func(*sql.DB, chan string, ...any)
}

func (c *ConnectDatabase) Exec(comm string, out chan string, args ...any) {
	c.Command[comm](c.Connection, out, args...)
}

func (c *ConnectDatabase) Query(comm string, out chan string, args ...any) {
	c.Command[comm](c.Connection, out, args...)
}

func NewDB(host, user, password, dbname string, port int) ConnectDatabase {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
							host, port, user, password, dbname)
	conn, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic("Проблема с подключением к базе данных")
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(time.Second * 3)

	db := ConnectDatabase{
		Quit: make(chan interface{}),
		Connection: conn,
	}

	go func() {
		<-db.Quit
		db.Connection.Close()
	}()

	return db
}