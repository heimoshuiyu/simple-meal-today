package main

import (
	"flag"
	"log"
	"smt/internal/pkg/smtbot"
)

var Debug bool
var Timeout int
var Token string
var RecordFile string
var AdminUsersID int

func init() {
	flag.BoolVar(&Debug, "debug", true, "Enable debug mod for tgbotapi")
	flag.IntVar(&Timeout, "timeout", 60, "Timeout of tgbotapi update")
	flag.StringVar(&Token, "token", "", "telegram bot api token")
	flag.StringVar(&RecordFile, "record", "record.json", "record of simple-meal-today bot")
	flag.IntVar(&AdminUsersID, "admin", 0, "Telegram admin user id")
}

func main() {
	flag.Parse()
	var err error

	// check token set
	if Token == "" {
		log.Fatal("Telegram bot api token not set")
	}

	// create new smt bot
	smtbot, err := smtbot.NewSmtBot(Token, Debug, Timeout, RecordFile)
	if err != nil {
		log.Fatal("Fail at creating smtbot " + err.Error())
	}

	// bot create success
	log.Println("TgBotAPI created successfully")

	log.Fatal(smtbot.Run())
}
