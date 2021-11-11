package recv

import (
	"time"

	"github.com/mavr/just-carrier/pkg/staff_msg"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/text/language"
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

	repo Repository

	bot *tb.Bot
}

func New(cfg Config, log *zap.Logger, repository Repository) (*recv, error) {
	b, err := tb.NewBot(tb.Settings{
		Token: cfg.Token,
		Poller: &tb.LongPoller{
			Timeout: defaultUpdateLongPollerTimeout,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "create new tg bot")
	}

	b.Handle("/start", func(ctx tb.Context) error {
		log.Info("receive start command", zap.String("lang_code", ctx.Sender().LanguageCode))

		tag := language.Make(ctx.Sender().LanguageCode)
		return ctx.Send(staff_msg.Translates.WelcomeMessage(tag), tb.ModeMarkdown)
	})
	b.Handle("/send", procCommandSend(log, repository))

	return &recv{
		log:  log,
		repo: repository,
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
