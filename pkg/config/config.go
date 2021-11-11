package config

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/mavr/just-carrier/pkg/http_server"
	"github.com/mavr/just-carrier/pkg/infra/redis"
	"github.com/mavr/just-carrier/pkg/infra/rmq"
	"github.com/mavr/just-carrier/pkg/recv"
)

type Config interface {
	Rabbit() rmq.Config
	Telegram() recv.Config
	HTTPApi() http_server.APIServiceConfig
	Redis() redis.Config
}

// Config containing configuration values
type config struct {
	// Debug set logs level output
	App struct {
		Debug bool `toml:"debug"`
		Port  int  `toml:"port"`
	} `toml:"application"`

	// Telegram bot token
	Bot struct {
		TGBotToken string `toml:"bot_token"`
	} `toml:"telegram"`

	// Rabbit
	RMQ struct {
		ConnectionString   string `toml:"connection_string"`
		ExchangeNewChat    string `toml:"new_chat_notification_exchange"`
		ExchangeNewMessage string `toml:"new_message_notification"`
	} `toml:"rabbit"`

	// Redis
	RedisDatabase struct {
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
		DB       int    `toml:"db"`
		Username string `toml:"username"`
		Password string `toml:"password"`
	} `toml:"redis"`
}

func New(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var c config
	if err := toml.Unmarshal(buf, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *config) Rabbit() rmq.Config {
	return rmq.Config{
		ConnectionString:          c.RMQ.ConnectionString,
		ExchangeNewChatNotificate: c.RMQ.ExchangeNewChat,
		ExchangeMessage:           c.RMQ.ExchangeNewMessage,
	}
}

func (c *config) Telegram() recv.Config {
	return recv.Config{
		Token: c.Bot.TGBotToken,
	}
}

func (c *config) HTTPApi() http_server.APIServiceConfig {
	return http_server.APIServiceConfig{
		AppRevision:  "0.0.1",
		AppVersion:   "0.0.1",
		AppDebugMode: c.App.Debug,
		ServPort:     c.App.Port,
	}
}

func (c *config) Redis() redis.Config {
	return redis.Config{
		Host:     c.RedisDatabase.Host,
		Port:     c.RedisDatabase.Port,
		Username: c.RedisDatabase.Username,
		Password: c.RedisDatabase.Password,
		Database: c.RedisDatabase.DB,
	}
}