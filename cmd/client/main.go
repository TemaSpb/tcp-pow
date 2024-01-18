package main

import (
	"log"

	"tcp-pow/internal/client"
	"tcp-pow/internal/config"
)

func main() {
	conf, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal(err)
	}

	cl := client.NewClient(conf)

	log.Println("start client")

	if err = cl.Run(); err != nil {
		log.Fatal(err)
	}
}
