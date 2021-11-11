package staff_msg

import "golang.org/x/text/language"

var (
	Translates *ts
)

func init() {
	Translates = &ts{}
}

type translate struct {
	m map[language.Tag]string
}

func (t *translate) Get(tag language.Tag) string {
	if r, ok := t.m[tag]; ok {
		return r
	}
	return t.m[language.English]
}

type ts struct{}

func (s *ts) WelcomeMessage(tag language.Tag) string {
	return msgWelcome.Get(tag)
}

func (s *ts) InvalidSendCommand(tag language.Tag) string {
	return msgInvalidSendCommand.Get(tag)
}
