package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/f97one/AirConCon/dataaccess"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func respondError(w *http.ResponseWriter, err error, sc int) {
	http.Error(*w, err.Error(), sc)
}

// ログインする。
func login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()

	logger.Traceln("リクエストボディを読み取り中")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusBadRequest)
		return
	}

	logger.Traceln("ユーザーデータをアンマーシャリング中")
	var reqUser *dataaccess.AppUser
	err = json.Unmarshal(body, &reqUser)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusBadRequest)
		return
	}

	logger.Traceln("ユーザー名のバリデート中")
	if strings.TrimSpace(reqUser.Username) == "" {
		logger.Errorln(err)
		respondError(&w, errors.New("username must not be empty"), http.StatusBadRequest)
		return
	}

	logger.Traceln("ユーザーを検索中")
	au, err := dataaccess.LoadByUsername(reqUser.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			// ユーザーなしの場合は401
			respondUnauthorized(w)
			return
		} else {
			logger.Errorln(err)
			respondError(&w, err, http.StatusInternalServerError)
			return
		}
	}

	// パスワードを照合
	logger.Traceln("パスワードを照合中")
	err = bcrypt.CompareHashAndPassword([]byte(au.Password), []byte(reqUser.Password))
	if err == nil {
		logger.Traceln("JWTを生成中")
		respondJwtToken(w, au)
	} else {
		// パスワード不一致も401
		logger.Traceln("パスワード不一致")
		respondUnauthorized(w)
	}
}

// JWTトークンを返送する。
func respondJwtToken(w http.ResponseWriter, au dataaccess.AppUser) (string, time.Time) {
	w.WriteHeader(http.StatusOK)
	jwtToken, expirationDate := genJwtToken(w, au.Username)

	err := dataaccess.PutToken(au.UserId, jwtToken, expirationDate)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return "", time.Now()
	}

	b, err := json.Marshal(&tokenResp{Token: jwtToken})
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return "", time.Now()
	}
	_, err = w.Write(b)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
	}

	return jwtToken, expirationDate
}

// ユーザーまたはパスワードが違う場合のレスポンスを作る
func respondUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	b, err := json.Marshal(&msgResp{Msg: "Invalid username or password"})
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
	}
	return
}

// ユーザーを追加する。
func subscribe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

// ログアウトさせる
func logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	claims, done, err := extractClaims(w, r)
	if done {
		return
	}
	username := claims["name"].(string)

	logger.Tracef("ユーザー %s を検索中", username)
	appUser, err := dataaccess.LoadByUsername(username)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

	logger.Traceln("カレントユーザーのJWTを削除")
	err = dataaccess.RemoveToken(appUser.UserId)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
	}
}

func extractClaims(w http.ResponseWriter, r *http.Request) (jwt.MapClaims, bool, error) {
	// JWTトークンの検証とClaimsをとってくる
	token, err := extractJwt(w, r)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return nil, true, err
	}
	if !token.Valid {
		logger.Errorln(err)
		respondError(&w, err, http.StatusBadRequest)
		return nil, true, nil
	}

	verifyKey, err := extractPublicKey(w)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusBadRequest)
		return nil, true, err
	}

	logger.Traceln("JWTからUsernameを抽出中")
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token.Raw, claims, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	return claims, false, err
}

// ユーザーのJWTを更新する
func auth(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	claims, done, _ := extractClaims(w, r)
	if done {
		return
	}

	username := claims["name"].(string)

	logger.Tracef("ユーザー %s を検索中", username)
	appUser, err := dataaccess.LoadByUsername(username)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

	logger.Traceln("JWTを生成中")
	respondJwtToken(w, appUser)
}
