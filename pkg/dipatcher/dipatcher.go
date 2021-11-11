package dipatcher

import (
	"encoding/json"

	"github.com/mavr/just-carrier/domain"
)

type newMessage struct {
	From      int64  `json:"from" bson:"from"`
	To        string `json:"to_recipient" bson:"to_recipient"`
	Text      string `json:"text" bson:"text"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	Callback  string `json:"callback_key" bson:"callback_key"`
}

type dispatcher struct {
	cliMsg  domain.Publisher
	cliChat domain.Publisher
}

func New(message, chat domain.Publisher) (*dispatcher, error) {
	return &dispatcher{
		cliMsg:  message,
		cliChat: chat,
	}, nil
}

func (d *dispatcher) PushMessage(msg *domain.Message, callbackKey string) error {
	m := newMessage{
		From:      msg.From,
		To:        msg.To,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt.Unix(),
		Callback: callbackKey,
	}
	raw, _ := json.Marshal(m)
	return d.cliMsg.Publish(raw)
}

func (d *dispatcher) PushNewChat(chat *domain.Chat) error {
	raw, _ := json.Marshal(chat)
	return d.cliChat.Publish(raw)
}
