package dataaccess

// app_user の構造体
type AppUser struct {
	// ユーザーID
	UserId int32 `db:"user_id"`
	// ユーザー名
	Username string `db:"username" json:"username"`
	// パスワード
	Password string `db:"password" json:"password"`
	// ダイジェスト値
	NeedPwChange bool `db:"need_pw_change"`
}

// jwt_token の構造体
type JwtToken struct {
	// ユーザーID
	UserId int32 `db:"user_id"`
	// 生成済みJWTトークン
	GeneratedToken string `db:"generated_token"`
	// トークンの有効期限
	ExpiresAt string `db:"expires_at"`
}
