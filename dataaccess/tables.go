package dataaccess

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	Freq36kHz float64 = 36
	Freq40kHz float64 = 40
	Freq56kHz float64 = 56
)

// app_user の構造体
type AppUser struct {
	// ユーザーID
	UserId int `db:"user_id"`
	// ユーザー名
	Username string `db:"username" json:"username"`
	// パスワード
	Password string `db:"password" json:"password"`
	// ダイジェスト値
	NeedPwChange bool `db:"need_pw_change"`
	// 管理者フラグ
	AdminFlag bool `db:"admin_flag"`
}

// jwt_token の構造体
type JwtToken struct {
	// ユーザーID
	UserId int `db:"user_id"`
	// 生成済みJWTトークン
	GeneratedToken string `db:"generated_token"`
	// トークンの有効期限
	ExpiresAt string `db:"expires_at"`
}

// schedule の構造体
type Schedule struct {
	// スケジュール番号
	ScheduleId string `db:"schedule_id"`
	// スケジュール名
	Name string `db:"name"`
	// オンオフ種別
	OnOff bool `db:"on_off"`
	// 実行時間
	ExecuteTime string `db:"execute_time"`
	// 実行スクリプト番号
	ScriptId string `db:"script_id"`
	// 実行日
	ExecDay []Timing
}

// timing の構造体
type Timing struct {
	// スケジュール番号
	ScheduleId string `db:"schedule_id"`
	// 実行日
	WeekdayId time.Weekday `db:"weekday_id"`
}

// scripts の構造体
type Scripts struct {
	// スクリプト番号
	ScriptId string `db:"script_id" json:"script_id"`
	// GPIOピン番号
	Gpio int `db:"gpio" json:"gpio"`
	// スクリプト名
	ScriptName string `db:"script_name" json:"name"`
	// サンプリング周波数
	Freq float64 `db:"freq" json:"freq"`
}

// 次回スケジュール返送用構造体
type NextSchedule struct {
	// スケジュール番号
	ScheduleId string `db:"schedule_id"`
	// スケジュール名
	Name string `db:"name"`
	// オンオフ種別
	OnOff bool `db:"on_off"`
	// 実行時間
	ExecuteTime string `db:"execute_time"`
	// 実行スクリプト番号
	ScriptId string `db:"script_id"`
	// 実行日
	WeekdayId int `db:"weekday_id"`
}

// job_schedule の構造体
type JobSchedule struct {
	// 登録したスケジュールID
	ScheduleId string `db:"schedule_id"`
	// ジョブID
	JobId int `db:"job_id"`
	// 実行コマンドライン
	CmdLine string `db:"cmd_line"`
	// 実行予定日時
	RunAt string `db:"run_at"`
}

func (au *AppUser) Validate(checkUsername bool, forceUserMinLength bool, forcePasswdMinLength bool) error {
	alphaNumeric := regexp.MustCompile("^[0-9a-z]+$")
	bCryptPrefix := regexp.MustCompile("^\\$2[aby]\\$")

	userMinLen := 0
	if forceUserMinLength {
		userMinLen = 5
	}
	passwdMinLen := 0
	if forcePasswdMinLength {
		passwdMinLen = 7
	}

	// username
	if checkUsername {
		if len(au.Username) <= userMinLen {
			txt := "username must not be empty"
			if userMinLen != 0 {
				txt = fmt.Sprintf("username must be greater than %d characters", userMinLen)
			}
			return errors.New(txt)
		}
		if len(au.Username) > 32 {
			return errors.New("username must be less than or equals to 32 characters")
		}
		if !alphaNumeric.MatchString(au.Username) {
			return errors.New("username must contain alphabet or numeric")
		}
	}

	// password
	if !bCryptPrefix.MatchString(au.Password) {
		if len(au.Password) <= passwdMinLen {
			txt := "password must not be empty"
			if passwdMinLen != 0 {
				txt = fmt.Sprintf("password must be greater than %d characters", passwdMinLen)
			}
			return errors.New(txt)
		}
		if len(au.Password) > 32 {
			return errors.New("password must be less than or equals to 32 characters")
		}
	}

	return nil
}

func (s *Scripts) Validate() error {
	// GPIO
	// 1～40
	if s.Gpio < 1 || s.Gpio > 40 {
		return errors.New("gpio must be between 1 to 40")
	}

	// ScriptName
	alphaNumericUnder := regexp.MustCompile("^[0-9a-zA-Z_]+$")
	if len(s.ScriptName) > 32 {
		return errors.New("ScriptName must be less than or equals to 32 characters")
	}
	if len(strings.TrimSpace(s.ScriptName)) == 0 {
		return errors.New("ScriptName must not be empty")
	}
	if !alphaNumericUnder.MatchString(s.ScriptName) {
		return errors.New("ScriptName must contain alphabet or number or underscore only")
	}

	// Freq
	if !(s.Freq == Freq36kHz || s.Freq == Freq40kHz || s.Freq == Freq56kHz) {
		return errors.New("Freq must be either 36, 40, or 56")
	}
	return nil
}
