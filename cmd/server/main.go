package main

import (
	"log"

	"tcp-pow/internal/config"
	"tcp-pow/internal/server"
)

func main() {
	conf, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("starting server")

	serv := server.NewServer(conf)

	if err = serv.Run(); err != nil {
		log.Fatal(err)
	}
}
