package recv

import (
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v3"
)

const (
	// Default buffer size for request to telegam api
	defaultUpdateLongPollerTimeout = 30 * time.Second
	defaultUpdatesBufferSize       = 4
	defaultUpdateTimer             = 3 * time.Second
)

type Config struct {
	Token string
}

type recv struct {
	log *zap.Logger

	repo Dipatcher

	bot *tb.Bot
}

func New(cfg Config, log *zap.Logger, dispatcher Dipatcher, repository Repository) (*recv, error) {
	b, err := tb.NewBot(tb.Settings{
		Token: cfg.Token,
		Poller: &tb.LongPoller{
			Timeout: defaultUpdateLongPollerTimeout,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "create new tg bot")
	}

	b.Handle("/start", procCommandStart(log, dispatcher))
	b.Handle("/send", procCommandSend(log, dispatcher, repository))

	return &recv{
		log:  log,
		repo: dispatcher,
		bot:  b,
	}, nil
}

func (r *recv) Start() error {
	go r.bot.Start()
	r.log.Info("receiver started")
	return nil
}

func (r *recv) Stop() error {
	r.bot.Stop()
	r.log.Info("receiver stopped")
	return nil
}
