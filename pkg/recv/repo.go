package recv

import "github.com/mavr/just-carrier/domain"

type Repository interface {
	PushMessage(*domain.Message) error
	PushNewChat(chat *domain.Chat) error
}
