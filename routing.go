package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"reflect"
	"runtime"
)

func configureRouting(mux *httprouter.Router) {
	mux.POST("/login", withLog(login))
	mux.POST("/logout", requireJwtHandler(withLog(logout)))
}

// withLog sends log to logger before calling Handle
func withLog(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handlerName := runtime.FuncForPC(reflect.ValueOf(handle).Pointer()).Name()
		logger.Tracef("On entering handler %s, with params = %s", handlerName, ps)

		handle(w, r, ps)
	}
}
