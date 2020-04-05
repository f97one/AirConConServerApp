package dataaccess

import (
	"crypto/sha1"
	"fmt"
	"github.com/f97one/AirConCon/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
	"io"
	"os"
	"path/filepath"
	"time"
)

var db *sqlx.DB

const (
	dialect    string = "sqlite3"
	dbFilename string = "airconcon.sqlite"
)

func InitDatabase() {
	log := utils.GetLogger()

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	dbFilePath := filepath.Join(currentDir, dbFilename)
	log.Tracef("DBgファイルのパス : %s", dbFilePath)
	database, err := sqlx.Open(dialect, fmt.Sprintf("file:%s?_mode=rw&_journal=WAL&_auto_vacuum=incremental", dbFilePath))
	if err != nil {
		log.Fatalln(err)
	}

	// マイグレーションを自動実行
	migration := &migrate.FileMigrationSource{Dir: "migrations/sqlite3"}
	_, err = migrate.Exec(database.DB, dialect, migration, migrate.Up)
	if err != nil {
		log.Fatalln(err)
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
