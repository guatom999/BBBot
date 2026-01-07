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
			Name:        "leave",
			Description: "force bot to leave channel",
		},
		{
			Name:        "getfollowers",
			Description: "get instagram followers",
		},
		{
			Name:        "greenapple",
			Description: "greenapple",
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
	// once.Do(func() {
	// 	instaBot = goinsta.New(m.cfg.User.Username, m.cfg.User.Password)
	// 	if err := instaBot.Login(); err != nil {
	// 		// log.Fatalf("Failed to login instagram: %v", err)
	// 		log.Println("Failed to login instagram: ", err)
	// 	}
	// })

	gcpCli, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	instaBot = nil

	botRepository := botRepositories.NewBotRepository(instaBot)
	botfoUseCase := botUseCases.NewBotUseCase(botRepository, m.cfg, m.discordServer.dg, gcpCli)
	botfoHandler := botHandlers.NewBotHandler(botfoUseCase, m.cfg)

	// go botfoUseCase.MonitoringTicketShopServer()

	// botfoUseCase.ScheduleGetFollowers(session)

	return &botInfoModule{
		module:     m,
		botUseCase: botfoUseCase,
		botHandler: botfoHandler,
	}
}

func (b *botInfoModule) Init() {
	b.module.discordServer.commands = commands

	// b.module.commandHandler["play"] = b.botHandler.Play
	b.module.commandHandler["leave"] = b.botHandler.Leave
	b.module.commandHandler["test"] = b.botHandler.Help
	b.module.commandHandler["donate"] = b.botHandler.Donate
	b.module.commandHandler["getfollowers"] = b.botHandler.GetFollowers
	b.module.commandHandler["greenapple"] = b.botHandler.DisconnectAllMembers
}
