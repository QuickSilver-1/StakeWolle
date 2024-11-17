package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"referal/internal/config"
	dh "referal/internal/dbhandlers"
	sb "referal/internal/serverbuilder"
	"referal/pkg/db"
	"referal/pkg/log"
)

func main() {
    closer := make(chan interface{})
    
    // Запуск очистки базы данных каждые час
    go db.ClearDB(time.Hour, db.DB.Connection, "code", closer)
    
    defer func() {
        closer <- 1
        close(closer)
    }()

    // Сборка обработчиков базы данных
    dh.CollectHandlers(&db.DB)

    flag.Parse()

    // Cоздания сервера
    port := strconv.Itoa(*config.AppConfig.HttpPort)
    appServer := sb.MakeServer(":" + port, 10, 10)
    
    log.Logger.Info(fmt.Sprintf("Сервер запущен на порту %s", port))
    
    // Запуск сервера
    err := appServer.ListenAndServe()
    
    if err != nil {
        log.Logger.Error(fmt.Sprintf("Ошибка при запуске сервера: %v", err))
    }
}
