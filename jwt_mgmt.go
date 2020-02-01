package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
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

	signBytes, err := ioutil.ReadFile("./" + privateKeyFile)
	if err != nil {
		logger.Errorln(err)
		respondErrorWithLog(&w, err, http.StatusInternalServerError)
		return "", now
	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		logger.Errorln(err)
		respondErrorWithLog(&w, err, http.StatusInternalServerError)
		return "", now
	}
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	// UUIDを生成して有効期限ごとに変えさせる
	u, err := uuid.NewRandom()
	if err != nil {
		logger.Errorln(err)
		respondErrorWithLog(&w, err, http.StatusInternalServerError)
		return "", now
	}
	claims["id"] = u.String()
	claims["name"] = username
	claims["exp"] = period.Unix()
	claims["iat"] = now.Unix()

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		logger.Errorln(err)
		respondErrorWithLog(&w, err, http.StatusInternalServerError)
		return "", now
	}
	return tokenString, period
}

func requireJwtHandler(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		verifyBytes, err := ioutil.ReadFile("./" + pkcs8PublicKeyFile)
		if err != nil {
			logger.Errorln(err)
			respondErrorWithLog(&w, err, http.StatusInternalServerError)
			return
		}
		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			logger.Errorln(err)
			respondErrorWithLog(&w, err, http.StatusInternalServerError)
			return
		}
		token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			_, err := token.Method.(*jwt.SigningMethodRSA)
			if !err {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			} else {
				return verifyKey, nil
			}
		})
		if err == nil && token.Valid {
			handle(w, r, ps)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

}
