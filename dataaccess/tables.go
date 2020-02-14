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
	ScriptId   string  `db:"script_id" json:"script_id"`
	Gpio       int     `db:"gpio" json:"gpio"`
	ScriptName string  `db:"script_name" json:"name"`
	Freq       float64 `db:"freq" json:"freq"`
}
