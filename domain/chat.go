package domain

type Chat struct {
	ID       int64  `json:"chat_id" bson:"chat_id"`
	UserID   int64  `json:"user_id" bson:"user_id"`
	Username string `json:"username" bson:"username"`
	LangCode string `json:"language_code" bson:"language_code"`
}
