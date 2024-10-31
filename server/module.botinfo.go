package server

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	"github.com/bwmarrin/discordgo"
	"github.com/guatom999/BBBot/modules/botHandlers"
	"github.com/guatom999/BBBot/modules/botRepositories"
	"github.com/guatom999/BBBot/modules/botUseCases"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "test",
			Description: "Test",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "feature",
					Description: "Test Bot",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		{
			Name:        "donate",
			Description: "promptpay qr",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "price",
					Description: "Bath Price",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
	}
)

type (
	IBotinfoModule interface {
		Init()
	}

	botInfoModule struct {
		module     *module
		botUseCase botUseCases.IBotUseCase
		botHandler botHandlers.IBotHandler
	}
)

func (m *module) BotinfoModule() IBotinfoModule {

	ctx := context.Background()

	gcpCli, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	botRepository := botRepositories.NewBotRepository()
	botfoUseCase := botUseCases.NewBotUseCase(botRepository, m.cfg, gcpCli)
	botfoHandler := botHandlers.NewBotHandler(botfoUseCase)

	return &botInfoModule{
		module:     m,
		botUseCase: botfoUseCase,
		botHandler: botfoHandler,
	}
}

func (b *botInfoModule) Init() {
	b.module.discordServer.commands = commands

	b.module.commandHandler["test"] = b.botHandler.Help
	b.module.commandHandler["donate"] = b.botHandler.Donate
}
