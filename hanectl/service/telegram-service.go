package service

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/utils"
	"strconv"
)

type TelegramService struct {
	config      config.IConfiguration
	telegramBot model.ITelegramBotAPI // *tgbotapi.BotAPI
}

func (t *TelegramService) connect(token string) {
	t.telegramBot = nil
	if utils.IsNotBlank(token) {
		bot, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			log.Error().Msgf("Could not create telegram bot: %v", err)
			bot = nil
		} else {
			log.Info().Msg("Telegram bot created")
		}
		t.telegramBot = bot
	}
}

func (t *TelegramService) SendMessage(message string) bool {
	result := false
	if t.telegramBot != nil {
		chatIds := t.config.GetList(config.TelegramChatIds)
		for _, chatId := range chatIds {
			telegramId, err := strconv.ParseInt(chatId.(string), 10, 64)
			if err == nil {
				log.Debug().Msgf("Sending telegram message %s", message)
				msg := tgbotapi.NewMessage(telegramId, message)
				_, err := t.telegramBot.Send(msg)
				if err != nil {
					result = false
					log.Error().Msgf("Could not send msg to telegram chat:%d err:%v", chatId, err)
				} else {
					result = true
				}
			} else {
				log.Error().Msgf("Could not convert chatId %d: %v", chatId, err)
			}
		}
	}
	return result
}

func getTelegramToken(cfg config.IConfiguration) string {
	botEnabled := cfg.GetBool(config.TelegramBotEnabled, false)
	if botEnabled {
		return cfg.GetStr(config.TelegramBotToken, "")
	}
	return ""
}

func NewTelegramService(cfg config.IConfiguration) model.ITelegramService {

	service := &TelegramService{
		config: cfg,
	}
	service.connect(getTelegramToken(cfg))
	return service
}
