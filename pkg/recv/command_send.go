package recv

import (
	"strings"

	"github.com/mavr/just-carrier/domain"
	"github.com/mavr/just-carrier/pkg/staff_msg"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	tb "gopkg.in/tucnak/telebot.v3"
)

func procCommandSend(log *zap.Logger, dispatcher Dipatcher, repo Repository) tb.HandlerFunc {
	return func(ctx tb.Context) error {
		text := ctx.Message()
		if text == nil {
			return errors.New("message empty")
		}

		to, message, err := parseCommandSend(text.Text)
		if err != nil {
			tag := language.Make(ctx.Sender().LanguageCode)
			if err := ctx.Send(staff_msg.Translates.InvalidSendCommand(tag), tb.ModeMarkdown); err != nil {
				log.Error("send staff message", zap.Error(err))
			}

			log.Error("message has invalid format", zap.Error(err), zap.String("message", text.Text))

			return nil
		}

		msg := &domain.Message{
			From:      ctx.Sender().ID,
			To:        to,
			Text:      message,
			Processed: false,
			CreatedAt: ctx.Message().Time(),
		}

		callbackKey, err := repo.SetSendCallback(ctx.Sender().Username, to)
		if err != nil {
			log.Error("set callback", zap.Error(err))
		}

		if err := dispatcher.PushMessage(msg, callbackKey); err != nil {
			log.Error("push message", zap.Error(err))
		}

		return nil
	}
}

func parseCommandSend(text string) (target string, message string, err error) {
	splitText := strings.SplitN(text, " ", 3)
	if len(splitText) < 3 {
		return "", "", errors.New("invalid count")
	}

	to := splitText[1]
	if (len(to) <= 2) || ([]rune(to)[0] != '@') {
		return "", "", errors.New("invalid target format")
	}

	return string([]rune(to)[1:]), splitText[2], nil
}
