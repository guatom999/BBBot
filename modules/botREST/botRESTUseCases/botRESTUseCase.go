package botRESTUseCases

import (
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/guatom999/BBBot/config"
	"github.com/guatom999/BBBot/modules/botREST"
	"github.com/guatom999/BBBot/utils"
)

type (
	EchoServerInterface interface {
		Send(message *botREST.InCommingMessage) error
	}

	echoServer struct {
		session *discordgo.Session
		cfg     *config.Config
	}
)

func NewEchoUserCase(session *discordgo.Session, cfg *config.Config) EchoServerInterface {
	return &echoServer{
		session: session,
		cfg:     cfg,
	}
}

func (s *echoServer) Send(message *botREST.InCommingMessage) error {

	embeded := &discordgo.MessageEmbed{
		Title:       "ci alert",
		Description: message.ProjectName,
		Timestamp:   utils.GetLocalBkkTime().Format("2006-01-02T15:04:05Z07:00"),
		Color:       0xFF0040,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "ci failed",
				Value:  fmt.Sprintf("error occurred cause %s", message.Status),
				Inline: true,
			},
		},
	}

	if _, err := s.session.ChannelMessageSendEmbed(s.cfg.App.ChannelID, embeded); err != nil {
		log.Printf("Error: Failed to send message: %s", err.Error())
		return errors.New("error: failed to send message")
	}

	return nil

}
