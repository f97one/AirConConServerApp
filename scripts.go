package main

import (
	"encoding/json"
	"github.com/f97one/AirConCon/dataaccess"
	"github.com/julienschmidt/httprouter"
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
	_, err = w.Write(b)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
}

func addScript(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
