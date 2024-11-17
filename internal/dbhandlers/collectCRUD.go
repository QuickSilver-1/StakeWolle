package dbhandlers

import (
	"database/sql"

	"referal/pkg/db"
	"referal/pkg/log"
)

// CollectHandlers инициализирует обработчики CRUD операций
func CollectHandlers(conn *db.ConnectDatabase) {
    log.Logger.Info("Сборка CRUD")

    // Инициализация команд в структуре ConnectDatabase
    conn.Command = map[string]func(*sql.DB, chan string, interface{}){
        "create":   createUser,    // Обработчик для создания пользователя
        "check":    checkUser,     // Обработчик для проверки пользователя для входа
        "generate": setCode,       // Обработчик для генерации реферального кода
        "get":      getCode,       // Обработчик для получения кода по email
        "delete":   delCode,       // Обработчик для удаления реферального кода
        "referals": getRefBD,      // Обработчик для получения рефералов по id
    }
}
