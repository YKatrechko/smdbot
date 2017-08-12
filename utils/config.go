package utils

import (
	"encoding/json"
	"log"
	"os"
)

//Configuration needed for plugins and bot
type Config struct {
	LogFile string
	DBFile  string
	Token   string
	//RedditUser     string
	//RedditPassword string
	//BrainFile      string
	//Administrators map[string]bool
}

//func (c *Config) IsAdmin(user string) bool {
//	_, ok := c.Administrators[user]
//	return ok
//}

func LoadConfig(f string) *Config {

	file, err := os.Open(f)
	if err != nil {
		log.Fatalln("Can't read config file")
	}
	defer file.Close()

	var config *Config
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		log.Fatalln("Can't parse json file")
	}

	return config
}

func DefaultConfig(cnfg *Config) {
	if cnfg.LogFile == "" {
		cnfg.LogFile = "smdbot.log"
	}
	if cnfg.DBFile == "" {
		cnfg.DBFile = "smdbase.db"
	}
}

func SaveConfig(cnfg *Config) {
	log.Println("Saving config...")
	log.Printf("Log file: %s\n", cnfg.LogFile)
	log.Printf("DB file: %s\n", cnfg.DBFile)
	//log.Printf("Token: %s\n", cnfg.Token)

	//TODO
}
