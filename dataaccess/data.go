package dataaccess

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
)

var Db *sqlx.DB

const dialect string = "sqlite3"

func init() {
	Db, err := sqlx.Open(dialect, "file:airconcon.sqlite?_mode=rw&_journal=WAL&_auto_vacuu=incremental")
	if err != nil {
		panic(err)
	}

	// マイグレーションを自動実行
	migration := &migrate.FileMigrationSource{Dir: "migrations/sqlite3"}
	_, err = migrate.Exec(Db.DB, dialect, migration, migrate.Up)
	if err != nil {
		panic(err)
	}
}
