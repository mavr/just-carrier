package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database int
}

type store struct {
	ctx context.Context
	db  *redis.Client
}

func New(ctx context.Context, cfg Config) (*store, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username: "",
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	return &store{
		ctx: ctx,
		db: rdb,
	}, nil
}

type callback struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func (s *store) SetSendCallback(fromUsername, toUserName string) (string, error) {
	uid := uuid.NewV4().String()
	payload := callback{
		From: fromUsername,
		To:   toUserName,
	}
	raw, _ := json.Marshal(payload)

	return uid, s.db.Set(s.ctx, uid, raw, 0).Err()
}
