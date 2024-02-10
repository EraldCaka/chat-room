package db

import (
	"database/sql"
	"github.com/EraldCaka/chat-room/util"
	_ "github.com/lib/pq"
)

type database struct {
	db *sql.DB
}

func NewDatabase() (*database, error) {
	db, err := sql.Open("postgres", util.DB_CONN_STR)
	if err != nil {
		return nil, err
	}

	return &database{db: db}, nil
}

func (d *database) GetDB() *sql.DB {
	return d.db
}

func (d *database) Ping() error {
	return d.db.Ping()
}

func (d *database) Close() {
	d.db.Close()
}
