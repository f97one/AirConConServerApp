package dataaccess

import (
	"database/sql"
	"github.com/f97one/AirConCon/utils"
)

func LoadByUsername(username string) (AppUser, error) {
	sqlStmt := "select user_id, username, password, need_pw_change from app_user where username = $1"
	au := AppUser{}
	err := db.QueryRowx(sqlStmt, username).StructScan(&au)
	if err != nil && err != sql.ErrNoRows {
		utils.GetLogger().Errorln(err)
	}

	return au, err
}

func findById(userId int) (AppUser, error) {
	sqlStmt := "select user_id, username, password, need_pw_change from app_user where user_id = $1"
	au := AppUser{}
	err := db.QueryRowx(sqlStmt, userId).StructScan(&au)
	if err != nil && err != sql.ErrNoRows {
		utils.GetLogger().Errorln(err)
	}

	return au, err
}
