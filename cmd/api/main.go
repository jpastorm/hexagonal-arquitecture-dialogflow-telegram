package main

import (
	"log"

	"github.com/jpastorm/dialogflowbot/cmd/api/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}