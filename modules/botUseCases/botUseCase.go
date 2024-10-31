package botUseCases

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/guatom999/BBBot/config"
	"github.com/guatom999/BBBot/modules/botRepositories"
	"github.com/guatom999/BBBot/utils"
	"github.com/skip2/go-qrcode"
)

type (
	IBotUseCase interface {
		Test() string
		RandomTarot()
		Donate(pctx context.Context, price string) string
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

	fmt.Println("convertPrice :::::::::>", convertPrice)

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

	// folderPath := "./temp"
	// filePath := folderPath + "/donate.png"

	// fmt.Println("filePath is ", filePath)

	// err = os.MkdirAll(folderPath, os.ModePerm)
	// if err != nil {
	// 	fmt.Println("Error creating folder:", err)
	// }

	// // Create the file in the specific folder
	// file, err := os.Create(filePath)
	// if err != nil {
	// 	fmt.Println("Error creating file:", err)

	// }
	// defer file.Close()

	// _, err = file.Write(png)
	// if err != nil {
	// 	log.Printf("Error: Failed to Write File: %v", err)
	// }

	// imageUrl := "localhost:27017/temp/dobate.png"

	return fileUrl
}
