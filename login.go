package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/f97one/AirConCon/dataaccess"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"strings"
)

func respondErrorWithLog(w *http.ResponseWriter, err error, sc int) {
	http.Error(*w, err.Error(), sc)
}

// ログインする。
func login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorln(err)
		respondErrorWithLog(&w, err, http.StatusBadRequest)
		return
	}

	var reqUser *dataaccess.AppUser
	err = json.Unmarshal(body, &reqUser)
	if err != nil {
		logger.Errorln(err)
		respondErrorWithLog(&w, err, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(reqUser.Username) == "" {
		logger.Errorln(err)
		respondErrorWithLog(&w, errors.New("username must not be empty"), http.StatusBadRequest)
		return
	}

	au, err := dataaccess.LoadByUsername(reqUser.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			// ユーザーなしの場合は404
			userNotFound(w)
			return
		} else {
			logger.Errorln(err)
			respondErrorWithLog(&w, err, http.StatusInternalServerError)
			return
		}
	}

	// パスワードを照合
	err = bcrypt.CompareHashAndPassword([]byte(au.Password), []byte(reqUser.Password))
	if err != nil {
		// パスワード不一致も404
		userNotFound(w)
		return
	} else {
		respondJwtToken(w, au)
	}
}

// JWTトークンを返送する。
func respondJwtToken(w http.ResponseWriter, au dataaccess.AppUser) {
	w.WriteHeader(http.StatusOK)
	jwtToken, expirationDate := genJwtToken(w, au.Username)

	err := dataaccess.PutToken(au.UserId, jwtToken, expirationDate)
	if err != nil {
		logger.Errorln(err)
		respondErrorWithLog(&w, err, http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(&tokenResp{Token: jwtToken})
	if err != nil {
		logger.Errorln(err)
		respondErrorWithLog(&w, err, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		logger.Errorln(err)
		respondErrorWithLog(&w, err, http.StatusInternalServerError)
	}
}

// ユーザーまたはパスワードが違う場合のレスポンスを作る
func userNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	b, err := json.Marshal(&msgResp{Msg: "Invalid username or password"})
	if err != nil {
		logger.Errorln(err)
		respondErrorWithLog(&w, err, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		logger.Errorln(err)
		respondErrorWithLog(&w, err, http.StatusInternalServerError)
	}
	return
}

// ユーザーを追加する。
func subscribe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
