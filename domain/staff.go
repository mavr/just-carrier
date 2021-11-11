package domain

import "golang.org/x/text/language"

type Staff interface {
	WelcomeMessage(tag language.Tag) string
	InvalidSendCommand(tag language.Tag) string
}
