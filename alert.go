package main

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strconv"
)

func alert(name string, status bool) {
	var msg string
	if(!status){
		msg = "Your Service `" + name + "` is offline!"
	}else{
		msg = "Your Service `" + name + "` is back online!"
	}

	if(config.Alert.Typ == "alertTelegram"){
		alertTelegram(config.Alert.Target, msg)
	}else{
		fmt.Println("[Alert] ", msg)
	}
}

func startAlert(){
	if(config.Alert.Typ == "alertTelegram"){
		startTelegram()
	}
}

func startTelegram(){
	var err error
	alertAPI.Telegram, err = tgbotapi.NewBotAPI(config.Alert.Auth)
	if err != nil {
		log.Panic(err)
	}
}

func alertTelegram(target, msg string) bool {
	if(len(target) > 3 && alertAPI.Telegram != nil && len(msg) > 5){
		chatInt, _ := strconv.Atoi(target)
		msg := tgbotapi.NewMessage(int64(chatInt), msg)
		alertAPI.Telegram.Send(msg)

		return true
	}else{
		return false
	}
}
