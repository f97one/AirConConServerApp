package dataaccess

import (
	"database/sql"
	"github.com/f97one/AirConCon/utils"
)

func GetAllSchedule() ([]Schedule, error) {
	logger := utils.GetLogger()

	schStmt := "select schedule_id, name, on_off, execute_time, script_id from schedule"
	var ret []Schedule
	rows, err := db.Queryx(schStmt)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil && err != sql.ErrNoRows {
		logger.Errorln(err)
		return nil, err
	}
	if rows != nil {
		for rows.Next() {
			s := &Schedule{}
			err = rows.StructScan(s)
			if err != nil {
				logger.Errorln(err)
				return nil, err
			}
			ret = append(ret, *s)
		}
	}

	timingStmt := "select schedule_id, weekday_id from timing where schedule_id = $1"
	if ret != nil {
		for _, val := range ret {
			timingResult := make([]Timing, 7)
			rows, err = db.Queryx(timingStmt, val.ScheduleId)
			if err != nil {
				logger.Errorln(err)
				return nil, err
			}
			if rows != nil {
				for rows.Next() {
					t := &Timing{}
					err = rows.StructScan(t)
					if err != nil {
						logger.Errorln(err)
						return nil, err
					}
					timingResult = append(timingResult, *t)
				}
				val.ExecDay = timingResult
			}
			_ = rows.Close()
		}
	}

	return ret, nil
}

func (s *Schedule) Save() error {
	logger := utils.GetLogger()

	tx, err := db.Beginx()
	if err != nil {
		logger.Errorln(err)
		return err
	}

	if len(s.ScheduleId) > 0 {
		// timing の関連レコードを 先にdelete
		delStmt := "delete from timing where schedule_id = $1"

		_, err = tx.Exec(delStmt, s.ScheduleId)
		if err != nil {
			logger.Errorln(err)
			_ = tx.Rollback()
			return err
		}
	}

	schStmt := `
update schedule set name = :name, on_off = :onOff, execute_time = :executeTime, script_id = :scriptId 
where script_id = :scriptId
`
	if s.ScheduleId == "" {
		schStmt = `
insert into schedule (schedule_id, name, on_off, execute_time, script_id) 
values (:scheduleId, :name, :onOff, :executeTime, :scriptId)
`
		s.ScheduleId = createKey()
	}

	bindValues := map[string]interface{}{
		"scheduleId":  s.ScheduleId,
		"name":        s.Name,
		"onOff":       s.OnOff,
		"executeTime": s.ExecuteTime,
		"scriptId":    s.ScriptId,
	}

	_, err = tx.NamedExec(schStmt, bindValues)
	if err != nil {
		logger.Errorln(err)
		_ = tx.Rollback()
		return err
	}

	insStmt := "insert into timing (schedule_id, weekday_id) values ($1, $2)"
	for _, timing := range s.ExecDay {
		timing.ScheduleId = s.ScheduleId
		_, err = tx.Exec(insStmt, s.ScheduleId, timing.WeekdayId)
		if err != nil {
			logger.Errorln(err)
			_ = tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		logger.Errorln(err)
		_ = tx.Rollback()
		return err
	}
	return nil
}
