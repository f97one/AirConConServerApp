package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/f97one/AirConCon/dataaccess"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"strings"
)

// ユーザーを追加する。
func subscribe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username, err := getUsernameFromClaims(w, r)
	if err != nil {
		return
	}

	logger.Tracef("ユーザー %s を検索中", username)
	appUser, err := dataaccess.LoadByUsername(username)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

	if !appUser.AdminFlag {
		// 管理者フラグなしはユーザー追加要求を拒否
		logger.Warnln("カレントユーザーに管理者権限なし")
		w.Header().Add(contentType, appJson)
		w.WriteHeader(http.StatusUnauthorized)
		msg := "Current user has no privilege to add other user."
		b, err := json.Marshal(&msgResp{Msg: msg})
		if err != nil {
			logger.Errorln(err)
			respondError(&w, err, http.StatusInternalServerError)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			logger.Errorln(err)
			respondError(&w, err, http.StatusInternalServerError)
			return
		}
		return
	}

	logger.Traceln("ユーザーの追加開始")
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
			// ユーザーなしの場合はユーザーを追加する
			logger.Tracef("ユーザー %s を追加中", reqUser.Username)
			encodedPw, err := bcrypt.GenerateFromPassword([]byte(reqUser.Password), bcrypt.DefaultCost)
			if err != nil {
				logger.Errorln(err)
				respondError(&w, err, http.StatusInternalServerError)
				return
			}
			u := dataaccess.AppUser{
				Username:     reqUser.Username,
				Password:     fmt.Sprintf("%s", encodedPw),
				NeedPwChange: false,
				AdminFlag:    false,
			}
			logger.Tracef("ユーザー %s を永続化中", reqUser.Username)
			err = dataaccess.CreateUser(u)
			if err != nil {
				logger.Errorln(err)
				respondError(&w, err, http.StatusInternalServerError)
				return
			}
			// CREATEDでレスポンスjsonを返す
			logger.Traceln("レスポンス作成中")
			b, err := json.Marshal(usernameResp{Username: reqUser.Username})
			w.Header().Add(contentType, appJson)
			w.WriteHeader(http.StatusCreated)
			_, err = w.Write(b)
			if err != nil {
				logger.Errorln(err)
				respondError(&w, err, http.StatusInternalServerError)
				return
			}
		} else {
			logger.Errorln(err)
			respondError(&w, err, http.StatusInternalServerError)
		}
	} else {
		// ユーザーがいたのでCONFLICTを返す
		logger.Warnf("ユーザー %s は既に追加されている", au.Username)
		msg := fmt.Sprintf("User %s already exists.", au.Username)
		logger.Traceln("レスポンス作成中")
		w.Header().Add(contentType, appJson)
		w.WriteHeader(http.StatusConflict)
		b, err := json.Marshal(&msgResp{Msg: msg})
		if err != nil {
			logger.Errorln(err)
			respondError(&w, err, http.StatusInternalServerError)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			logger.Errorln(err)
			respondError(&w, err, http.StatusInternalServerError)
			return
		}
	}
}

// パスワードを変更する。
func changePassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username, err := getUsernameFromClaims(w, r)
	if err != nil {
		return
	}

	logger.Tracef("ユーザー %s を検索中", username)
	appUser, err := dataaccess.LoadByUsername(username)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}

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

	// ユーザーなしの場合はユーザーを追加する
	logger.Tracef("ユーザー %s のパスワードを変更中", username)
	encodedPw, err := bcrypt.GenerateFromPassword([]byte(reqUser.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	err = dataaccess.UpdatePassword(fmt.Sprintf("%s", encodedPw), appUser.UserId)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return
	}
	respondJwtToken(w, appUser)
}

// ログアウトして自身を削除する。
func unsubscribe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username, err := getUsernameFromClaims(w, r)
	if err != nil {
		return
	}

	logger.Tracef("ユーザー %s を検索中", username)
	appUser, err := dataaccess.LoadByUsername(username)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusNotFound)
		return
	}

	logger.Tracef("ユーザー %s を削除中", username)
	err = dataaccess.DeleteById(appUser.UserId)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	// ボディなし
}
