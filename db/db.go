package db

import (
	"context"
	"encoding/json"
	"fmt"
	"polaris/ent"
	"polaris/ent/downloadclients"
	"polaris/ent/indexers"
	"polaris/ent/series"
	"polaris/ent/settings"
	"polaris/log"

	"entgo.io/ent/dialect"
	tmdb "github.com/cyruzin/golang-tmdb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Client struct {
	ent *ent.Client
}

func Open() (*Client, error) {
	client, err := ent.Open(dialect.SQLite, "file:polaris.db?cache=shared&_fk=1")
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
		return "zh-CN"
	}
	return lang
}

func (c *Client) AddWatchlist(path string, detail *tmdb.TVDetails, episodes []int) (*ent.Series, error) {
	count := c.ent.Series.Query().Where(series.TmdbID(int(detail.ID))).CountX(context.Background())
	if (count > 0) {
		return nil, fmt.Errorf("tv series %s already in watchlist", detail.Name)
	}
	r, err := c.ent.Series.Create().
		SetTmdbID(int(detail.ID)).
		SetPath(path).
		SetOverview(detail.Overview).
		SetName(detail.Name).
		SetOriginalName(detail.OriginalName).
		SetPosterPath(detail.PosterPath).
		AddEpisodeIDs(episodes...).
		Save(context.TODO())
	return r, err
}

func (c *Client) GetWatchlist() []*ent.Series {
	list, err := c.ent.Series.Query().All(context.TODO())
	if err != nil {
		log.Infof("query wtach list error: %v", err)
		return nil
	}
	return list
}

func (c *Client) SaveEposideDetail(d *ent.Episode) (int, error) {
	ep, err := c.ent.Episode.Create().
		SetAirDate(d.AirDate).
		SetSeasonNumber(d.SeasonNumber).
		SetEpisodeNumber(d.EpisodeNumber).
		SetOverview(d.Overview).
		SetTitle(d.Title).Save(context.TODO())

	return ep.ID,err
}

type TorznabSetting struct {
	URL    string `json:"url"`
	ApiKey string `json:"api_key"`
}

func (c *Client) SaveTorznabInfo(name string, setting TorznabSetting) error {
	data, err := json.Marshal(setting)
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}
	_, err = c.ent.Indexers.Create().
		SetName(name).SetImplementation(IndexerTorznabImpl).SetPriority(1).SetSettings(string(data)).Save(context.TODO())
	if err != nil {
		return errors.Wrap(err, "save db")
	}
	return nil
}

func (c *Client) GetAllTorznabInfo() map[string]TorznabSetting {
	res := c.ent.Indexers.Query().Where(indexers.Implementation(IndexerTorznabImpl)).AllX(context.TODO())
	var m = make(map[string]TorznabSetting, len(res))
	for _, r := range res {
		var ss TorznabSetting
		err := json.Unmarshal([]byte(r.Settings), &ss)
		if err != nil {
			log.Errorf("unmarshal torznab %s error: %v", r.Name, err)
			continue
		}
		m[r.Name] = ss
	}
	return m
}

func (c *Client) SaveTransmission(name, url, user, password string) error {
	_, err := c.ent.DownloadClients.Create().SetEnable(true).SetImplementation("transmission").
		SetName(name).SetURL(url).SetUser(user).SetPassword(password).Save(context.TODO())

	return err
}

func (c *Client) GetTransmission() *ent.DownloadClients {
	dc, err := c.ent.DownloadClients.Query().Where(downloadclients.Implementation("transmission")).First(context.TODO())
	if err != nil {
		log.Errorf("no transmission client found: %v", err)
		return nil
	}
	return dc
}
