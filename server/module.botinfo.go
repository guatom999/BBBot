package server

import (
	"context"
	"log"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/ahmdrz/goinsta/v2"
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
					Name: "feature",

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
		{
			Name:        "getfollowers",
			Description: "get instagram followers",
		},
	}
	instaBot *goinsta.Instagram
	once     sync.Once
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

func (m *module) BotinfoModule(session *discordgo.Session) IBotinfoModule {

	ctx := context.Background()
	once.Do(func() {
		instaBot = goinsta.New(m.cfg.User.Username, m.cfg.User.Password)
		if err := instaBot.Login(); err != nil {
			log.Fatalf("Failed to login instagram: %s", err.Error())
		}
	})

	gcpCli, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	botRepository := botRepositories.NewBotRepository(instaBot)
	botfoUseCase := botUseCases.NewBotUseCase(botRepository, m.cfg, gcpCli)
	botfoHandler := botHandlers.NewBotHandler(botfoUseCase)

	botfoUseCase.ScheduleGetFollowers(session)

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
	b.module.commandHandler["getfollowers"] = b.botHandler.GetFollowers
}
