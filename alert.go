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
		alertTelegram(config.Alert.Target, config.Alert.Auth, msg)
	}else{
		fmt.Println("[Alert] ", msg)
	}
}

func alertTelegram(target, auth, msg string) bool {
	if(len(target) > 3 && len(auth) > 20 && len(msg) > 5){
		bot, err := tgbotapi.NewBotAPI(auth)
		if err != nil {
			log.Panic(err)
		}
		chatInt, _ := strconv.Atoi(target)
		msg := tgbotapi.NewMessage(int64(chatInt), msg)
		bot.Send(msg)
		return true
	}else{
		return false
	}
}
