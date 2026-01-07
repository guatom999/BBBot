package botHandlers

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/guatom999/BBBot/config"
	"github.com/guatom999/BBBot/modules/botUseCases"
)

type (
	IBotHandler interface {
		Help(s *discordgo.Session, i *discordgo.InteractionCreate)
		Donate(s *discordgo.Session, i *discordgo.InteractionCreate)
		GetFollowers(s *discordgo.Session, i *discordgo.InteractionCreate)
		DisconnectAllMembers(s *discordgo.Session, i *discordgo.InteractionCreate)
		// Play(s *discordgo.Session, i *discordgo.InteractionCreate)
		Leave(s *discordgo.Session, i *discordgo.InteractionCreate)
	}

	botHandler struct {
		cfg        *config.Config
		botUseCase botUseCases.IBotUseCase
	}
)

var ChannelID string

func NewBotHandler(botUseCase botUseCases.IBotUseCase, cfg *config.Config) IBotHandler {
	return &botHandler{botUseCase: botUseCase, cfg: cfg}
}

func (h *botHandler) Help(s *discordgo.Session, i *discordgo.InteractionCreate) {
	command := i.ApplicationCommandData()
	messageContent := command.Options[0].StringValue()

	fmt.Println("Test")

	_ = messageContent

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: h.botUseCase.Test(),
			Embeds: []*discordgo.MessageEmbed{
				{
					Image: &discordgo.MessageEmbedImage{
						URL: "https://www.housesamyan.com/assets/uploads/movie/poster_web_path/20240925174230_4E9A127B-B7E3-4D05-AFC9-A343C944FD57.jpeg",
					},
				},
			},
		},
	})

}

func (h *botHandler) Donate(s *discordgo.Session, i *discordgo.InteractionCreate) {

	ctx := context.Background()

	command := i.ApplicationCommandData()
	messageContent := command.Options[0].StringValue()

	_ = messageContent
	_ = ctx

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Image: &discordgo.MessageEmbedImage{
						// URL: "https://www.housesamyan.com/assets/uploads/movie/poster_web_path/20240925174230_4E9A127B-B7E3-4D05-AFC9-A343C944FD57.jpeg",
						URL: fmt.Sprint(h.botUseCase.Donate(ctx, messageContent)),
					},
				},
			},
		},
	})

}

// func (h *botHandler) Play(s *discordgo.Session, i *discordgo.InteractionCreate) {

// 	ctx := context.Background()

// 	// time.Sleep(5 * time.Second)

// 	unfollwedList := h.botUseCase.GetFollowers(ctx)

// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
// 		Data: &discordgo.InteractionResponseData{
// 			Content: unfollwedList,
// 		},
// 	})

// 	// utils.DownloadYouTubeAudio("https://www.youtube.com/watch?v=8mLCaXvVPas")
// 	utils.DownloadYouTubeAudio("https://www.youtube.com/watch?v=Vw1mNzIIBpw")

// 	guild, err := s.State.Guild(GuildID)
// 	if err != nil {
// 		fmt.Println("Error finding guild:", err)
// 		return
// 	}

// 	for _, vs := range guild.VoiceStates {
// 		if vs.ChannelID != "" {
// 			ChannelID = vs.ChannelID
// 			break
// 		}
// 	}

// 	vc, err := s.ChannelVoiceJoin(GuildID, ChannelID, false, false)
// 	if err != nil {
// 		fmt.Println("error: failed to join Voice Channel")
// 		return
// 	}
// 	defer vc.Disconnect()

// 	pkg.PlayAudio(vc)

// 	// s.ChannelMessageSend(ChannelID, "You must be in a voice channel!")

// }

func (h *botHandler) DisconnectAllMembers(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := h.botUseCase.DisconnectAllMembers(s, h.cfg.App.GuildID); err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error: " + err.Error(),
			},
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Green apple",
		},
	})
}

func (h *botHandler) Leave(s *discordgo.Session, i *discordgo.InteractionCreate) {
	vc, err := s.ChannelVoiceJoin(h.cfg.App.GuildID, "", false, false)
	if err == nil {
		vc.Disconnect()
	}
}

func (h *botHandler) GetFollowers(s *discordgo.Session, i *discordgo.InteractionCreate) {

	ctx := context.Background()

	// time.Sleep(5 * time.Second)

	unfollwedList := h.botUseCase.GetFollowers(ctx)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: unfollwedList,
		},
	})

}
