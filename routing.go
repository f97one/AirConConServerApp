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

	// 全スケジュール取得
	mux.GET("/allschedule", requireJwtHandler(withLog(allSchedules)))
	// スケジュール追加
	mux.PUT("/schedule/add", requireJwtHandler(withLog(addSchedule)))
	// 指定スケジュール取得
	mux.GET("/schedule/:scheduleId", requireJwtHandler(withLog(getSchedule)))
	// 指定スケジュール更新
	mux.POST("/schedule/edit/:scheduleId", requireJwtHandler(withLog(updateSchedule)))
	// 指定スケジュール削除
	mux.DELETE("/schedule/drop/:scheduleId", requireJwtHandler(withLog(removeSchedule)))

	// 全スクリプト取得
	mux.GET("/scripts", requireJwtHandler(withLog(allScripts)))
	// スクリプト追加
	mux.PUT("/scripts/add", requireJwtHandler(withLog(addScript)))
	// 指定スクリプト取得
	mux.GET("/scripts/:scriptId", requireJwtHandler(withLog(getScript)))
	// 指定スクリプト更新
	mux.POST("/scripts/edit/:scriptId", requireJwtHandler(withLog(updateScript)))
	// 指定スクリプト削除
	mux.DELETE("/scripts/drop/:scriptId", requireJwtHandler(withLog(removeScript)))
}

// withLog sends log to logger before calling Handle
func withLog(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handlerName := runtime.FuncForPC(reflect.ValueOf(handle).Pointer()).Name()
		logger.Tracef("On entering handler %s, with params = %s", handlerName, ps)

		handle(w, r, ps)
	}
}
