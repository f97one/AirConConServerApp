package dataaccess

import (
	"database/sql"
	"errors"
	"github.com/f97one/AirConCon/utils"
)

func GetScript(scriptId string) (*Scripts, error) {
	sqlStmt := "select script_id, gpio, script_name, freq from scripts where script_id = $1"

	var s *Scripts
	err := db.QueryRowx(sqlStmt, scriptId).StructScan(s)
	if err != nil && err != sql.ErrNoRows {
		utils.GetLogger().Errorln(err)
		return nil, err
	}
	return s, err
}

func GetScriptByName(scriptName string) (*Scripts, error) {
	sqlStmt := "select script_id, gpio, script_name, freq from scripts where script_name = $1"

	var s *Scripts
	err := db.QueryRowx(sqlStmt, scriptName).StructScan(s)
	if err != nil && err != sql.ErrNoRows {
		utils.GetLogger().Errorln(err)
		return nil, err
	}
	return s, err
}

func (s *Scripts) Save() error {
	logger := utils.GetLogger()

	// 衝突するレコードがないか調べる
	sc, err := GetScriptByName(s.ScriptName)
	if err != nil && err != sql.ErrNoRows {
		logger.Errorln(err)
		return err
	}
	if sc != nil {
		err = errors.New("same script name found, abort")
		return err
	}

	sqlStmt := "update scripts set gpio = :gpio, script_name = :scriptId, freq = :freq where script_id = :scriptId"
	if len(s.ScriptId) == 0 {
		sqlStmt = "insert into scripts (script_id, gpio, script_name, freq) values (:scriptId, :gpio, :scriptName, :freq)"
		s.ScriptId = createKey()
	}

	bindValues := map[string]interface{}{
		"scriptId":   s.ScriptId,
		"gpio":       s.Gpio,
		"scriptName": s.ScriptName,
		"freq":       s.Freq,
	}

	_, err = db.NamedExec(sqlStmt, bindValues)
	if err != nil {
		logger.Errorln(err)
	}
	return err
}

func GetAllScripts() ([]Scripts, error) {
	logger := utils.GetLogger()

	sqlStmt := "select script_id, gpio, script_name, freq from scripts"
	ret := make([]Scripts, 0)
	rows, err := db.Queryx(sqlStmt)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil && err != sql.ErrNoRows {
		logger.Errorln(err)
		return nil, err
	}
	for rows.Next() {
		var s Scripts
		err = rows.StructScan(s)
		if err != nil {
			logger.Errorln(err)
			return nil, err
		}
		ret = append(ret, s)
	}
	return ret, nil
}
