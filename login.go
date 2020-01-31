package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// ログインする。
func login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("login called")
}

// ユーザーを追加する。
func subscribe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
