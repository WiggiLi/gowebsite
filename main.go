package main

import (
	"gowebsite/api"
	"gowebsite/app"
	"gowebsite/dal"
	"log"
)

func run(errc chan<- error) {
	//time.Sleep(time.Second * 10)

	// TODO: init DAL here for MS SQL
	db, err := dal.NewMsSQL("localhost", 1434)
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
