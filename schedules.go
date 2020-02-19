package main

import (
	"database/sql"
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

	var sch []scheduleResp
	for _, val := range schedules {

		var wd []int
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

	w.Header().Add(contentType, appJson)
	if len(sch) == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
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
	w.Header().Add(contentType, appJson)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(b)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
}

// 指定スケジュールを取得する。
func getSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	scheduleId := ps.ByName("scheduleId")
	if len(scheduleId) == 0 {
		logger.Errorln("スケジュール番号を渡されなかった")
		respondError(&w, nil, http.StatusBadRequest)
	}

	sch, err := dataaccess.GetSchedule(scheduleId)
	if err != nil {
		logger.Errorln(err)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			respondError(&w, err, http.StatusInternalServerError)
			return
		}
	}
	var weekdays []int
	for _, v := range sch.ExecDay {
		weekdays = append(weekdays, int(v.WeekdayId))
	}

	ret := scheduleResp{
		ScheduleId: sch.ScheduleId,
		Name:       sch.Name,
		OnOff:      utils.BoolToOnOff(sch.OnOff),
		Weekday:    weekdays,
		Time:       sch.ExecuteTime,
		ScriptId:   sch.ScriptId,
	}

	b, err := json.Marshal(ret)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Add(contentType, appJson)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
}

// スケジュールを更新する。
func updateSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	scheduleId := ps.ByName("scheduleId")
	if len(scheduleId) == 0 {
		logger.Errorln("スケジュール番号を渡されなかった")
		respondError(&w, nil, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	schReq := &scheduleResp{}
	err = json.Unmarshal(body, schReq)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	var weekdayTiming []dataaccess.Timing
	for _, weekday := range schReq.Weekday {
		weekdayTiming = append(weekdayTiming, dataaccess.Timing{
			ScheduleId: schReq.ScheduleId,
			WeekdayId:  time.Weekday(weekday),
		})
	}

	sch, err := dataaccess.GetSchedule(scheduleId)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Errorln(err)
			w.Header().Add(contentType, appJson)
			w.WriteHeader(http.StatusNotFound)
			msg := msgResp{Msg: err.Error()}
			b, err := json.Marshal(msg)
			if err != nil {
				logger.Errorln(err)
				respondError(&w, err, http.StatusInternalServerError)
				return
			}
			_, err = w.Write(b)
			if err != nil {
				logger.Errorln(err)
				respondError(&w, err, http.StatusInternalServerError)
				return
			}
		}
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	sch.ExecDay = weekdayTiming
	sch.ExecuteTime = schReq.Time
	sch.OnOff = utils.OnOffToBool(schReq.OnOff)
	sch.Name = schReq.Name
	sch.ScriptId = schReq.ScriptId

	err = sch.Save()
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

	schResp := scheduleResp{
		ScheduleId: sch.ScriptId,
		Name:       sch.Name,
		OnOff:      utils.BoolToOnOff(sch.OnOff),
		Weekday:    schReq.Weekday,
		Time:       sch.ExecuteTime,
		ScriptId:   sch.ScriptId,
	}
	b, err := json.Marshal(schResp)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add(contentType, appJson)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
	}
}

// 指定スケジュールを削除する。
func removeSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	scheduleId := ps.ByName("scheduleId")
	if len(scheduleId) == 0 {
		m := "schedule id must not be empty"
		logger.Errorln(m)
		msg := msgResp{Msg: m}
		b, err := json.Marshal(msg)
		if err != nil {
			logger.Errorln(err)
			respondError(&w, err, http.StatusInternalServerError)
			return
		}
		w.Header().Add(contentType, appJson)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(b)
		if err != nil {
			logger.Errorln(err)
			respondError(&w, err, http.StatusInternalServerError)
		}
		return
	}

	err := dataaccess.DeleteSchedule(scheduleId)
	if err != nil {
		logger.Errorln(err)
		msg := msgResp{Msg: err.Error()}
		b, e := json.Marshal(msg)
		if e != nil {
			logger.Errorln(e)
			respondError(&w, e, http.StatusInternalServerError)
			return
		}
		w.Header().Add(contentType, appJson)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, err := w.Write(b)
		if err != nil {
			logger.Errorln(err)
			respondError(&w, err, http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
