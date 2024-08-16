package config

import (
	"encoding/json"
	"log"
	"os"
)

type CFG struct {
	PG_host      string `json:"pg_host"`
	PG_port      string `json:"pg_port"`
	PG_user      string `json:"pg_user"`
	PG_password  string `json:"pg_password"`
	PG_bdname    string `json:"pg_bdname"`
	Con_Attempts int    `json:"max_connection_attempts"`
}

func Read_Config() CFG {
	var config CFG
	data := make([]byte, 1024)
	file, err := os.Open("Server/cmd/config.cfg")
	if err != nil {
		// DEBUG confcle
		file, err = os.Open("config.cfg")
		if err != nil {
			log.Fatalln("Config.cfg reading: config.cfg is not found.")
		}
	}
	len, err := file.Read(data)
	if err != nil {
		log.Fatalln("Config.cfg reading: can't read config.cfg")
	}
	defer file.Close()
	data = append([]byte{byte(123)}, data[0:len]...)
	data = append(data[0:len], byte(125))
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalln("Config.cfg reading: config.txt unmarshalled incorrectly.")
	}
	log.Println("Config.cfg reading: config.cfg was readed")
	return config
}
