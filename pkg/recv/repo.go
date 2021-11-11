package recv

import "github.com/mavr/just-carrier/domain"

type Repository interface {
	PushMessage(*domain.Message) error
	//SetChat(c *domain.Chat) error
}
