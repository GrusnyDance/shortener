package main

import (
	"github.com/joho/godotenv"
	"log"
	"shortener/internal/app"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	app.Start()
}
