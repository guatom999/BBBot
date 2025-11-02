package server

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/guatom999/BBBot/config"
)

var (
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

type (
	IDiscordServer interface {
		Start()
	}

	discordServer struct {
		cfg      *config.Config
		dg       *discordgo.Session
		commands []*discordgo.ApplicationCommand
	}
)

func NewDiscordServer(cfg *config.Config) IDiscordServer {

	dg, err := discordgo.New("Bot " + cfg.App.Token)
	if err != nil {
		log.Fatalf("Invalid bot token :%v", err.Error())
	}

	return &discordServer{
		cfg:      cfg,
		dg:       dg,
		commands: make([]*discordgo.ApplicationCommand, 0),
	}
}

func (s *discordServer) Start() {

	s.dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v # %v", s.State.User.Username, s.State.User.Discriminator)
	})

	if err := s.dg.Open(); err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	module := InitModule(s, s.cfg)
	module.BotinfoModule(s.dg).Init()

	registeredCommands := make([]*discordgo.ApplicationCommand, len(s.commands))
	for i, v := range s.commands {
		cmd, err := s.dg.ApplicationCommandCreate(s.dg.State.User.ID, s.cfg.App.GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.dg.Close()

	s.dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		h, ok := module.GetCommandHandler()[i.ApplicationCommandData().Name]
		if ok {
			h(s, i)
		}
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl + C to exit")

	restServer := NewEchoServer(s.dg, s.cfg)
	restServer.Start(context.Background())

	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")

		for _, v := range registeredCommands {
			err := s.dg.ApplicationCommandDelete(s.dg.State.User.ID, s.cfg.App.GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
		restServer.GracefulShutdown(context.Background(), stop)
	}

	log.Println("Gracefully shutting down.")

}
