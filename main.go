package main

//package lingua

import (
	"log"
	"net/http"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	IsDebug   = true
	Token     = ""
	ServerUrl = "https://www.google.com/"
	Address   = "0.0.0.0:443"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(Token)
	if err != nil {
		log.Panic(err)
	}

	me, _ := bot.GetMe()

	bot.Debug = IsDebug

	if IsDebug {

		DeleteWebHookIfSet(bot)

		log.Println("Info about the bot:")
		log.Printf("ID: %d", me.ID)
		log.Printf("FirstName: '%s'", me.FirstName)
		log.Printf("LastName: '%s'", me.LastName)
		log.Printf("UserName: '%s'", me.UserName)
		log.Printf("LanguageCode: '%s'", me.LanguageCode)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	var updates tgbotapi.UpdatesChannel

	if IsDebug {

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 10
		updates, err = bot.GetUpdatesChan(u)

	} else {

		_, err = bot.SetWebhook(tgbotapi.NewWebhook(ServerUrl + bot.Token))
		if err != nil {
			log.Fatal(err)
		}

		updates = bot.ListenForWebhook("/" + bot.Token)
		go http.ListenAndServe(Address, nil)
	}

	for update := range updates {

		if update.Message == nil {
			continue
		}

		log.Println("Found updates")
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		log.Printf("%+v\n", update)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID

		log.Println("Sending message")
		bot.Send(msg)
	}
}

func DeleteWebHookIfSet(bot *tgbotapi.BotAPI) {

	hook, _ := bot.GetWebhookInfo()

	if !hook.IsSet() {
		return
	}

	log.Printf("Found webhook: '%s'", hook.URL)
	log.Printf("With error: '%s'", hook.LastErrorMessage)
	log.Printf("Pending updates: %d", hook.PendingUpdateCount)
	log.Println("Deleting...")

	response, _ := bot.RemoveWebhook()
	if response.Ok {
		log.Println("Deleted sucessfully...")
	} else{
		log.Printf("Failed to delete with status: %s", response.Description)
	}
}