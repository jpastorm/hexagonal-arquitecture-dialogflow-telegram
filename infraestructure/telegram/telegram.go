package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jpastorm/dialogflowbot/domain/product"
	"github.com/jpastorm/dialogflowbot/infraestructure/dialogflow"
	"github.com/jpastorm/dialogflowbot/model"
)

type Logger interface {
	Warnf(format string, args ...interface{})
}
type Usecase interface {
	DetectIntentText(sessionID, text, languageCode string) (string, error)
}
type Telegram struct {
	Token          string
	logger         Logger
	DialogFlow     dialogflow.Usecase
	productUsecase product.UseCase
}

func New(logger Logger, dialogflow dialogflow.Usecase, productUseCase product.UseCase, token string) *Telegram {
	return &Telegram{logger: logger, DialogFlow: dialogflow, productUsecase: productUseCase, Token: token}
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
		// fmt.Println("MIJSON")

		// fmt.Println(update.Message.From.UserName)
		// fmt.Println(update.Message.Text)
		// fmt.Print(update.Message.Chat.UserName)
		dfResponse, action, err := t.DialogFlow.DetectIntentText(update.Message.From.UserName, update.Message.Text, "es")
		if err != nil {
			t.logger.Warnf("error %v", err)
		}
		fmt.Println(dfResponse)
		fmt.Println(action)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, dfResponse)
		//msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
		message, err := t.typeOfAction(action)
		if err != nil {
			t.logger.Warnf(err.Error())
		}
		if message != "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
			bot.Send(msg)
		}
	}

}

func (t Telegram) typeOfAction(action string) (string, error) {
	var prefix = "input."

	switch action {
	case prefix + "lista":
		list, err := t.getAllProducts()
		if err != nil {
			return "", err
		}

		return list, nil
	}
	return "", nil
}

func (t Telegram) getAllProducts() (string, error) {
	var list string

	fields := model.Fields{}
	sortsFields := model.SortFields{}
	pag := model.Pagination{}
	products, err := t.productUsecase.GetAllWhere(fields, sortsFields, pag)
	if err != nil {
		return list, err
	}

	for i, product := range products {
		list += fmt.Sprintf("%d. %s   --> %.2f \n", i+1, product.Name, product.Price)
	}
	return list, nil
}
