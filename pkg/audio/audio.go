package pkg

// import (
// 	"fmt"
// 	"io"
// 	"os"
// 	"os/exec"

// 	"github.com/bwmarrin/discordgo"
// 	"golang.org/x/crypto/bcrypt"
// )

// func PlayAudio(vc *discordgo.VoiceConnection) {
// 	file, err := os.Open("./audioTest.mp3")
// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		return
// 	}
// 	defer file.Close()

// 	bcrypt.CompareHashAndPassword([]byte("123"), []byte("123"))

// 	// ใช้ ffmpeg แปลง MP3 เป็น PCM (16-bit, 48kHz, 2-channel)
// 	ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
// 	ffmpeg.Stdin = file

// 	ffmpegOut, err := ffmpeg.StdoutPipe()
// 	if err != nil {
// 		fmt.Println("Error creating StdoutPipe:", err)
// 		return
// 	}

// 	err = ffmpeg.Start()
// 	if err != nil {
// 		fmt.Println("Error starting ffmpeg:", err)
// 		return
// 	}

// 	// เรียก sendOpus เพื่อแปลง PCM เป็น Opus และส่งไป Voice Channel
// 	sendOpus(vc, ffmpegOut)
// }

// // sendOpus แปลง PCM เป็น Opus และส่งไปยัง Discord
// func sendOpus(vc *discordgo.VoiceConnection, pcm io.Reader) {
// 	vc.Speaking(true)
// 	defer vc.Speaking(false)

// 	// สร้าง Opus Encoder
// 	opusEncoder, err := gopus.NewEncoder()
// 	if err != nil {
// 		fmt.Println("Error creating Opus encoder:", err)
// 		return
// 	}

// 	buf := make([]byte, 3840) // 20ms audio frame
// 	for {
// 		n, err := pcm.Read(buf)
// 		if err != nil {
// 			fmt.Println("Finished playing audio:", err)
// 			break
// 		}

// 		// แปลง PCM เป็น Opus
// 		opusData, err := opusEncoder.Encode(buf[:n], 3840) // 960 samples = 20ms frame
// 		if err != nil {
// 			fmt.Println("Error encoding Opus:", err)
// 			break
// 		}

// 		vc.OpusSend <- opusData
// 	}
// }
