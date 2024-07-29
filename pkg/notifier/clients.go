package notifier

import (
	"encoding/json"

	"github.com/nikoksr/notify/service/bark"
	"github.com/nikoksr/notify/service/dingding"
	po "github.com/nikoksr/notify/service/pushover"
	"github.com/nikoksr/notify/service/telegram"
	"github.com/pkg/errors"
)

type PushoverConfig struct {
	UserKey  string `json:"user_key"`
	GroupKey string `json:"group_key"`
	AppToken    string `json:"app_token"`
}

func NewPushoverClient(s string) (NotificationClient, error) {
	var cfg PushoverConfig
	if err := json.Unmarshal([]byte(s), &cfg); err != nil {
		return nil, errors.Wrap(err, "json")
	}

	c := po.New(cfg.AppToken)
	if cfg.UserKey != "" {
		c.AddReceivers(cfg.UserKey)
	}
	if cfg.GroupKey != "" {
		c.AddReceivers(cfg.GroupKey)
	}

	return &Notifier{service: c}, nil
}

type DingTalkConfig struct {
	Token  string `json:"token"`
	Secret string `json:"secret"`
}

func NewDingTalkClient(s string) (NotificationClient, error) {
	var cfg DingTalkConfig
	if err := json.Unmarshal([]byte(s), &cfg); err != nil {
		return nil, errors.Wrap(err, "json")
	}

	svc := dingding.New(&dingding.Config{
		Token:  cfg.Token,
		Secret: cfg.Secret,
	})
	return &Notifier{service: svc}, nil
}

type TelegramConfig struct {
	Token  string `json:"token"`
	ChatID int64  `json:"chat_id"`
}

func NewTelegramClient(s string) (NotificationClient, error) {
	var cfg TelegramConfig
	if err := json.Unmarshal([]byte(s), &cfg); err != nil {
		return nil, errors.Wrap(err, "json")
	}

	svc, err := telegram.New(cfg.Token)
	if err != nil {
		panic(err)
	}
	svc.AddReceivers(cfg.ChatID)
	return &Notifier{service: svc}, nil
}


type BarkConfig struct {
	DeviceKey string `json:"device_key"`
	URL string `json:"url"`
}

func NewbarkClient(s string) (NotificationClient, error) {
	var cfg BarkConfig
	if err := json.Unmarshal([]byte(s), &cfg); err != nil {
		return nil, errors.Wrap(err, "json")
	}
	url := cfg.URL
	if url == "" {
		url = bark.DefaultServerURL
	}
	b := bark.NewWithServers(cfg.DeviceKey, url)
	return &Notifier{service: b}, nil
}