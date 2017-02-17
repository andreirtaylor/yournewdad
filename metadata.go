package kaa

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
	"os"
)

// remember to defer db.close
func getDB(req *http.Request) (*sql.DB, error) {
	connectionName := mustGetenv("CLOUDSQL_CONNECTION_NAME")
	user := mustGetenv("CLOUDSQL_USER")
	password := mustGetenv("CLOUDSQL_PASSWORD")

	sqlStr := fmt.Sprintf("%s:%s@cloudsql(%s)/", user, password, connectionName)

	if appengine.IsDevAppServer() {
		sqlStr = "root@/Kaa" // dev server has no password baby
	}

	db, err := sql.Open(
		"mysql",
		sqlStr)

	return db, err
}

func SaveMove(bo *MoveRequest, req *http.Request) {
	ctx := appengine.NewContext(req)

	db, err := getDB(req)
	defer db.Close()
	if err != nil {
		log.Errorf(ctx, "Could not get DB %v", err)
		return
	}
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		//log.Errorf("%s environment variable not set.", k)
	}
	return v
}
