package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/f97one/AirConCon/dataaccess"
	"github.com/f97one/AirConCon/utils"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	queueId             string = "n"
	scriptExecutionFile string = "airconcon_cmd.sh"
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

	b, err := json.Marshal(sch)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add(contentType, appJson)
	if len(sch) == 0 {
		w.WriteHeader(http.StatusNotFound)
		b, err = json.Marshal(msgResp{Msg: "Schedule not found."})
		if err != nil {
			logger.Errorln(err)
			respondError(&w, err, http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusOK)
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
	// scheduleRespのバリデータでエラーを400で返す
	err = schResp.validate()
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusBadRequest)
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

	if scheduleId == "next" {
		nextSchedule(w, r, ps)
		return
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
	// scheduleRespのバリデータでエラーを400で返す
	err = schReq.validate()
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusBadRequest)
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

// 次回スケジュールを返す。
func nextSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Add(contentType, appJson)

	next, err := dataaccess.GetNextSchedule()
	if err != nil {
		logger.Errorln(err)
		sc := http.StatusInternalServerError
		if err == sql.ErrNoRows {
			sc = http.StatusNotFound
		}
		w.WriteHeader(sc)
		b, err := json.Marshal(msgResp{Msg: err.Error()})
		if err != nil {
			logger.Error(err)
			respondError(&w, err, http.StatusInternalServerError)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			logger.Error(err)
			respondError(&w, err, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	nextResp := nextScheduleResp{
		ScheduleId: next.ScheduleId,
		Name:       next.Name,
		OnOff:      utils.BoolToOnOff(next.OnOff),
		WeekdayId:  next.WeekdayId,
		Time:       next.ExecuteTime,
		ScriptId:   next.ScriptId,
	}

	b, err := json.Marshal(nextResp)
	if err != nil {
		logger.Error(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		logger.Error(err)
		respondError(&w, err, http.StatusInternalServerError)
	}
}

// 次回予定をシステムに登録する。
func registerNext(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if checkPlatform(w) {
		return
	}
	logger.Tracef("%s で実行", runtime.GOOS)

	// 次回分の実行情報を得る
	next, err := dataaccess.GetNextSchedule()
	if err != nil {
		logger.Error(err)
		if err == sql.ErrNoRows {
			b, err := json.Marshal(msgResp{Msg: err.Error()})
			if err != nil {
				logger.Error(err)
				respondError(&w, err, http.StatusInternalServerError)
				return
			}
			w.Header().Add(contentType, appJson)
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write(b)
			if err != nil {
				logger.Error(err)
				respondError(&w, err, http.StatusInternalServerError)
			}
			return
		}
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	logger.Traceln(next)

	script, err := dataaccess.GetScript(next.ScriptId)
	if err != nil {
		logger.Error(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	currentTime := time.Now()
	deltas := next.WeekdayId - int(currentTime.Weekday())
	if deltas < 0 {
		deltas = deltas + 7
	}

	// 実行時間のパース
	hhmm := strings.Split(next.ExecuteTime, ":")
	hh, err := strconv.Atoi(hhmm[0])
	if err != nil {
		logger.Error(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	mm, err := strconv.Atoi(hhmm[1])
	if err != nil {
		logger.Error(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	y, m, d := currentTime.AddDate(0, 0, deltas).Date()
	// at の -t オプションの書式に合わせる
	runAt := time.Date(y, m, d, hh, mm, 0, 0, time.Local).Format("200601021504.05")

	// テンプレートへのバインド値の作成
	bindValues := map[string]interface{}{
		"Freq":         script.Freq,
		"Gpio":         script.Gpio,
		"SignalDbFile": conf.SignalDbFile,
		"ScriptName":   script.ScriptName,
	}
	logger.Traceln(bindValues)

	cmdTmpl, err := template.ParseFiles("cmdtmpl/playback.txt")
	if err != nil {
		logger.Error(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

	file, err := os.Create("airconcon_cmd.sh")
	if err != nil {
		logger.Error(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

	err = cmdTmpl.Execute(file, bindValues)
	if err != nil {
		logger.Error(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

	// at の登録
	cmd := exec.Command("at", "-M", "-f", scriptExecutionFile, "-t", runAt, "-q", queueId)
	cmdline := cmd.String()
	logger.Tracef("登録される cmdline : %s", cmdline)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := msgResp{Msg: err.Error()}
		b, err := json.Marshal(msg)
		if err != nil {
			logger.Error(err)
			respondError(&w, err, http.StatusInternalServerError)
			return
		}
		w.Header().Add(contentType, appJson)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write(b)
		if err != nil {
			logger.Error(err)
			respondError(&w, err, http.StatusInternalServerError)
		}
		return
	}
	// ジョブIDを取り出す
	logger.Tracef("実行後のStdin : %s", out)
	lines := strings.Split(string(out), "\n")
	var retLine []string
	for _, line := range lines {
		if strings.HasPrefix(line, "job ") {
			logger.Tracef("at(1)の出力 : %s", line)
			retLine = strings.Split(line, " ")
			break
		}
	}
	if len(retLine) == 0 {
		msg := msgResp{Msg: "cmd execution seems to fail."}
		b, err := json.Marshal(msg)
		if err != nil {
			logger.Error(err)
			respondError(&w, err, http.StatusInternalServerError)
			return
		}
		w.Header().Add(contentType, appJson)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write(b)
		if err != nil {
			logger.Error(err)
			respondError(&w, err, http.StatusInternalServerError)
		}
		return
	}
	jobId, err := strconv.Atoi(retLine[1])
	if err != nil {
		logger.Error(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

	jobSchedule := dataaccess.JobSchedule{
		ScheduleId: next.ScheduleId,
		JobId:      jobId,
		CmdLine:    cmdline,
		RunAt:      runAt,
	}
	err = jobSchedule.Save()
	if err != nil {
		logger.Error(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func checkPlatform(w http.ResponseWriter) bool {
	if !(runtime.GOOS == "linux" || runtime.GOOS == "freebsd") {
		body := fmt.Sprintf("This endpoint does not work on this platform (%s).", runtime.GOOS)
		b, err := json.Marshal(msgResp{Msg: body})
		if err != nil {
			logger.Error(err)
			respondError(&w, err, http.StatusInternalServerError)
			return true
		}
		w.Header().Add(contentType, appJson)
		w.WriteHeader(http.StatusNotImplemented)
		_, err = w.Write(b)
		if err != nil {
			logger.Error(err)
			respondError(&w, err, http.StatusInternalServerError)
		}
		return true
	}
	return false
}

// 次回予定をキャンセルする。
func cancelNext(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if checkPlatform(w) {
		return
	}
	logger.Tracef("%s で実行", runtime.GOOS)

	atq := exec.Command("at", "-l", "-q", queueId)
	logger.Tracef("実行する cmdline : %s", atq.String())

	out, err := atq.CombinedOutput()
	if err != nil {
		msg := msgResp{Msg: err.Error()}
		b, err := json.Marshal(msg)
		if err != nil {
			logger.Error(err)
			respondError(&w, err, http.StatusInternalServerError)
			return
		}
		w.Header().Add(contentType, appJson)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write(b)
		if err != nil {
			logger.Error(err)
			respondError(&w, err, http.StatusInternalServerError)
		}
		return
	}
	logger.Tracef("実行結果 : %s", out)

	lines := strings.Split(string(out), "\n")

	var jobIds []string
	for _, line := range lines {
		// 先頭８文字だけ切り取る
		logger.Tracef("行データ : %s", line)
		idx := strings.Index(line, "\t")
		logger.Tracef("index : %d", idx)
		if idx == -1 {
			continue
		}
		jobNo := line[:idx]
		jobIds = append(jobIds, jobNo)
	}

	jsSlice, err := dataaccess.GetRegisteredNext(jobIds)
	if err != nil {
		logger.Error(err)
		if err == sql.ErrNoRows {
			msg := msgResp{Msg: err.Error()}
			b, err := json.Marshal(msg)
			if err != nil {
				logger.Error(err)
				respondError(&w, err, http.StatusInternalServerError)
				return
			}
			w.Header().Add(contentType, appJson)
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write(b)
			if err != nil {
				logger.Error(err)
				respondError(&w, err, http.StatusInternalServerError)
			}
			return
		}
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

	var ids string
	var removeJobId []string
	for jIdx, js := range jsSlice {
		removeJobId = append(removeJobId, strconv.Itoa(js.JobId))
		if jIdx == 0 {
			ids = fmt.Sprintf("%d", js.JobId)
		} else {
			ids = fmt.Sprintf("%s %d ", ids, js.JobId)
		}
	}

	atrm := exec.Command("at", "-r", ids)
	logger.Tracef("実行する cmdline : %s", atrm.String())
	out, err = atrm.CombinedOutput()
	logger.Tracef("実行結果 : %s", out)
	if err != nil {
		logger.Error(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

	// ジョブスケジュールの消込
	err = dataaccess.RemoveNext(removeJobId)
	if err != nil {
		logger.Error(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
