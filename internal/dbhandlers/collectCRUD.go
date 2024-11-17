package dbhandlers

import (
	"database/sql"
	"referal/pkg/db"
)

func CollectHandlers(conn *db.ConnectDatabase) {
	conn.Command = map[string]func(*sql.DB, chan string, interface{}){
		"create":   createUser,
		"check":    checkUser,
		"checkid":  checkUserID,
		"generate": setCode,
		"get":      getCode,
		"delete":   delCode,
		"referals": getRefBD,
	}
}