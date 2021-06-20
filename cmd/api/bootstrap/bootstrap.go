package bootstrap

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jpastorm/dialogflowbot/domain/user"
	userHandler "github.com/jpastorm/dialogflowbot/infraestructure/handler/user"
	userStorage "github.com/jpastorm/dialogflowbot/infraestructure/postgres/user"
	"github.com/jpastorm/dialogflowbot/infraestructure/response"
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
	userUsecase := user.New(userStorage.New(db))
	userHandler.NewRouter(api, userUsecase)
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
