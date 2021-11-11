package config

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/mavr/just-carrier/pkg/http_server"
	"github.com/mavr/just-carrier/pkg/infra/rmq"
	"github.com/mavr/just-carrier/pkg/recv"
)

type Config interface {
	Rabbit() rmq.Config
	Telegram() recv.Config
	HTTPApi() http_server.APIServiceConfig
}

// Config containing configuration values
type config struct {
	// Debug set logs level output
	App struct {
		Debug bool `toml:"debug"`
		Port  int  `toml:"port"`

		ExchangeNewChat    string `toml:"new_chat_notification_exchange"`
		ExchangeNewMessage string `toml:"new_message_notification"`
	}

	// Telegram bot token
	Bot struct {
		TGBotToken string `toml:"telegram_bot_token"`
	}

	// Rabbit
	RMQ struct {
		ConnectionString string `toml:"connection_string"`
	}
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
		ConnectionString: c.RMQ.ConnectionString,
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
