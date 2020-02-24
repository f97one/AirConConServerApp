package dataaccess

import (
	"database/sql"
	"github.com/f97one/AirConCon/utils"
)

func (js *JobSchedule) Save() error {
	logger := utils.GetLogger()
	_, err := GetCondition(js.ScheduleId)
	if err != nil && err != sql.ErrNoRows {
		logger.Errorln(err)
		return err
	}

	stmt := "update job_schedule set job_id = :jobId, cmd_line = :cmdLine, run_at = :runAt where schedule_id = :scheduleId"
	if err != nil {
		stmt = "insert into job_schedule (schedule_id, job_id, cmd_line, run_at) values (:scheduleId, :jobId, :cmdLine, :runAt)"
	}

	tx, err := db.Beginx()
	if err != nil {
		logger.Errorln(err)
		return err
	}

	bindValues := map[string]interface{}{
		"scheduleId": js.ScheduleId,
		"jobId":      js.JobId,
		"cmdLine":    js.CmdLine,
		"runAt":      js.RunAt,
	}

	_, err = tx.NamedExec(stmt, bindValues)
	if err != nil {
		logger.Errorln(err)
		err = tx.Rollback()
		if err != nil {
			logger.Errorln(err)
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		logger.Errorln(err)
	}

	return nil
}

func GetCondition(scheduleId string) (*JobSchedule, error) {
	stmt := "select schedule_id, job_id, cmd_line, run_at from job_schedule where schedule_id = $1"
	var js JobSchedule
	err := db.QueryRowx(stmt, scheduleId).StructScan(&js)
	if err != nil {
		return nil, err
	}
	return &js, nil
}
