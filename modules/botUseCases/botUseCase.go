package botUseCases

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/bwmarrin/discordgo"
	"github.com/guatom999/BBBot/config"
	"github.com/guatom999/BBBot/modules/botRepositories"
	"github.com/guatom999/BBBot/utils"
	"github.com/robfig/cron"
	"github.com/skip2/go-qrcode"
)

type (
	IBotUseCase interface {
		Test() string
		RandomTarot()
		Donate(pctx context.Context, price string) string
		ScheduleGetFollowers(session *discordgo.Session)
		GetFollowers(pctx context.Context) string
	}

	botUseCase struct {
		botRepo botRepositories.IBotRepository
		cfg     *config.Config
		cl      *storage.Client
	}
)

func NewBotUseCase(botRepo botRepositories.IBotRepository, config *config.Config, cli *storage.Client) IBotUseCase {
	return &botUseCase{botRepo: botRepo, cfg: config, cl: cli}
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
		log.Printf("Error: Failed to Generate QR Code: %v", err)
	}

	fileUrl, err := utils.UploadFile(u.cfg, u.cl, ctx, png)
	if err != nil {
		log.Printf("Error: Failed to Generate QR Code: %v", err)
		return ""
	}

	return fileUrl
}

func (u *botUseCase) ScheduleGetFollowers(session *discordgo.Session) {

	location, _ := time.LoadLocation("Asia/Bangkok")

	c := cron.NewWithLocation(location)
	if err := c.AddFunc("@every 1h", func() {
		go func() {
			if err := u.botRepo.GetFollowers("l3adzboss", false); err != nil {
				log.Printf("Error: Failed to Get Followers: %v", err)
			}
		}()

	}); err != nil {
		log.Printf("Error: Failed to AddFunc: %v", err)
	}

	if err := c.AddFunc("@every 1h", func() {
		time.Sleep(time.Second * 20)
		go func() {
			time.Sleep(time.Second * 10)
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

	c.Start()
}

func (u *botUseCase) GetFollowers(pctx context.Context) string {

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
