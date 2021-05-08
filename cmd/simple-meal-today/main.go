package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"smt/internal/pkg/smtbot"
)

var Debug bool
var Timeout int
var Token string
var TokenFile string
var RecordFile string
var AdminUsersID int
var DatabaseName string

func init() {
	flag.BoolVar(&Debug, "debug", true, "Enable debug mod for tgbotapi")
	flag.IntVar(&Timeout, "timeout", 60, "Timeout of tgbotapi update")
	flag.StringVar(&Token, "token", "", "telegram bot api token")
	flag.StringVar(&TokenFile, "tokenfile", "", "telegram bot api token save in json file")
	flag.StringVar(&RecordFile, "record", "record.json", "record of simple-meal-today bot")
	flag.IntVar(&AdminUsersID, "admin", 0, "Telegram admin user id")
	flag.StringVar(&DatabaseName, "db", "db.sqlite", "database name for sqlite3")
}

type TokenJsonFileStruct struct {
	Token string
}

func main() {
	flag.Parse()
	var err error

	// read token from file
	if TokenFile != "" {
		tokenJson := new(TokenJsonFileStruct)
		tokenJsonFile, err := os.Open(TokenFile)
		if err != nil {
			log.Fatal("Can not open tokenJson file " + err.Error())
		}
		err = json.NewDecoder(tokenJsonFile).Decode(tokenJson)
		if err != nil {
			log.Fatal("Can not decode tokenJson file " + err.Error())
		}
		if tokenJson.Token == "" {
			log.Fatal("Token not set in json file")
		}
		Token = tokenJson.Token
		tokenJsonFile.Close()
	}

	// check token
	if Token == "" {
		log.Fatal("TgBotAPI Token not set")
	}

	// create new smt bot
	smtbot, err := smtbot.NewSmtBot(Token, AdminUsersID, Debug, Timeout, RecordFile, DatabaseName)
	if err != nil {
		log.Fatal("Fail at creating smtbot " + err.Error())
	}

	// bot create success
	log.Println("TgBotAPI created successfully")

	log.Fatal(smtbot.Run())
}
