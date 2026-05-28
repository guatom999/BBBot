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
		Title:       fmt.Sprintf("CI Alert for Project %s", message.ProjectName),
		Description: message.ProjectName,
		Timestamp:   utils.GetLocalBkkTime().Format("2006-01-02T15:04:05Z07:00"),
		Color:       message.Color,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: func(status string) string {
					switch status {
					case "success":
						return "✅ Status"
					case "failed":
						return "❌ Status"
					default:
						return "❓ Status"
					}
				}(message.Status),
				Value:  message.Message,
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
