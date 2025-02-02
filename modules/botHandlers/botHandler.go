package botHandlers

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/guatom999/BBBot/modules/botUseCases"
)

type (
	IBotHandler interface {
		Help(s *discordgo.Session, i *discordgo.InteractionCreate)
		Donate(s *discordgo.Session, i *discordgo.InteractionCreate)
		GetFollowers(s *discordgo.Session, i *discordgo.InteractionCreate)
	}

	botHandler struct {
		botUseCase botUseCases.IBotUseCase
	}
)

func NewBotHandler(botUseCase botUseCases.IBotUseCase) IBotHandler {
	return &botHandler{botUseCase: botUseCase}
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
