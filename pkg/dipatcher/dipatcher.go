package dipatcher

import (
	"encoding/json"

	"github.com/mavr/just-carrier/domain"
)

type dispatcher struct {
	cli domain.Publisher
}

func New(publisher domain.Publisher) (*dispatcher, error) {
	return &dispatcher{
		cli: publisher,
	}, nil
}

func (d *dispatcher) PushMessage(msg *domain.Message) error {
	raw, _ := json.Marshal(msg)
	return d.cli.Publish(raw)
}
