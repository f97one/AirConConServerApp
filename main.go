package main

import (
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"github.com/f97one/AirConCon/dataaccess"
	"github.com/f97one/AirConCon/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"time"
)

var conf *utils.AppConfig
var logger *logrus.Logger

const (
	currentUser string = "current_user_token"
	contentType string = "Content-Type"
	appJson     string = "application/json; charset=UTF-8"
)

func init() {
	gob.Register(dataaccess.AppUser{})
	gob.Register(dataaccess.JwtToken{})
}

func main() {
	flag.Parse()
	conf = utils.Load(flag.Arg(0))
	logger = utils.GetLogger()

	logger.Traceln("config :", conf)
	dataaccess.InitDatabase()

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
type scheduleResp struct {
	ScheduleId string `json:"schedule_id"`
	Name       string `json:"name"`
	OnOff      string `json:"on_off"`
	Weekday    []int  `json:"weekday"`
	Time       string `json:"time"`
	ScriptId   string `json:"script_id"`
}

// 次回のコントロールスケジュールを返却するレスポンスのJSON
type nextScheduleResp struct {
	ScheduleId string `json:"schedule_id"`
	Name       string `json:"name"`
	OnOff      string `json:"on_off"`
	WeekdayId  int    `json:"weekday_id"`
	Time       string `json:"time"`
	ScriptId   string `json:"script_id"`
}

func (s *scheduleResp) validate() error {
	alphaNumericUnder := regexp.MustCompile("^[0-9A-Za-z_]+$")

	// name
	// 33文字以上
	if len(s.Name) > 32 {
		return errors.New("name must be less than or equals to 32 characters")
	}
	if !alphaNumericUnder.MatchString(s.Name) {
		return errors.New("name must contain alphabet or number or underscore only")
	}

	// onOff
	// on, off 以外
	if !(s.OnOff == "on" || s.OnOff == "off") {
		return errors.New("on_off must be set 'on' or 'off'")
	}

	// weekday
	// スライスサイズ
	if len(s.Weekday) == 0 {
		return errors.New("weekday needs to set at least 1 day")
	}
	if len(s.Weekday) > 7 {
		return errors.New("weekday exceeds 7 days")
	}
	// 重複
	work := make(map[int]bool)
	var uniq []int
	for _, i := range s.Weekday {
		if !work[i] {
			work[i] = true
			uniq = append(uniq, i)
		}
	}
	if len(s.Weekday) != len(uniq) {
		return errors.New("weekday duplicates")
	}
	// 0～6
	for _, i := range s.Weekday {
		if i < 0 || i > 6 {
			return errors.New("weekday must be between 0 and 6")
		}
	}

	// time
	if s.Time == "" {
		return errors.New("time must not be empty")
	}
	_, err := time.Parse("15:04", s.Time)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid time value %s, must be between 00:00 - 23:59", s.Time))
	}

	// scriptId
	lowerHexChars := regexp.MustCompile("^[0-9a-f]{40}$")
	if !lowerHexChars.MatchString(s.ScriptId) {
		return errors.New("script_id must be 40 digits of lower cased hexadecimal")
	}

	return nil
}
