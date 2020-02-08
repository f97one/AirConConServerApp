package dataaccess

import (
	"errors"
	"github.com/f97one/AirConCon/utils"
	"time"
)

func PutToken(userId int, token string, expiration time.Time) error {
	logger := utils.GetLogger()

	tx, err := db.Beginx()
	if err != nil {
		logger.Errorln(err)
		return err
	}

	bindVals := map[string]interface{}{
		"userId":    userId,
		"token":     token,
		"expiresAt": expiration.Format(time.RFC3339),
	}

	var jt JwtToken
	findSql := "select user_id, generated_token, expires_at from jwt_token where user_id = $1"
	err = tx.QueryRowx(findSql, userId).StructScan(&jt)

	var mergeSql string
	if err != nil {
		// ROWが空エラー = レコードなしなのでINSERT
		mergeSql = "insert into jwt_token (user_id, generated_token, expires_at) values (:userId, :token, :expiresAt)"
	} else {
		// エラーなし = レコードありなのでUPDATE
		mergeSql = "update jwt_token set generated_token = :token, expires_at = :expiresAt where user_id = :userId"
	}
	_, err = db.NamedExec(mergeSql, bindVals)
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

func RemoveToken(userId int) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	sqlStmt := "delete from jwt_token where user_id = :userId"

	bind := map[string]interface{}{
		"userId": userId,
	}

	result, err := tx.NamedExec(sqlStmt, bind)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
