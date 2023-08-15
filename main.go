package main

import (
	"log"

	"github.com/connoraubry/chess-server/server"
	"github.com/joho/godotenv"
)

type Game struct {
	ID      int64
	Fen     string
	Done    bool
	PgnPath string
}

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s := server.New()
	s.Run()
}
