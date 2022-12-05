package telegram

import (
	"checker/config"
	"checker/logger"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot  *tgbotapi.BotAPI
	conf *config.Config
	log  logger.Log
}

func NewTg(conf *config.Config, log logger.Log) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(conf.TokenTg)
	if err != nil {
		return nil, err
	}
	log.Info("Start telegram bot")
	return &Bot{bot: bot, conf: conf, log: log}, err
}
func (bot *Bot) ReConnect() error {
	b, err := tgbotapi.NewBotAPI(bot.conf.TokenTg)
	if err != nil {
		return err
	}
	bot.bot = b
	return err
}

func (bot *Bot) Send(v interface{}) error {

	text, err := json.Marshal(v)
	if err != nil {
		return err
	}
	mesconf := tgbotapi.NewMessage(bot.conf.IdTelegram, string(text))
	_, err = bot.bot.Send(mesconf)
	return nil
}
