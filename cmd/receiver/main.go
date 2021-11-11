package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/mavr/just-carrier/pkg/config"
	"github.com/mavr/just-carrier/pkg/dipatcher"
	"github.com/mavr/just-carrier/pkg/http_server"
	"github.com/mavr/just-carrier/pkg/infra/rmq"
	"github.com/mavr/just-carrier/pkg/recv"
	"go.uber.org/zap"
)

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	cfg, err := config.New("conf/config.toml")
	if err != nil {
		log.Error("create config failed", zap.Error(err))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	amqp := rmq.New(cfg.Rabbit(), log.Named("rabbitmq"))
	pubMessage, err := amqp.NewPublisher(cfg.Rabbit().ExchangeMessage, cfg.Rabbit().ExchangeMessage, "")
	if err != nil {
		log.Error("create publisher", zap.Error(err))
		return
	}
	pubChat, err := amqp.NewPublisher(cfg.Rabbit().ExchangeNewChatNotificate, cfg.Rabbit().ExchangeNewChatNotificate, "")
	if err != nil {
		log.Error("create publisher", zap.Error(err))
		return
	}

	dispatcher, err := dipatcher.New(pubMessage, pubChat)
	if err != nil {
		log.Error("create dispatcher", zap.Error(err))
		return
	}

	receiver, err := recv.New(cfg.Telegram(), log, dispatcher)
	if err != nil {
		log.Error("create receiver", zap.Error(err))
		return
	}
	_ = receiver.Start()

	api := http_server.NewAPIService(cfg.HTTPApi(), log)
	_ = api.Run(ctx)

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt)

	_ = <-signalCh

	if err := receiver.Stop(); err != nil {
		log.Error("receiver stop", zap.Error(err))
	}

	cancel()
}
