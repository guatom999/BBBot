package main

import (
	"context"
	"log"
	"os"

	"github.com/guatom999/BBBot/config"
	"github.com/guatom999/BBBot/server"
)

func main() {

	ctx := context.Background()

	_ = ctx

	cfg := config.GetConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is required")
		}
		log.Printf("choosen env is :%v", os.Args[1])
		return os.Args[1]
	}())

	server.NewDiscordServer(&cfg).Start()

}
