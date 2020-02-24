package dataaccess

import "time"

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
