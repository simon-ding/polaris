package notifier

import (
	"context"
	"polaris/pkg/utils"

	"github.com/nikoksr/notify"
)

type HandlerFunc func(string) (NotificationClient, error)

type NotificationClient interface {
	SendMsg(msg string) error
}

type Notifier struct {
	service notify.Notifier
}

func (s *Notifier) SendMsg(msg string) error {
	notifier := notify.New()
	notifier.UseServices(s.service)
	return notifier.Send(context.TODO(), "Polaris", msg)
}

var handler = utils.Map[string, HandlerFunc]{}

func init() {
	handler.Store("pushover", NewPushoverClient)
	handler.Store("dingtalk", NewDingTalkClient)
	handler.Store("telegram", NewTelegramClient)
	handler.Store("bark", NewbarkClient)
}

func Gethandler(name string) (HandlerFunc, bool) {
	return handler.Load(name)
}
