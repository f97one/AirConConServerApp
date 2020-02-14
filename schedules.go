package main

import (
	"encoding/json"
	"github.com/f97one/AirConCon/dataaccess"
	"github.com/f97one/AirConCon/utils"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"time"
)

// 保存中のすべてのスケジュールを返す。
func allSchedules(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	schedules, err := dataaccess.GetAllSchedule()
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

	sch := make([]scheduleResp, len(schedules))
	for _, val := range schedules {

		wd := make([]int, 7)
		for _, v := range val.ExecDay {
			wd = append(wd, int(v.WeekdayId))
		}

		s := scheduleResp{
			ScheduleId: val.ScheduleId,
			Name:       val.Name,
			OnOff:      utils.BoolToOnOff(val.OnOff),
			Weekday:    wd,
			Time:       val.ExecuteTime,
			ScriptId:   val.ScriptId,
		}
		sch = append(sch, s)
	}

	if len(sch) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	b, err := json.Marshal(sch)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
	}
}

func addSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	schResp := &scheduleResp{}
	err = json.Unmarshal(body, schResp)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	var weekdayTiming []dataaccess.Timing
	for _, weekday := range schResp.Weekday {
		weekdayTiming = append(weekdayTiming, dataaccess.Timing{
			ScheduleId: schResp.ScheduleId,
			WeekdayId:  time.Weekday(weekday),
		})
	}

	sch := dataaccess.Schedule{
		ScheduleId:  schResp.ScheduleId,
		Name:        schResp.Name,
		OnOff:       utils.OnOffToBool(schResp.OnOff),
		ExecuteTime: schResp.Time,
		ScriptId:    schResp.ScriptId,
		ExecDay:     weekdayTiming,
	}

	err = sch.Save()
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	schResp.ScheduleId = sch.ScheduleId
	b, err := json.Marshal(schResp)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(b)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
}
