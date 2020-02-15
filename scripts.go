package main

import (
	"encoding/json"
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
