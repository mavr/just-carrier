package tg

import (
	"net/http"
	"time"

	"github.com/pkg/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Config struct {
	Token string
}

type bot struct {
	*tgbotapi.BotAPI
}

func New(cfg Config) (*bot, error) {
	httpCli := &http.Client{
		Timeout: time.Second * 60,
	}
	tgCli, err := tgbotapi.NewBotAPIWithClient(cfg.Token, httpCli)
	if err != nil {
		return nil, errors.Wrap(err, "create new telegram client")
	}

	bot := &bot{
		tgCli,
	}

	return bot, nil
}
