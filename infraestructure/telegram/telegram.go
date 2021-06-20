package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jpastorm/dialogflowbot/infraestructure/dialogflow"
)

type Logger interface {
	Warnf(format string, args ...interface{})
}
type Usecase interface {
	DetectIntentText(sessionID, text, languageCode string) (string, error)
}
type Telegram struct {
	Token      string
	logger     Logger
	DialogFlow dialogflow.Usecase
}

func New(logger Logger, dialogflow dialogflow.Usecase, token string) *Telegram {
	return &Telegram{logger: logger, DialogFlow: dialogflow, Token: token}
}

func (t Telegram) RunService() {
	bot, err := tgbotapi.NewBotAPI(t.Token)
	if err != nil {
		t.logger.Warnf(err.Error())
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		t.logger.Warnf(err.Error())
	}
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		//MarshalIndent
		fmt.Println("MIJSON")

		fmt.Println(update.Message.From.UserName)
		fmt.Println(update.Message.Text)
		fmt.Print(update.Message.Chat.UserName)
		dfResponse, err := t.DialogFlow.DetectIntentText(update.Message.From.UserName, update.Message.Text, "es")
		if err != nil {
			t.logger.Warnf("error %v", err)
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, dfResponse)
		msg.ReplyToMessageID = update.Message.MessageID
		fmt.Println(dfResponse)
		bot.Send(msg)
	}
}
