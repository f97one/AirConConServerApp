package dataaccess

import (
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
