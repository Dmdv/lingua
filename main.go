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

		// Process commands

		if update.Message.IsCommand(){

			// format is just a word
			log.Printf("This is command: '%s'", update.Message.Command())
			log.Printf("Arguments: '%s'", update.Message.CommandArguments())

			var text string

			switch update.Message.Command() {
			case "about":
				text = "Карманный переводчик с поддержкой множества языков"
			case "help":
				text = "При помощи команд из списка выберите пару языков, " +
					"словарь по умолчанию, количество словарных " +
					"статей на одно сообщение"
			case "from":
				text = "Вы выбрали язык, с которого надо перевести"
				// TODO: Reply buttons
			case "to":
				text = "Вы выбрали язык, на который надо перевести"
				// TODO: Reply buttons
			case "dic":
				text = "Вы выбрали словарь по умолчанию"
				// TODO: Reply buttons
			case "count":
				text = "Вы выбрали количество словарных статей на одно сообщение"
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			bot.Send(msg)
			continue
		}

		// Process callbacks

		// Reply with translation

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		/*
		str1 := "Back"
		markup := tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
				{
					tgbotapi.InlineKeyboardButton{Text: "Вперёд", CallbackData: &str1},
					tgbotapi.InlineKeyboardButton{Text: "Назад", CallbackData: &str1},
				},
			},
		}
		*/

		markup := CreateInlineButtons()
		// msg.ReplyToMessageID = update.Message.MessageID
		msg.ReplyMarkup = markup

		log.Println("Sending message")
		bot.Send(msg)
	}
}
func CreateKeyboardButtons() tgbotapi.ReplyKeyboardMarkup {
	btn1 := tgbotapi.NewKeyboardButton("KeyboardButton1")
	btn2 := tgbotapi.NewKeyboardButton("KeyboardButton2")
	row := tgbotapi.NewKeyboardButtonRow(btn1, btn2)
	return tgbotapi.NewReplyKeyboard(row)
}

func CreateInlineButtons() tgbotapi.InlineKeyboardMarkup {
	kbd1 := tgbotapi.NewInlineKeyboardButtonData("   Назад  ", "Back")
	kbd2 := tgbotapi.NewInlineKeyboardButtonData("   Вперёд   ", "Next")
	markup := tgbotapi.NewInlineKeyboardRow(kbd1, kbd2)
	return tgbotapi.NewInlineKeyboardMarkup(markup)
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