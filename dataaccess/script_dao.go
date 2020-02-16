package dataaccess

import (
	"database/sql"
	"github.com/f97one/AirConCon/utils"
)

func GetScript(scriptId string) (*Scripts, error) {
	sqlStmt := "select script_id, gpio, script_name, freq from scripts where script_id = $1"

	s := &Scripts{}
	err := db.QueryRowx(sqlStmt, scriptId).StructScan(s)
	if err != nil && err != sql.ErrNoRows {
		utils.GetLogger().Errorln(err)
		return nil, err
	}
	return s, err
}

func GetScriptByName(scriptName string) (*Scripts, error) {
	sqlStmt := "select script_id, gpio, script_name, freq from scripts where script_name = $1"

	s := &Scripts{}
	err := db.QueryRowx(sqlStmt, scriptName).StructScan(s)
	if err != nil && err != sql.ErrNoRows {
		utils.GetLogger().Errorln(err)
		return nil, err
	}
	return s, err
}

func (s *Scripts) Save() (*Scripts, error) {
	logger := utils.GetLogger()

	sqlStmt := "update scripts set gpio = :gpio, script_name = :scriptName, freq = :freq where script_id = :scriptId"
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

	_, err := db.NamedExec(sqlStmt, bindValues)
	if err != nil {
		logger.Errorln(err)
	}
	return s, err
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
		s := &Scripts{}
		err = rows.StructScan(s)
		if err != nil {
			logger.Errorln(err)
			return nil, err
		}
		ret = append(ret, *s)
	}
	return ret, nil
}
