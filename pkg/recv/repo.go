package recv

import "github.com/mavr/just-carrier/domain"

type Dipatcher interface {
	PushMessage(msg *domain.Message, callbackKey string) error
	PushNewChat(chat *domain.Chat) error
}

type Repository interface {
	SetSendCallback(fromUsername, toUserName string) (string, error)
}


