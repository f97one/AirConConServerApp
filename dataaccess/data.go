package dataaccess

import (
	"crypto/sha1"
	"fmt"
	"github.com/f97one/AirConCon/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
	"io"
	"time"
)

var db *sqlx.DB

const dialect string = "sqlite3"

func init() {
	initDatabase()
}

func initDatabase() {
	database, err := sqlx.Open(dialect, "file:airconcon.sqlite?_mode=rw&_journal=WAL&_auto_vacuum=incremental")
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

func createKey() string {
	hash := sha1.New()
	_, err := io.WriteString(hash, time.Now().String())
	if err != nil {
		utils.GetLogger().Errorln(err)
		return ""
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}
