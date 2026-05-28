package botHandlers

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/guatom999/BBBot/config"
	"github.com/guatom999/BBBot/modules/botUseCases"
	"github.com/jonas747/dca"
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

	fmt.Println("user used donate command")

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
						URL: fmt.Sprintf("attachment://%s.png", messageContent),
					},
				},
			},
			Files: []*discordgo.File{
				{
					Name:        fmt.Sprintf("%s.png", messageContent),
					ContentType: "image/png",
					Reader:      h.botUseCase.Donate(ctx, messageContent),
				},
			},
		},
	})

}

// func (h *botHandler) Play(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 	command := i.ApplicationCommandData()

// 	// รับ URL จาก command option
// 	var youtubeURL string
// 	if len(command.Options) > 0 {
// 		youtubeURL = command.Options[0].StringValue()
// 	}

// 	if youtubeURL == "" {
// 		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 			Type: discordgo.InteractionResponseChannelMessageWithSource,
// 			Data: &discordgo.InteractionResponseData{
// 				Content: "❌ กรุณาระบุ YouTube URL",
// 			},
// 		})
// 		return
// 	}

// 	// หา Voice Channel ที่ user อยู่
// 	guild, err := s.State.Guild(h.cfg.App.GuildID)
// 	if err != nil {
// 		fmt.Println("Error finding guild:", err)
// 		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 			Type: discordgo.InteractionResponseChannelMessageWithSource,
// 			Data: &discordgo.InteractionResponseData{
// 				Content: "❌ ไม่พบ Guild",
// 			},
// 		})
// 		return
// 	}

// 	// หา Voice Channel ที่ user ที่เรียก command อยู่
// 	var userVoiceChannelID string
// 	for _, vs := range guild.VoiceStates {
// 		if vs.UserID == i.Member.User.ID {
// 			userVoiceChannelID = vs.ChannelID
// 			break
// 		}
// 	}

// 	if userVoiceChannelID == "" {
// 		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 			Type: discordgo.InteractionResponseChannelMessageWithSource,
// 			Data: &discordgo.InteractionResponseData{
// 				Content: "❌ กรุณาเข้า Voice Channel ก่อนใช้คำสั่งนี้",
// 			},
// 		})
// 		return
// 	}

// 	// ตอบกลับว่ากำลังเล่นเพลง
// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
// 		Data: &discordgo.InteractionResponseData{
// 			Content: fmt.Sprintf("🎵 กำลังโหลดและเล่นเพลงจาก: %s", youtubeURL),
// 		},
// 	})

// 	// Download audio จาก YouTube
// 	if err := utils.DownloadYouTubeAudio(youtubeURL); err != nil {
// 		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
// 			Content: fmt.Sprintf("❌ ไม่สามารถโหลดเพลงได้: %v", err),
// 		})
// 		return
// 	}

// 	// Join Voice Channel
// 	vc, err := s.ChannelVoiceJoin(h.cfg.App.GuildID, userVoiceChannelID, false, false)
// 	if err != nil {
// 		fmt.Println("error: failed to join Voice Channel", err)
// 		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
// 			Content: "❌ ไม่สามารถเข้าร่วม Voice Channel ได้",
// 		})
// 		return
// 	}
// 	defer vc.Disconnect()

// 	// เล่นเพลง
// 	PlayAudioDCA(vc)
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

// PlayAudioDCA เล่นไฟล์ MP3 โดยใช้ dca library
func PlayAudioDCA(vc *discordgo.VoiceConnection) error {
	// รอให้ voice connection พร้อม
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			fmt.Println("Timeout waiting for voice connection to be ready")
			return fmt.Errorf("voice connection not ready")
		case <-ticker.C:
			fmt.Printf("Waiting... Ready: %v, OpusSend: %v\n", vc.Ready, vc.OpusSend != nil)
			if vc.Ready && vc.OpusSend != nil {
				goto READY
			}
		}
	}
READY:
	fmt.Println("Voice connection is ready!")

	// ตั้งค่า speaking
	if err := vc.Speaking(true); err != nil {
		fmt.Println("Error setting speaking:", err)
		return err
	}
	defer vc.Speaking(false)

	// ใช้ dca encode
	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 96
	opts.Application = dca.AudioApplicationAudio

	encodeSession, err := dca.EncodeFile("./audioTest.mp3", opts)
	if err != nil {
		fmt.Println("Error encoding audio file:", err)
		return err
	}
	defer encodeSession.Cleanup()

	fmt.Println("Starting to stream audio...")

	// Manual loop to send opus frames
	frameCount := 0
	for {
		frame, err := encodeSession.OpusFrame()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("End of audio file after %d frames\n", frameCount)
				break
			}
			fmt.Println("Error reading opus frame:", err)
			break
		}

		frameCount++
		if frameCount%100 == 0 {
			fmt.Printf("Sent %d frames...\n", frameCount)
		}

		// ส่ง frame ไปยัง Discord
		select {
		case vc.OpusSend <- frame:
		case <-time.After(5 * time.Second):
			fmt.Println("Timeout sending opus frame")
			return fmt.Errorf("timeout sending audio")
		}
	}

	fmt.Printf("Finished playing audio - Total frames: %d\n", frameCount)
	return nil
}
