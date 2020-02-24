package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type AppConfig struct {
	ProductionMode bool   `json:"production_mode"`
	ListenAddr     string `json:"listen_addr"`
	ListenPort     int    `json:"listen_port"`
	SignalDbFile   string `json:"signal_db_file"`
}

var conf *AppConfig

func Load(path string) *AppConfig {
	if conf == nil {
		if len(path) == 0 {
			conf = &AppConfig{
				ProductionMode: false,
				ListenAddr:     "0.0.0.0",
				ListenPort:     8080,
				SignalDbFile:   "signal_codes.db",
			}
		} else {
			file, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatalln(err)
			}

			err = json.Unmarshal(file, &conf)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	return conf
}
