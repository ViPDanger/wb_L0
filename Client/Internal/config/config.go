package config

import (
	"encoding/json"
	"log"
	"os"
)

type CFG struct {
	Adress       string `json:"adress"`
	Port         string `json:"port"`
	Con_Attempts int    `json:"max_connection_attempts"`
	Nats_host    string `json:"nats_host"`
	Nats_port    string `json:"nats_port"`
}

func Read_Config() CFG {
	var config CFG
	data := make([]byte, 1024)
	file, err := os.Open("Client/cmd/config.cfg")
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
