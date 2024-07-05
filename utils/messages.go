package utils

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Message struct {
	bot *tgbotapi.BotAPI
}

func NewMessages(bot *tgbotapi.BotAPI) *Message {
	return &Message{
		bot: bot,
	}
}

func SendMessage(bot *tgbotapi.BotAPI, chatID int64, messageText string) {
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ParseMode = "markdown"
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Error while sending reply:", err)
	}
}

func SendReply(bot *tgbotapi.BotAPI, update tgbotapi.Update, messageText string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "markdown"
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Error while sending reply:", err)
	}
}

func SendVideo(bot *tgbotapi.BotAPI, chatID int64, outputFile string) error {
	videoMsg := tgbotapi.NewVideoUpload(chatID, outputFile)
	_, err := bot.Send(videoMsg)
	if err != nil {
		log.Println("Error while sending video:", err)
		return err
	}

	return nil
}

func SendFile(bot *tgbotapi.BotAPI, chatID int64, outputFile string) error {
	fileMsg := tgbotapi.NewDocumentUpload(chatID, outputFile)
	_, err := bot.Send(fileMsg)
	if err != nil {
		log.Println("Error while sending video:", err)
		return err
	}

	return nil
}
