package log

import (
	"encoding/json"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DBLogWritter struct {
    db *sqlx.DB
}

func NewDBWritter() *DBLogWritter {
    if _, err := os.Stat("log.db"); err != nil {
        os.Create("log.db")
    }

    db, err := sqlx.Open("sqlite3", "log.db")

    if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

    if statement, err := db.Prepare(`
      CREATE TABLE IF NOT EXISTS log (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        time_stamp TIMESTAMPTZ,
        request_id TEXT,
        level TEXT,
        message TEXT,
        attributes TEXT
      ); 
    `); err != nil {
		panic(err)
    } else {
        statement.Exec()
    }

    return &DBLogWritter{db: db}
}

func (w *DBLogWritter) Write(requestId string, level string, msg string, timeStamp time.Time, attrs map[string]any) error {
    a, _ := json.Marshal(attrs)
    _, err := w.db.Exec(`insert into log (time_stamp, request_id, level, message, attributes) values (?,?,?,?,?)`, timeStamp, requestId, level, msg, a)
    return err
}
