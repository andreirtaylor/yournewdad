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

func saveGame(g *GameStartRequest, req *http.Request) {
	ctx := appengine.NewContext(req)

	db, err := getDB(req)
	if err != nil {
		log.Errorf(ctx, "Could not get DB %v", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO Games(GameId, Width, Height) VALUES(?,?,?)")

	if err != nil {
		log.Errorf(ctx, "Unable to prepare game saving statement: %v", err)
	}
	_, err = stmt.Exec(g.GameId, g.Width, g.Height)
	if err != nil {
		log.Errorf(ctx, "Error executing game save statement: %v", err)
	}
}

func SaveMove(bo *MoveRequest, req *http.Request) {
	//	lastId, err := res.LastInsertId()
	//	if err != nil {
	//		log.Errorf(ctx, "Could not get lask %v", err)
	//	}
	//	rowCnt, err := res.RowsAffected()
	//	if err != nil {
	//		log.Errorf(ctx, "Could not get lask %v", err)
	//		log.Fatal(err)
	//	}
	//	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		//log.Errorf("%s environment variable not set.", k)
	}
	return v
}
