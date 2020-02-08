package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"reflect"
	"runtime"
)

func configureRouting(mux *httprouter.Router) {
	// ログイン
	mux.POST("/login", withLog(login))
	// ログアウト
	mux.POST("/logout", requireJwtHandler(withLog(logout)))
	// JWT更新
	mux.POST("/auth", requireJwtHandler(withLog(auth)))
	// ユーザー追加
	mux.PUT("/adduser", requireJwtHandler(withLog(subscribe)))
	// パスワード更新
	mux.POST("/passwd", requireJwtHandler(withLog(changePassword)))
	// ユーザー削除
	mux.DELETE("/unsubscribe", requireJwtHandler(withLog(unsubscribe)))
}

// withLog sends log to logger before calling Handle
func withLog(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handlerName := runtime.FuncForPC(reflect.ValueOf(handle).Pointer()).Name()
		logger.Tracef("On entering handler %s, with params = %s", handlerName, ps)

		handle(w, r, ps)
	}
}
