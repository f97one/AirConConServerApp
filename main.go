package main

import (
	"flag"
	"fmt"
	"github.com/f97one/AirConCon/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
)

var conf *utils.AppConfig
var logger *logrus.Logger

func main() {
	flag.Parse()
	conf = utils.Load(flag.Arg(0))
	logger = utils.GetLogger()

	logger.Traceln("config :", conf)

	mux := httprouter.New()

	configureRouting(mux)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.ListenAddr, conf.ListenPort),
		Handler: mux,
	}

	logger.Traceln("Starting server.")
	err := server.ListenAndServe()
	if err != nil {
		logger.Fatalln(err)
	}
}

// JWTトークンを返却するレスポンスのJSON
type tokenResp struct {
	Token string `json:"token"`
}

// ユーザー追加に成功したユーザーを返却するレスポンスのJSON
type usernameResp struct {
	Username string `json:"username"`
}

// なんらかのメッセージを返却するレスポンスのJSON
type msgResp struct {
	Msg string `json:"msg"`
}

// コントロールスケジュールを返却するレスポンスのJSON(1要素分)
//type scheduleResp struct {
//	ScheduleId int `json:"schedule_id"`
//	Name string `json:"name"`
//	OnOff string `json:"on_off"`
//	Weekday []int `json:"weekday"`
//	Time time.Time `json:"time"`
//	ScriptId int `json:"script_id"`
//}
