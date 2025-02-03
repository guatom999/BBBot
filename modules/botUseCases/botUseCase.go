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
	if err := c.AddFunc("@every 25m", func() {
		go func() {
			if err := u.botRepo.GetFollowers(u.cfg.User.Username, u.cfg.User.Password, "l3adzboss"); err != nil {
				log.Printf("Error: Failed to Get Followers: %v", err)
			}
		}()

	}); err != nil {
		log.Printf("Error: Failed to AddFunc: %v", err)
	}

	if err := c.AddFunc("@every 5m", func() {
		go func() {
			time.Sleep(time.Second * 10)
			lastFollowers := u.botRepo.GetLastFollowers()
			nowFollowers := u.botRepo.GetNowFollowers()

			diff := u.difference(lastFollowers, nowFollowers)
			if diff != "" {
				session.ChannelMessageSend("1171470545225793576", fmt.Sprintf("Unfollowed are: %s", diff))
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

	fmt.Printf("lastFollowers length: %d\n", len(lastFollowers))
	fmt.Printf("nowFollowers length: %d\n", len(nowFollowers))

	return u.difference(lastFollowers, nowFollowers)
}

func (u *botUseCase) difference(a, b []string) string {

	fmt.Println("a", len(a), "b", len(b))
	m := make(map[string]struct{}, 0)

	for _, s := range b {
		m[s] = struct{}{}
	}

	var diff string
	for _, s := range a {
		if _, found := m[s]; !found {
			diff += s + "\n"
		}
	}

	return diff
}
