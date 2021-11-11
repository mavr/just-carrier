package recv

import (
	"github.com/mavr/just-carrier/domain"
	"github.com/mavr/just-carrier/pkg/staff_msg"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	tb "gopkg.in/tucnak/telebot.v3"
)

func procCommandStart(log *zap.Logger, repo Dipatcher) tb.HandlerFunc {
	return func(ctx tb.Context) error {
		log.Info("receive start command", zap.String("username", ctx.Sender().Username))

		tag := language.Make(ctx.Sender().LanguageCode)
		chat := &domain.Chat{
			ID:       ctx.Chat().ID,
			Username: ctx.Sender().Username,
			UserID:   ctx.Sender().ID,
			LangCode: tag.String(),
		}

		if err := repo.PushNewChat(chat); err != nil {
			log.Error("push new chat notification", zap.Error(err))
		}

		return ctx.Send(staff_msg.Translates.WelcomeMessage(tag), tb.ModeMarkdown)
	}
}
