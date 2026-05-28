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
	// ลบไฟล์เก่าก่อน (ถ้ามี) เพื่อโหลดเพลงใหม่
	if _, err := os.Stat("./audioTest.mp3"); err == nil {
		os.Remove("./audioTest.mp3")
	}

	// โหลดเพลงใหม่จาก URL ที่ได้รับ
	cmd := exec.Command("yt-dlp", "--ffmpeg-location", "C:\\ffmpeg\\bin\\", "-f", "bestaudio", "--extract-audio", "--audio-format", "mp3", "-o", fileName, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Download failed:", err)
		return err
	}

	fmt.Println("Downloaded:", fileName)
	return nil
}
