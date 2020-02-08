package dataaccess

import (
	"database/sql"
	"github.com/f97one/AirConCon/utils"
)

func LoadByUsername(username string) (AppUser, error) {
	sqlStmt := "select user_id, username, password, need_pw_change, admin_flag from app_user where username = $1"
	au := AppUser{}
	err := db.QueryRowx(sqlStmt, username).StructScan(&au)
	if err != nil && err != sql.ErrNoRows {
		utils.GetLogger().Errorln(err)
	}

	return au, err
}

func findById(userId int) (AppUser, error) {
	sqlStmt := "select user_id, username, password, need_pw_change, admin_flag from app_user where user_id = $1"
	au := AppUser{}
	err := db.QueryRowx(sqlStmt, userId).StructScan(&au)
	if err != nil && err != sql.ErrNoRows {
		utils.GetLogger().Errorln(err)
	}

	return au, err
}

func CreateUser(user AppUser) error {
	sqlStmt := "insert into app_user (username, password, need_pw_change, admin_flag) values (:username, :password, :needPwChange, :adminFlag)"
	logger := utils.GetLogger()

	tx, err := db.Beginx()
	if err != nil {
		logger.Errorln(err)
		return err
	}

	bindValues := map[string]interface{}{
		"username":     user.Username,
		"password":     user.Password,
		"needPwChange": boolToInt(user.NeedPwChange),
		"adminFlag":    boolToInt(user.AdminFlag),
	}

	_, err = tx.NamedExec(sqlStmt, bindValues)
	if err != nil {
		logger.Errorln(err)
		err = tx.Rollback()
		if err != nil {
			logger.Errorln(err)
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		logger.Errorln(err)
		return err
	}
	return nil
}
