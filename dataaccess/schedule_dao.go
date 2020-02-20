package dataaccess

import (
	"database/sql"
	"fmt"
	"github.com/f97one/AirConCon/utils"
	"time"
)

func GetAllSchedule() ([]Schedule, error) {
	logger := utils.GetLogger()

	schStmt := "select schedule_id, name, on_off, execute_time, script_id from schedule"
	var ret []Schedule
	schRows, err := db.Queryx(schStmt)
	if schRows != nil {
		defer schRows.Close()
	}
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	for schRows.Next() {
		s := &Schedule{}
		err = schRows.StructScan(s)
		if err != nil {
			logger.Errorln(err)
			return nil, err
		}
		ret = append(ret, *s)
	}

	timingStmt := "select schedule_id, weekday_id from timing where schedule_id = $1"
	if ret != nil {
		for idx, schVal := range ret {
			var timingResult []Timing
			timingRows, err := db.Queryx(timingStmt, schVal.ScheduleId)
			if err != nil {
				logger.Errorln(err)
				return nil, err
			}

			for timingRows.Next() {
				var t Timing
				err = timingRows.StructScan(&t)
				if err != nil {
					logger.Errorln(err)
					return nil, err
				}
				timingResult = append(timingResult, t)
			}
			ret[idx].ExecDay = timingResult
			_ = timingRows.Close()
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

func GetSchedule(scheduleId string) (*Schedule, error) {
	logger := utils.GetLogger()
	schStmt := "select schedule_id, name, on_off, execute_time, script_id from schedule where schedule_id = $1"
	var s Schedule
	err := db.QueryRowx(schStmt, scheduleId).StructScan(&s)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}

	timingStmt := "select schedule_id, weekday_id from timing where schedule_id = $1"
	rows, err := db.Queryx(timingStmt, scheduleId)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	defer rows.Close()

	var timings []Timing
	for rows.Next() {
		var t Timing
		err = rows.StructScan(&t)
		if err != nil {
			logger.Errorln(err)
			return nil, err
		}
		timings = append(timings, t)
	}

	s.ExecDay = timings

	return &s, nil
}

func DeleteSchedule(scheduleId string) error {
	logger := utils.GetLogger()
	_, err := GetSchedule(scheduleId)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	tx, err := db.Beginx()
	if err != nil {
		logger.Errorln(err)
		return err
	}

	timingStmt := "delete from timing where schedule_id = $1"
	_, err = tx.Exec(timingStmt, scheduleId)
	if err != nil {
		logger.Errorln(err)
		err = tx.Rollback()
		if err != nil {
			logger.Errorln(err)
		}
		return err
	}

	schStmt := "delete from schedule where schedule_id = $1"
	_, err = tx.Exec(schStmt, scheduleId)
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
		err = tx.Rollback()
		if err != nil {
			logger.Errorln(err)
		}
		return err
	}
	return nil
}

// 次回スケジュールを取得する。
func GetNextSchedule() (*NextSchedule, error) {
	logger := utils.GetLogger()

	currentTime := time.Now().Format("15:04")
	dayOfWeek := time.Now().Weekday()

	bindValues := map[string]interface{}{
		"currentTime": currentTime,
		"dayOfWeek":   int(dayOfWeek),
	}

	sqlFmt := `select sc.schedule_id, sc.name, sc.on_off, sc.execute_time, sc.script_id, tm.weekday_id from schedule sc
inner join timing tm on sc.schedule_id = tm.schedule_id
where sc.execute_time %s :currentTime
and tm.weekday_id %s :dayOfWeek
order by sc.execute_time, tm.weekday_id
limit 1`
	sql1 := fmt.Sprintf(sqlFmt, ">", ">=")
	sql2 := fmt.Sprintf(sqlFmt, ">=", ">")
	rows, err := db.NamedQuery(sql1, bindValues)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warnln(err)
			rows, err = db.NamedQuery(sql2, bindValues)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	ret := &NextSchedule{}
	_ = rows.Next()
	err = rows.StructScan(ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
