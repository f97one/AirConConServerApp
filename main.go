package main

import (
	"flag"
	"fmt"
	"github.com/f97one/AirConCon/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
)

var conf *utils.AppConfig
var logger *logrus.Logger

func main() {
	flag.Parse()
	conf = utils.Load(flag.Arg(0))
	logger = utils.GetLogger()

	logger.Traceln("config :", conf)

	mux := httprouter.New()

	configureRouting(mux)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.ListenAddr, conf.ListenPort),
		Handler: mux,
	}

	logger.Traceln("Starting server.")
	err := server.ListenAndServe()
	if err != nil {
		logger.Fatalln(err)
	}
}
