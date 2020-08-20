package main

import (
	"log"

	"github.com/WiggiLi/gowebsite/api"
	"github.com/WiggiLi/gowebsite/app"
	"github.com/WiggiLi/gowebsite/dal"
)

func run(errc chan<- error) {
	db, err := dal.NewPSQL("localhost", 1434)
	if err != nil {
		errc <- err
		return
	}

	application := app.NewApplication(db, errc)
	server := api.NewWebServer(application)

	server.Start(errc)
}

func main() {
	log.Print("Server is preparing to start...")

	errc := make(chan error)
	go run(errc)
	if err := <-errc; err != nil {
		log.Fatal(err)
	}
}
