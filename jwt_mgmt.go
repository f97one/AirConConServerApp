package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/f97one/AirConCon/dataaccess"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	privateKeyFile     = "airconcon_jwt_rsa"
	publicKeyFile      = privateKeyFile + ".pub"
	pkcs8PublicKeyFile = publicKeyFile + ".pkcs8"
	validPeriod        = time.Hour * 72 // 現在時刻から72時間(=3日)
)

func genJwtToken(w http.ResponseWriter, username string) (string, time.Time) {
	now := time.Now()
	period := now.Add(validPeriod)

	logger.Traceln("秘密鍵を読み取り中")
	signBytes, err := ioutil.ReadFile("./" + privateKeyFile)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return "", now
	}
	logger.Traceln("秘密鍵の抽出中")
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return "", now
	}

	logger.Traceln("JWT Claims 生成中")
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	// UUIDを生成して有効期限ごとに変えさせる
	u, err := uuid.NewRandom()
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return "", now
	}
	claims["id"] = u.String()
	claims["name"] = username
	claims["exp"] = period.Unix()
	claims["iat"] = now.Unix()

	logger.Traceln("JWTを署名中")
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return "", now
	}
	return tokenString, period
}

func requireJwtHandler(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token, err := extractJwt(w, r)
		if err != nil {
			return
		}
		if token.Valid {
			logger.Traceln("トークンに問題なし")

			// 有効期限切れかどうかを確認する
			verifyKey, err := extractPublicKey(w)
			if err != nil {
				logger.Errorln(err)
				respondError(&w, err, http.StatusBadRequest)
				return
			}

			claims := jwt.MapClaims{}
			_, err = jwt.ParseWithClaims(token.Raw, claims, func(token *jwt.Token) (interface{}, error) {
				return verifyKey, nil
			})
			// claimsに格納したUnixタイムスタンプがなぜかfloat64扱いになっているのでキャストで対応
			exp := claims["exp"].(float64)
			expiresAt := time.Unix(int64(exp), 0)
			if expiresAt.Before(time.Now()) {
				msg := "Given token has expired."
				logger.Errorln(msg)
				w.WriteHeader(http.StatusUnauthorized)
				b, err := json.Marshal(&msgResp{Msg: msg})
				if err != nil {
					logger.Errorln(err)
					respondError(&w, err, http.StatusBadRequest)
					return
				}
				_, err = w.Write(b)
				if err != nil {
					logger.Errorln(err)
					respondError(&w, err, http.StatusBadRequest)
					return
				}
				return
			}

			// claimのユーザーが存在するか否かを確認する
			username := claims["name"].(string)
			_, err = dataaccess.LoadByUsername(username)
			if err != nil {
				logger.Errorln(err)
				msg := msgResp{Msg: fmt.Sprintf("user %s not found", username)}
				b, err := json.Marshal(msg)
				if err != nil {
					logger.Errorln(err)
					respondError(&w, err, http.StatusInternalServerError)
					return
				}
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, err = w.Write(b)
				if err != nil {
					logger.Errorln(err)
					respondError(&w, err, http.StatusInternalServerError)
					return
				}
			}
			logger.Tracef("user %s accessed within the expiration date.", username)

			handle(w, r, ps)
		} else {
			logger.Warnln("不正トークンを検知")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

}

func extractJwt(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	verifyKey, err := extractPublicKey(w)
	if err != nil {
		return nil, err
	}

	logger.Traceln("JWTの署名を検証中")
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		_, err := token.Method.(*jwt.SigningMethodRSA)
		if !err {
			logger.Warnf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return verifyKey, nil
		}
	})
	return token, err
}

func extractPublicKey(w http.ResponseWriter) (*rsa.PublicKey, error) {
	logger.Traceln("公開鍵を読み込み中")
	verifyBytes, err := ioutil.ReadFile("./" + pkcs8PublicKeyFile)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return nil, err
	}

	logger.Traceln("公開鍵を抽出中")
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		logger.Errorln(err)
		respondError(&w, err, http.StatusInternalServerError)
		return nil, err
	}
	return verifyKey, nil
}
