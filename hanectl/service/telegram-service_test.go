package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/test/mock_test"
	"testing"
)

func TestSendMessage(t *testing.T) {

	chatIds := []interface{}{"12345"}

	cfgMock := new(mock_test.ConfigurationMock)
	botMock := new(mock_test.TelegramBotApiMock)
	srv := TelegramService{
		config: cfgMock,
		telegramBot: botMock,
	}
	cfgMock.On("GetList", config.TelegramChatIds).Return(chatIds)
	botMock.On("Send", mock.AnythingOfTypeArgument("tgbotapi.MessageConfig")).Return(tgbotapi.Message{}, nil)

	send := srv.SendMessage("message")
	assert.True(t, send)

	botMock.AssertExpectations(t)
	cfgMock.AssertExpectations(t)
}
