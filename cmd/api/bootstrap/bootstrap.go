package bootstrap

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jpastorm/dialogflowbot/domain/product"
	"github.com/jpastorm/dialogflowbot/domain/user"
	"github.com/jpastorm/dialogflowbot/infraestructure/dialogflow"
	userHandler "github.com/jpastorm/dialogflowbot/infraestructure/handler/user"
	productStorage "github.com/jpastorm/dialogflowbot/infraestructure/postgres/product"
	userStorage "github.com/jpastorm/dialogflowbot/infraestructure/postgres/user"
	"github.com/jpastorm/dialogflowbot/infraestructure/response"
	"github.com/jpastorm/dialogflowbot/infraestructure/telegram"
)

// Logger interface
type Logger interface {
	Fatalf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

func Run() error {
	config := newConfiguration("./configuration.json")
	db := newPSQLDatabase(config)
	logger := newLogrus(config.LogFolder, true)
	api := newEcho(config, response.HTTPErrorHandler)
	loadSignatures(config, logger)
	//USER
	userUsecase := user.New(userStorage.New(db))
	userHandler.NewRouter(api, userUsecase)
	//PRODUCT
	productUsecase := product.New(productStorage.New(db))
	//DIALOGFLOW
	df := dialogflow.New(logger, config.DialogFlow.ProjectID)
	//TELEGRAM
	tg := telegram.New(logger, df, productUsecase, config.Telegram.Token)
	tg.RunService()

	port := fmt.Sprintf(":%d", config.PortHTTP)
	return api.Start(port)
}

func loadSignatures(conf Configuration, logger Logger) {
	priv := conf.PrivateFileSign
	publ := conf.PublicFileSign

	fpriv, err := ioutil.ReadFile(priv)
	checkErr(err, fmt.Sprintf("no se pudo leer el archivo de firma privado %s", fpriv))

	fpubl, err := ioutil.ReadFile(publ)
	checkErr(err, fmt.Sprintf("no se pudo leer el archivo de firma público %s", fpubl))
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
