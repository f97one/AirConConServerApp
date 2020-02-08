package dataaccess

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
)

var db *sqlx.DB

const dialect string = "sqlite3"

func init() {
	initDatabase()
}

func initDatabase() {
	database, err := sqlx.Open(dialect, "file:airconcon.sqlite?_mode=rw&_journal=WAL&_auto_vacuu=incremental")
	if err != nil {
		panic(err)
	}

	// マイグレーションを自動実行
	migration := &migrate.FileMigrationSource{Dir: "migrations/sqlite3"}
	_, err = migrate.Exec(database.DB, dialect, migration, migrate.Up)
	if err != nil {
		panic(err)
	}
	db = database
}

// trueを1に、falseを0に変換する。
func boolToInt(b bool) int {
	ret := 0
	if b {
		ret = 1
	}

	return ret
}
