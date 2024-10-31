package server

import (
	"github.com/bwmarrin/discordgo"
	"github.com/guatom999/BBBot/config"
)

type (
	IModules interface {
		GetCommandHandler() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
		BotinfoModule() IBotinfoModule
	}

	module struct {
		discordServer  *discordServer
		cfg            *config.Config
		commandHandler map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	}
)

func InitModule(s *discordServer, cfg *config.Config) IModules {
	return &module{
		discordServer:  s,
		cfg:            cfg,
		commandHandler: make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate), 0),
	}
}

func (m *module) GetCommandHandler() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return m.commandHandler
}
