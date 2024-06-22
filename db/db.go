package db

import (
	"context"
	"polaris/ent"
	"polaris/ent/settings"
	"polaris/log"

	tmdb "github.com/cyruzin/golang-tmdb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Client struct {
	ent *ent.Client
}

func Open() (*Client, error) {
	client, err := ent.Open("sqlite3", "file:polaris.db?cache=shared&_fk=1")
	if err != nil {
		return nil, errors.Wrap(err, "failed opening connection to sqlite")
	}
	//defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, errors.Wrap(err, "failed creating schema resources")
	}
	return &Client{
		ent: client,
	}, nil
}

func (c *Client) GetSetting(key string) string {
	v, err := c.ent.Settings.Query().Where(settings.Key(key)).Only(context.TODO())
	if err != nil {
		log.Errorf("get setting by key: %s error: %v", key, err)
		return ""
	}
	return v.Value
}

func (c *Client) SetSetting(key, value string) error {
	v, err := c.ent.Settings.Query().Where(settings.Key(key)).Only(context.TODO())
	if err != nil {
		log.Infof("create new setting")
		_, err := c.ent.Settings.Create().SetKey(key).SetValue(value).Save(context.TODO())
		return err
	}
	_, err = c.ent.Settings.UpdateOneID(v.ID).SetValue(value).Save(context.TODO())
	return err
}

func (c *Client) GetLanguage() string {
	lang := c.GetSetting(SettingLanguage)
	log.Infof("get application language: %s", lang)
	if lang == "" {
		return "zh_CN"
	}
	return lang
}

func (c *Client) AddWatchlist(path string, detail *tmdb.TVDetails) error {
	_, err := c.ent.Series.Create().
		SetTmdbID(int(detail.ID)). 
		SetPath(path). 
		SetOverview(detail.Overview). 
		SetTitle(detail.Name). 
		SetOriginalName(detail.OriginalName). 
		Save(context.TODO())
	return err
}

func (c *Client) GetWatchlist() []*ent.Series {
	list, err := c.ent.Series.Query().All(context.TODO())
	if err != nil {
		log.Infof("query wtach list error: %v", err)
		return nil
	}
	return list
}
