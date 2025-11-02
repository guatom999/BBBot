package utils

import (
	"fmt"
	"os"
	"os/exec"
)

var (
	fileName = "audioTest.mp3"
)

func DownloadYouTubeAudio(url string) error {

	_, err := os.Open("./audioTest.mp3")
	if err != nil {
		cmd := exec.Command("yt-dlp", "--ffmpeg-location", "C:\\ffmpeg\\bin\\", "-f", "bestaudio", "--extract-audio", "--audio-format", "mp3", "-o", fileName, url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Println("Download failed:", err)
			return err
		}
	}

	fmt.Println("Downloaded:", fileName)
	return nil
}
