package dbconn

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func GetPostgresDb(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err == nil {
		err = db.Ping()
	}

	return db, err
}

func ClosePostgresDbConn(db *sql.DB) error {
	return db.Close()
}
