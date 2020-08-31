package config

import (
	jsoniter "github.com/json-iterator/go"
	"log"
	"os"
)

var Config189 Cloud189Config

func LoadCloud189Config() {
	config := os.Getenv("CONFIG")
	err := jsoniter.Unmarshal([]byte(config), &Config189)
	if err != nil {
		log.Fatal("err：", err)
	}
}

type Cloud189Config struct {
	User       string     `json:"user"`
	Password   string     `json:"password"`
	RootId     string     `json:"root_id"`
	PwdDirId   []PwdDirId `json:"pwd_dir_id"`
	HideFileId string     `json:"hide_file_id"`
}
type PwdDirId struct {
	Id  string `json:"id"`
	Pwd string `json:"pwd"`
}
