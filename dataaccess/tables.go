package dataaccess

// app_user の構造体
type AppUser struct {
	// ユーザーID
	UserId int32 `db:"user_id"`
	// ユーザー名
	Username string `db:"username"`
	// パスワード
	Password string `db:"password"`
	// ダイジェスト値
	DigestId *string `db:"digest_id"`
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
