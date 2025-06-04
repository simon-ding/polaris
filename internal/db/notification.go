package db

import (
	"context"
	"encoding/json"
	"polaris/ent"
	"polaris/ent/notificationclient"
	"polaris/pkg/notifier"
	"strings"

	"github.com/pkg/errors"
)

func (c *client) GetAllNotificationClients2() ([]*ent.NotificationClient, error) {
	return c.ent.NotificationClient.Query().All(context.TODO())
}

func (c *client) GetAllNotificationClients() ([]*NotificationClient, error) {
	all, err := c.ent.NotificationClient.Query().All(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "query db")
	}
	var all1 []*NotificationClient
	for _, item := range all {
		cl, err := toNotificationClient(item)
		if err != nil {
			return nil, errors.Wrap(err, "convert")
		}
		all1 = append(all1, cl)
	}
	return all1, nil
}

func (c *client) AddNotificationClient(name, service string, setting string, enabled bool) error {
	// data, err := json.Marshal(setting)
	// if err != nil {
	// 	return errors.Wrap(err, "json")
	// }
	service = strings.ToLower(service)
	count, err := c.ent.NotificationClient.Query().Where(notificationclient.Name(name)).Count(context.Background())
	if err == nil && count > 0 {
		//update exist one
		return c.ent.NotificationClient.Update().Where(notificationclient.Name(name)).SetService(service).
			SetSettings(setting).SetEnabled(enabled).Exec(context.Background())
	}

	return c.ent.NotificationClient.Create().SetName(name).SetService(service).
		SetSettings(setting).SetEnabled(enabled).Exec(context.Background())
}

func (c *client) DeleteNotificationClient(id int) error {
	_, err := c.ent.NotificationClient.Delete().Where(notificationclient.ID(id)).Exec(context.Background())
	return err
}

func (c *client) GetNotificationClient(id int) (*NotificationClient, error) {
	noti, err := c.ent.NotificationClient.Query().Where(notificationclient.ID(id)).First(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}

	return toNotificationClient(noti)
}

func toNotificationClient(cl *ent.NotificationClient) (*NotificationClient, error) {
	var settings interface{}
	switch cl.Service {
	case "pushover":
		settings = notifier.PushoverConfig{}
	case "dingtalk":
		settings = notifier.DingTalkConfig{}
	case "telegram":
		settings = notifier.TelegramConfig{}
	case "bark":
		settings = notifier.BarkConfig{}
	case "serverchan":
		settings = notifier.ServerChanConfig{}
	}
	err := json.Unmarshal([]byte(cl.Settings), &settings)
	if err != nil {
		return nil, errors.Wrap(err, "json")
	}
	return &NotificationClient{
		ID:       cl.ID,
		Name:     cl.Name,
		Service:  cl.Service,
		Enabled:  cl.Enabled,
		Settings: settings,
	}, nil

}

type NotificationClient struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Service  string      `json:"service"`
	Enabled  bool        `json:"enabled"`
	Settings interface{} `json:"settings"`
}
