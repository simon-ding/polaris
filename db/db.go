package db

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"polaris/ent"
	"polaris/ent/downloadclients"
	"polaris/ent/history"
	"polaris/ent/indexers"
	"polaris/ent/series"
	"polaris/ent/settings"
	"polaris/ent/storage"
	"polaris/log"
	"slices"
	"time"

	"entgo.io/ent/dialect"
	tmdb "github.com/cyruzin/golang-tmdb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Client struct {
	ent *ent.Client
}

func Open() (*Client, error) {
	os.Mkdir("./db", 0666)
	client, err := ent.Open(dialect.SQLite, "file:./db/polaris.db?cache=shared&_fk=1")
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

func (c *Client) AddWatchlist(storageId int, nameEn string, detail *tmdb.TVDetails, episodes []int, res ResolutionType) (*ent.Series, error) {
	count := c.ent.Series.Query().Where(series.TmdbID(int(detail.ID))).CountX(context.Background())
	if count > 0 {
		return nil, fmt.Errorf("tv series %s already in watchlist", detail.Name)
	}
	if res == "" {
		res = R1080p
	}
	if storageId == 0 {
		r, err := c.ent.Storage.Query().Where(storage.Default(true)).First(context.TODO())
		if err == nil {
			log.Infof("use default storage: %v", r.Name)
			storageId = r.ID
		}
	}

	r, err := c.ent.Series.Create().
		SetTmdbID(int(detail.ID)).
		SetStorageID(storageId).
		SetOverview(detail.Overview).
		SetName(detail.Name).
		SetNameEn(nameEn).
		SetOriginalName(detail.OriginalName).
		SetPosterPath(detail.PosterPath).
		SetAirDate(detail.FirstAirDate).
		SetResolution(res.String()).
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

type SeriesDetails struct {
	*ent.Series
	Episodes []*ent.Episode `json:"episodes"`
}

func (c *Client) GetSeriesDetails(id int) *SeriesDetails {
	se := c.ent.Series.Query().Where(series.ID(id)).FirstX(context.TODO())
	ep := se.QueryEpisodes().AllX(context.Background())
	return &SeriesDetails{
		Series:   se,
		Episodes: ep,
	}
}

func (c *Client) SaveEposideDetail(d *ent.Episode) (int, error) {
	ep, err := c.ent.Episode.Create().
		SetAirDate(d.AirDate).
		SetSeasonNumber(d.SeasonNumber).
		SetEpisodeNumber(d.EpisodeNumber).
		SetOverview(d.Overview).
		SetTitle(d.Title).Save(context.TODO())

	return ep.ID, err
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
	count := c.ent.Indexers.Query().Where(indexers.Name(name)).CountX(context.TODO())
	if count > 0 {
		c.ent.Indexers.Update().Where(indexers.Name(name)).SetSettings(string(data)).Save(context.TODO())
		return err
	}

	_, err = c.ent.Indexers.Create().
		SetName(name).SetImplementation(IndexerTorznabImpl).SetPriority(1).SetSettings(string(data)).Save(context.TODO())
	if err != nil {
		return errors.Wrap(err, "save db")
	}

	return nil
}

func (c *Client) DeleteTorznab(id int) {
	c.ent.Indexers.Delete().Where(indexers.ID(id)).Exec(context.TODO())
}

type TorznabInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	TorznabSetting
}

func (c *Client) GetAllTorznabInfo() []*TorznabInfo {
	res := c.ent.Indexers.Query().Where(indexers.Implementation(IndexerTorznabImpl)).AllX(context.TODO())

	var l = make([]*TorznabInfo, 0, len(res))
	for _, r := range res {
		var ss TorznabSetting
		err := json.Unmarshal([]byte(r.Settings), &ss)
		if err != nil {
			log.Errorf("unmarshal torznab %s error: %v", r.Name, err)
			continue
		}
		l = append(l, &TorznabInfo{
			ID:             r.ID,
			Name:           r.Name,
			TorznabSetting: ss,
		})
	}
	return l
}

func (c *Client) SaveTransmission(name, url, user, password string) error {
	count := c.ent.DownloadClients.Query().Where(downloadclients.Name(name)).CountX(context.TODO())
	if count != 0 {
		err := c.ent.DownloadClients.Update().Where(downloadclients.Name(name)).
			SetURL(url).SetUser(user).SetPassword(password).Exec(context.TODO())
		return err
	}

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

func (c *Client) GetAllDonloadClients() []*ent.DownloadClients {
	cc, err := c.ent.DownloadClients.Query().All(context.TODO())
	if err != nil {
		log.Errorf("no download client")
		return nil
	}
	return cc
}

func (c *Client) DeleteDownloadCLient(id int) {
	c.ent.DownloadClients.Delete().Where(downloadclients.ID(id)).Exec(context.TODO())
}

// Storage is the model entity for the Storage schema.
type StorageInfo struct {
	Name           string `json:"name"`
	Implementation string `json:"implementation"`
	Path           string `json:"path"`
	User           string `json:"user"`
	Password       string `json:"password"`
	Default        bool   `json:"default"`
}

func (c *Client) AddStorage(s StorageInfo) error {
	if !slices.Contains(StorageImplementations(), s.Implementation) {
		return fmt.Errorf("implementation not supported: %v", s.Implementation)
	}
	count := c.ent.Storage.Query().Where(storage.Name(s.Name)).CountX(context.TODO())
	if count > 0 {
		//storage already exist, edit exist one
		return c.ent.Storage.Update().Where(storage.Name(s.Name)).
			SetImplementation(s.Implementation).
			SetPath(s.Path).
			SetUser(s.User).
			SetDefault(s.Default).
			SetPassword(s.Password).Exec(context.TODO())
	}
	countAll := c.ent.Storage.Query().CountX(context.TODO())
	if countAll == 0 {
		log.Infof("first storage, make it default: %s", s.Name)
		s.Default = true
	}
	_, err := c.ent.Storage.Create().SetName(s.Name).
		SetImplementation(s.Implementation).
		SetPath(s.Path).
		SetUser(s.User).
		SetDefault(s.Default).
		SetPassword(s.Password).Save(context.TODO())
	if err != nil {
		return err
	}
	if s.Default {
		return c.SetDefaultStorageByName(s.Name)
	}
	return nil
}

func (c *Client) GetAllStorage() []*ent.Storage {
	data, err := c.ent.Storage.Query().Where(storage.Deleted(false)).All(context.TODO())
	if err != nil {
		log.Errorf("get storage: %v", err)
		return nil
	}
	return data
}

func (c *Client) GetStorage(id int) *ent.Storage {
	r, err := c.ent.Storage.Query().Where(storage.ID(id)).First(context.TODO())
	if err != nil {
		//use default storage
		return c.ent.Storage.Query().Where(storage.Default(true)).FirstX(context.TODO())
	}
	return r
}

func (c *Client) DeleteStorage(id int) error {
	return c.ent.Storage.Update().Where(storage.ID(id)).SetDeleted(true).Exec(context.TODO())
}

func (c *Client) SetDefaultStorage(id int) error {
	err := c.ent.Storage.Update().Where(storage.ID(id)).SetDefault(true).Exec(context.TODO())
	if err != nil {
		return err
	}
	err = c.ent.Storage.Update().Where(storage.Or(storage.ID(id))).SetDefault(false).Exec(context.TODO())
	return err
}

func (c *Client) SetDefaultStorageByName(name string) error {
	err := c.ent.Storage.Update().Where(storage.Name(name)).SetDefault(true).Exec(context.TODO())
	if err != nil {
		return err
	}
	err = c.ent.Storage.Update().Where(storage.Or(storage.Name(name))).SetDefault(false).Exec(context.TODO())
	return err
}


func (c *Client) SaveHistoryRecord(h ent.History) (*ent.History,error) {
	return c.ent.History.Create().SetSeriesID(h.SeriesID).SetEpisodeID(h.EpisodeID).SetDate(time.Now()). 
		SetCompleted(h.Completed).SetTargetDir(h.TargetDir).SetSourceTitle(h.SourceTitle).SetSaved(h.Saved).Save(context.TODO())
}

func (c *Client) SetHistoryComplete(id int) error {
	return c.ent.History.Update().Where(history.ID(id)).SetCompleted(true).Exec(context.TODO())
}

func (c *Client) GetHistories() ent.Histories {
	h, err := c.ent.History.Query().All(context.TODO())
	if err != nil {
		return nil
	}
	return h
}

func (c *Client) GetHistory(id int) *ent.History {
	return c.ent.History.Query().Where(history.ID(id)).FirstX(context.TODO())
}


func (c *Client) DeleteHistory(id int) error {
	_, err := c.ent.History.Delete().Where(history.ID(id)).Exec(context.Background())
	return err
}


func (c *Client) GetDownloadDir() string {
	r, err := c.ent.Settings.Query().Where(settings.Key(SettingDownloadDir)).First(context.TODO())
	if err != nil {
		return "/downloads"
	}
	return r.Value
}