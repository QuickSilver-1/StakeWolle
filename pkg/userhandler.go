package pkg

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/QuickSilver-1/StakeWolle/internal/config"
	_ "github.com/lib/pq"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	config.PgHost, config.PgPort, config.PgUser, config.PgPass, сonfig.PgName)
    db, err := sql.Open("postgres", psqlInfo)
	fmt.Print(db, err)
}
