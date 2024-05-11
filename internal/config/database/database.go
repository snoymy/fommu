package database

import (
	"app/internal/config"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewConnection() *sqlx.DB {
    dburl := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", config.DB.DBUser, config.DB.DBPass, config.DB.DBHost, config.DB.DBPort, config.DB.DBName)

    db, err := sqlx.Connect("pgx", dburl)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}

    return db 
}

func TestConnection(db *sqlx.DB) error {
    result := 0
    err := db.Get(&result, "select 1")
    if err != nil {
        return err 
    }

    if result != 1 {
        return errors.New("Result unexpected, expect 1")
    }

    return nil
}
