package main

//package lingua

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
)

const (
	Token = ""
	ServerUrl = "https://www.google.com/"
	Address = "0.0.0.0:443"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(Token)
	if err != nil {
		log.Panic(err)
	}

	me, _ := bot.GetMe()

	log.Println(me.FirstName)
	log.Println(me.LastName)
	log.Println(me.UserName)
	log.Println(me.LanguageCode)
	log.Println(me.ID)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(ServerUrl + bot.Token))

	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)

	go http.ListenAndServe(Address, nil)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		log.Printf("%+v\n", update)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}