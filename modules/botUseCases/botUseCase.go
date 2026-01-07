package botUseCases

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/bwmarrin/discordgo"
	"github.com/guatom999/BBBot/config"
	"github.com/guatom999/BBBot/modules/botRepositories"
	"github.com/guatom999/BBBot/pkg/requests"
	"github.com/guatom999/BBBot/utils"
	"github.com/robfig/cron"
	"github.com/skip2/go-qrcode"
)

type (
	IBotUseCase interface {
		Test() string
		RandomTarot()
		Donate(pctx context.Context, price string) string
		DisconnectAllMembers(session *discordgo.Session, guildID string) error
		ScheduleGetFollowers(session *discordgo.Session)
		GetFollowers(pctx context.Context) string
		MonitoringTicketShopServer()
	}

	botUseCase struct {
		botRepo botRepositories.IBotRepository
		cfg     *config.Config
		session *discordgo.Session
		cl      *storage.Client
	}
)

func NewBotUseCase(
	botRepo botRepositories.IBotRepository,
	config *config.Config,
	session *discordgo.Session,
	cli *storage.Client,
) IBotUseCase {
	return &botUseCase{
		botRepo: botRepo,
		cfg:     config,
		session: session,
		cl:      cli,
	}
}

func (u *botUseCase) Test() string {
	return "test"
}

func (u *botUseCase) RandomTarot() {

}

func (u *botUseCase) Donate(pctx context.Context, price string) string {

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	convertPrice, err := strconv.Atoi(price)
	if err != nil {
		log.Printf("Error: Prices is Invalid: %v", err)
	}

	rawQrcode := utils.GenQRCode(u.cfg, convertPrice)

	var png []byte

	png, err = qrcode.Encode(rawQrcode, qrcode.Medium, 256)
	if err != nil {
		log.Printf("Error: Failed to Encode Qr Code: %v", err)
	}

	fileUrl, err := utils.UploadFile(u.cfg, u.cl, ctx, png)
	if err != nil {
		log.Printf("Error: Failed to Generate QR Code: %v", err)
		return ""
	}

	return fileUrl
}

func (u *botUseCase) DisconnectAllMembers(session *discordgo.Session, guildID string) error {
	guild, err := session.State.Guild(guildID)
	if err != nil {
		guild, err = session.Guild(guildID)
		if err != nil {
			return fmt.Errorf("failed to get guild: %w", err)
		}
	}
	disconnectedCount := 0

	for _, voiceState := range guild.VoiceStates {
		if voiceState.UserID == session.State.User.ID {
			continue
		}

		err := session.GuildMemberMove(guildID, voiceState.UserID, nil)
		if err != nil {
			log.Printf("Failed to disconnect user %s: %v", voiceState.UserID, err)
			continue
		}

		disconnectedCount++

		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("Disconnected %d members from voice channels", disconnectedCount)
	return nil
}

func (u *botUseCase) ScheduleGetFollowers(session *discordgo.Session) {

	location, _ := time.LoadLocation("Asia/Bangkok")

	var file *os.File

	if _, err := os.Stat("followers_last.csv"); os.IsNotExist(err) {
		file, err = os.Create("followers_last.csv")
		if err != nil {
			log.Fatalf("Error creating followers_last.csv: %v", err)
		}

		if err := u.botRepo.GetFollowers("l3adzboss", true); err != nil {
			log.Fatalf("Error getting followers: %v", err)
		}
	}

	if _, err := os.Stat("followers_now.csv"); os.IsNotExist(err) {
		file, err = os.Create("followers_now.csv")
		if err != nil {
			log.Fatalf("Error creating followers_now.csv: %v", err)
		}

		if err := u.botRepo.GetFollowers("l3adzboss", false); err != nil {
			log.Fatalf("Error getting followers: %v", err)
		}

	}

	defer file.Close()

	c := cron.NewWithLocation(location)
	if err := c.AddFunc("@every 2h", func() {
		go func() {
			if err := u.botRepo.GetFollowers("l3adzboss", false); err != nil {
				log.Printf("Error: Failed to Get Followers: %v", err)
			}

			time.Sleep(time.Second * 60)

			lastFollowers := u.botRepo.GetLastFollowers()
			nowFollowers := u.botRepo.GetNowFollowers()

			diff := u.difference(lastFollowers, nowFollowers)
			if diff != "" {
				session.ChannelMessageSend(u.cfg.App.ChannelID, fmt.Sprintf("Unfollowed are: %s", diff))
			}
		}()

	}); err != nil {
		log.Printf("Error: Failed to AddFunc: %v", err)
	}

	// if err := c.AddFunc("@every 2h", func() {
	// 	time.Sleep(time.Second * 20)
	// 	go func() {
	// 		lastFollowers := u.botRepo.GetLastFollowers()
	// 		nowFollowers := u.botRepo.GetNowFollowers()

	// 		diff := u.difference(lastFollowers, nowFollowers)
	// 		if diff != "" {
	// 			session.ChannelMessageSend(u.cfg.App.ChannelID, fmt.Sprintf("Unfollowed are: %s", diff))
	// 		}
	// 	}()
	// }); err != nil {
	// 	log.Printf("Error: Failed to AddFunc: %v", err)
	// }

	c.Start()
}

func (u *botUseCase) MonitoringTicketShopServer() {

	monitorList := []string{
		"http://45.144.167.78:8090/movie/health",
		"http://45.144.167.78:8103/payment/health",
		"http://45.144.167.78:8102/ticket/health",
		"http://45.144.167.78:8101/inventory/health",
		"http://45.144.167.78:8100/user/health",
	}

	serviceNameList := []string{
		"movie",
		"payment",
		"ticket",
		"inventory",
		"user",
	}

	for {
		time.Sleep(time.Second * 10)
		for i, list := range monitorList {
			_, err := requests.NewRequest().Get(list)
			if err != nil {
				log.Printf("Error: Failed to Get Health: %v", err)
				u.session.ChannelMessageSend(u.cfg.App.ChannelID, fmt.Sprintf("มีเรื่องแล้ว service %s แตก", serviceNameList[i]))
			}
		}

	}

}

func (u *botUseCase) GetFollowers(pctx context.Context) string {

	if _, err := os.Stat("followers_last.csv"); os.IsNotExist(err) {
		file, err := os.Create("followers_last.csv")
		if err != nil {
			log.Fatalf("Error creating followers_last.csv: %v", err)
		}
		file.Close()
	}

	if _, err := os.Stat("followers_now.csv"); os.IsNotExist(err) {
		file, err := os.Create("followers_now.csv")
		if err != nil {
			log.Fatalf("Error creating followers_now.csv: %v", err)
		}
		file.Close()
	}

	lastFollowers := u.botRepo.GetLastFollowers()
	nowFollowers := u.botRepo.GetNowFollowers()

	return u.difference(lastFollowers, nowFollowers)
}

func (u *botUseCase) difference(last, now []string) string {

	m := make(map[string]struct{}, 0)
	var diff string

	if len(last) > len(now) {
		for _, s := range now {
			m[s] = struct{}{}
		}

		for _, s := range last {
			if _, found := m[s]; !found {
				diff += s + "\n"
			}
		}
		return diff

	} else if len(last) < len(now) {
		u.botRepo.GetFollowers("l3adzboss", true)
	}

	return ""
}
