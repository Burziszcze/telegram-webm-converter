package main

import (
	"log"

	"github.com/BurntSushi/toml"
	utils "github.com/Burziszcze/telegram-webm-converter/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/fsnotify.v1"
)

type Config struct {
	Telegram BotConfig
}

type BotConfig struct {
	APIKey string
}

var config Config

func main() {
	if err := loadConfig(); err != nil {
		log.Fatalf("Error reading initial configuration: %v", err)
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("BError initializing file monitor: %v", err)
	}
	bot, err := tgbotapi.NewBotAPI(config.Telegram.APIKey)
	defer watcher.Close()

	if err := watcher.Add("config.toml"); err != nil {
		log.Fatalf("Error adding file to monitor: %v", err)
	}
	if err != nil {
		log.Fatalf("Unable to create Telegram Bot API client: %v", err)
	}
	log.Printf("Bot launched: %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	converter := utils.NewWebmConverter(bot)
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Unable to create update channel %v", err)
	}
	go watchConfigFile(watcher)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		converter.HandleMessage(update.Message)
	}
}

func loadConfig() error {
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		return err
	}

	return nil
}

func watchConfigFile(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if err := loadConfig(); err != nil {
					log.Printf("Error while reading changed configuration: %v", err)
				}
				log.Println("Updated configuration from file.")
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("File monitoring error: %v", err)
		}
	}
}
