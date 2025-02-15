package main

import (
	"log"
	"os"
	"wallet/cmd/app"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	if err := app.App.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
