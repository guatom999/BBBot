package main

import (
	"context"

	"github.com/guatom999/BBBot/config"
	"github.com/guatom999/BBBot/server"
)

func main() {

	ctx := context.Background()

	_ = ctx

	cfg := config.GetConfig()

	server.NewDiscordServer(&cfg).Start()

}
