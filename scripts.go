package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/f97one/AirConCon/dataaccess"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

// すべてのスクリプトを返す。
func allScripts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	scripts, err := dataaccess.GetAllScripts()
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(scripts)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	if len(scripts) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
}

func addScript(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	sc := &dataaccess.Scripts{}
	err = json.Unmarshal(body, sc)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	s, err := sc.Save()
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(s)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
}

// 番号指定でスクリプトを取得する。
func getScript(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	scriptId := ps.ByName("scriptId")
	if len(scriptId) == 0 {
		logger.Errorln("スクリプト番号を渡されなかった")
		respondError(&w, nil, http.StatusBadRequest)
		return
	}

	script, err := dataaccess.GetScript(scriptId)
	if err != nil {
		logger.Errorln(err)
		if err == sql.ErrNoRows {
			respondError(&w, err, http.StatusNotFound)
		} else {
			respondError(&w, err, http.StatusInternalServerError)
		}
		return
	}
	b, err := json.Marshal(script)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
}

// 指定番号のスクリプトを更新する。
func updateScript(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	scriptId := ps.ByName("scriptId")
	if len(scriptId) == 0 {
		logger.Errorln("スクリプト番号が渡されなかった")
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
	sc := dataaccess.Scripts{}
	err = json.Unmarshal(body, &sc)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

	// 要求データのスクリプト名をもつデータがないかを確認
	existingSc1, err := dataaccess.GetScriptByName(sc.ScriptName)
	if err == nil {
		if existingSc1.ScriptId != scriptId {
			msg := msgResp{Msg: fmt.Sprintf("given script name %s already exists.", sc.ScriptName)}
			b, err := json.Marshal(msg)
			if err != nil {
				logger.Errorln(err)
				respondError(&w, err, http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusConflict)
			w.Header().Add(contentType, appJson)
			_, err = w.Write(b)
			if err != nil {
				logger.Errorln(err)
				respondError(&w, err, http.StatusInternalServerError)
			}
			return
		}
	} else if err != sql.ErrNoRows {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

	existingSc, err := dataaccess.GetScript(scriptId)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := msgResp{Msg: fmt.Sprintf("given script id %s doesn't exist.", scriptId)}
			b, err := json.Marshal(msg)
			if err != nil {
				logger.Errorln(err)
				respondError(&w, err, http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			w.Header().Add(contentType, appJson)
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

	existingSc.ScriptName = sc.ScriptName
	existingSc.Freq = sc.Freq
	existingSc.Gpio = sc.Gpio
	newSc, err := existingSc.Save()
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(newSc)
	w.Header().Add(contentType, appJson)
	_, err = w.Write(b)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
}

// 指定番号のスクリプトを削除する。
func removeScript(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	scriptId := ps.ByName("scriptId")
	if len(scriptId) == 0 {
		logger.Errorln("スクリプト番号が渡されなかった")
		respondError(&w, nil, http.StatusBadRequest)
		return
	}

	err := dataaccess.DeleteScript(scriptId)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := msgResp{Msg: fmt.Sprintf("script id %s not found.", scriptId)}
			b, err := json.Marshal(msg)
			if err != nil {
				logger.Errorln(err)
				respondError(&w, err, http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			w.Header().Add(contentType, appJson)
			_, err = w.Write(b)
			if err != nil {
				logger.Errorln(err)
				respondError(&w, err, http.StatusInternalServerError)
			}
		} else {
			logger.Errorln(err)
			respondError(&w, err, http.StatusInternalServerError)
		}
		return
	}

	// NO_CONTENT を返す
	w.WriteHeader(http.StatusNoContent)

}
