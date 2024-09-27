package utils

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Converter struct {
	bot *tgbotapi.BotAPI
}

func NewWebmConverter(bot *tgbotapi.BotAPI) *Converter {
	return &Converter{
		bot: bot,
	}
}

func (s *Converter) HandleMessage(msg *tgbotapi.Message) {
	if msg == nil {
		log.Println("Received nil message")
		return
	}

	var fileURL string
	if msg.Document != nil && strings.HasSuffix(msg.Document.FileName, ".webm") {
		fileID := msg.Document.FileID
		file, err := s.bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
		if err != nil {
			log.Println("Error getting file:", err)
			return
		}
		fileURL = file.Link(s.bot.Token)
	}
	if msg.Text != "" && strings.HasSuffix(msg.Text, ".webm") {
		fileURL = msg.Text
	}
	if fileURL != "" {
		go s.processFile(msg, fileURL)
	}
}
func (s *Converter) processFile(msg *tgbotapi.Message, fileURL string) {
	EnsureTempDir()
	fileName := filepath.Base(fileURL)
	localFileName := "./temp/" + fileName
	// Download the file
	out, err := os.Create(localFileName)
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer out.Close()
	err = s.DownloadFile(fileURL, out)
	if err != nil {
		log.Println("Error downloading file:", err)
		DeleteFile(localFileName)
		return
	}

	// Convert the file
	outputFile, err := ConvertWebmToMp4(localFileName)
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && (exitError.ExitCode() == 187 || exitError.ExitCode() == 8) {
			log.Printf("Error converting file with exit status %d, trying with pixel adjustment or alternative settings", exitError.ExitCode())
			outputFile, err = ConvertWebmToMp4WithPixel(localFileName)
			if err != nil {
				log.Println("Error converting file with pixel adjustment or alternative settings:", err)
				SendMessage(s.bot, msg.Chat.ID, fmt.Sprintf("Error converting file: %s", err))
				DeleteFile(localFileName)
				DeleteFile(outputFile)
				return
			}
		} else {
			log.Println("Error converting file:", err)
			SendMessage(s.bot, msg.Chat.ID, fmt.Sprintf("Error converting file: %s", err))
			DeleteFile(localFileName)
			DeleteFile(outputFile)
			return
		}
	}

	// Check if the output file is non-empty
	if fileInfo, err := os.Stat(outputFile); err != nil || fileInfo.Size() == 0 {
		log.Println("Converted file is empty or does not exist.")
		SendMessage(s.bot, msg.Chat.ID, "The converted file is empty or could not be created.")
		DeleteFile(localFileName)
		DeleteFile(outputFile)
		return
	}

	// Send the converted video
	err = SendVideo(s.bot, msg.Chat.ID, outputFile)
	if err != nil {
		log.Printf("Error while sending video: %v", err)
		DeleteFile(localFileName)
		DeleteFile(outputFile)
	}

	log.Printf("%s conversion completed.", fileName)
	dstDir := "./media/downloads"
	err = MoveFile(outputFile, dstDir)
	if err != nil {
		log.Println("Error moving file:", err)
		DeleteFile(localFileName)
		DeleteFile(outputFile)
		return
	}
	log.Printf("File %s has been moved to %s.", outputFile, dstDir)
	DeleteFile(localFileName)
}

func DeleteFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	err = file.Close()
	if err != nil {
		return err
	}
	// Usuń plik
	err = os.Remove(filePath)
	if err != nil {
		return err
	}
	log.Printf("Temporary files have been cleaned up.")

	return nil
}

func (s *Converter) DownloadFile(url string, dest *os.File) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(dest, resp.Body)

	return err
}

func ConvertWebmToMp4(inputFile string) (string, error) {
	outputFile := strings.TrimSuffix(inputFile, ".webm") + ".mp4"
	cmd := exec.Command("ffmpeg", "-loglevel", "error", "-i", inputFile, outputFile)

	var stderr strings.Builder
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("ffmpeg error: %s", stderr.String())
	}
	return outputFile, err
}

func ConvertWebmToMp4WithPixel(inputFile string) (string, error) {
	outputFile := strings.TrimSuffix(inputFile, ".webm") + ".mp4"
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-vf", "pad='width=ceil(iw/2)*2:height=ceil(ih/2)*2'", outputFile)

	var stderr strings.Builder
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("ffmpeg error: %s", stderr.String())
	}
	return outputFile, err
}

func CopyFile(src, dstDir string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dst := filepath.Join(dstDir, filepath.Base(src))
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}
	err = os.Chmod(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	return nil
}

func MoveFile(src string, dstDir string) error {
	err := os.MkdirAll(dstDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = CopyFile(src, dstDir)
	if err != nil {
		return err
	}
	err = os.Remove(src)
	if err != nil {
		return err
	}
	return nil
}

func EnsureTempDir() {
	tempDir := "./temp"
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		log.Printf("%s doesn't exist, I'll create it for you...", tempDir)
		err := os.Mkdir(tempDir, 0755)
		if err != nil {
			log.Fatalf("Error creating temp directory: %v", err)
		}
	}
}

func CleanupTempFiles() {
	tempDir := "./temp"
	files, err := os.ReadDir(tempDir)
	if err != nil {
		log.Printf("Error reading temp directory: %v", err)
		return
	}
	for _, file := range files {
		err = os.Remove(filepath.Join(tempDir, file.Name()))
		if err != nil {
			log.Printf("Error removing file %s: %v", file.Name(), err)
		}
	}
	log.Printf("Temporary files have been cleaned up.")
}
