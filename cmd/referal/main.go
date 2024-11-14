package main

import (
	"fmt"
	c "referal/internal/config"
	dh "referal/internal/dbhandlers"
	sb "referal/internal/serverbuilder"
	"referal/pkg/db"
	"referal/pkg/log"
	"strconv"
	"time"
)

func main() {
	close := make(chan interface{})
	go db.ClearDB(time.Hour, db.DB.Connection, "code", close)
	defer func() {
		close<-1
	}()

	dh.CollectHandlers(&db.DB)

	appServer := sb.MakeServer(":" + strconv.Itoa(*c.AppConfig.HttpPort), 10, 10)
	log.Logger.Info("start server")
	err := appServer.ListenAndServe()
	fmt.Print(err)
}