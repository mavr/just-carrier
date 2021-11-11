package domain

import (
	"fmt"
	"time"
)

type Message struct {
	From      int64       `json:"from" bson:"from"`
	To        string    `json:"to_recipient" bson:"to_recipient"`
	Text      string    `json:"text" bson:"text"`
	Processed bool      `json:"is_processed" bson:"is_processed"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

func (m Message) String() string {
	var txt string
	if len(m.Text) > 16 {
		txt = m.Text[:16]
	} else {
		txt = m.Text
	}
	return fmt.Sprintf("%d->%s; (%s)", m.From, m.To, txt)
}
