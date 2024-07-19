package db

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"polaris/ent"
	"polaris/ent/downloadclients"
	"polaris/ent/episode"
	"polaris/ent/history"
	"polaris/ent/indexers"
	"polaris/ent/media"
	"polaris/ent/settings"
	"polaris/ent/storage"
	"polaris/log"
	"polaris/pkg/utils"
	"strings"
	"time"

	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Client struct {
	ent *ent.Client
}

func Open() (*Client, error) {
	os.Mkdir(DataPath, 0777)
	client, err := ent.Open(dialect.SQLite, fmt.Sprintf("file:%v/polaris.db?cache=shared&_fk=1", DataPath))
	if err != nil {
		return nil, errors.Wrap(err, "failed opening connection to sqlite")
	}
	//defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, errors.Wrap(err, "failed creating schema resources")
	}
	c := &Client{
		ent: client,
	}

	return c, nil
}

func (c *Client) init() {
	c.generateJwtSerectIfNotExist()
	if err := c.generateDefaultLocalStorage(); err != nil {
		log.Errorf("generate default storage: %v", err)
	}

	downloadDir := c.GetSetting(SettingDownloadDir)
	if downloadDir == "" {
		log.Infof("set default download dir")
		c.SetSetting(downloadDir, "/downloads")
	}
}

func (c *Client) generateJwtSerectIfNotExist() {
	v := c.GetSetting(JwtSerectKey)
	if v == "" {
		log.Infof("generate jwt serect")
		key := utils.RandString(32)
		c.SetSetting(JwtSerectKey, key)
	}
}

func (c *Client) generateDefaultLocalStorage() error {
	n, _ := c.ent.Storage.Query().Count(context.TODO())
	if n != 0 {
		return nil
	}
	log.Infof("add default storage")
	return c.AddStorage(&StorageInfo{
		Name:           "local",
		Implementation: "local",
		Default:        true,
		Settings: map[string]string{
			"tv_path":    "/data/tv",
			"movie_path": "/data/movies",
		},
	})
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

func (c *Client) AddMediaWatchlist(m *ent.Media, episodes []int) (*ent.Media, error) {
	count := c.ent.Media.Query().Where(media.TmdbID(m.TmdbID)).CountX(context.Background())
	if count > 0 {
		return nil, fmt.Errorf("tv series %s already in watchlist", m.NameEn)
	}

	if m.StorageID == 0 {
		r, err := c.ent.Storage.Query().Where(storage.And(storage.Default(true), storage.Deleted(false))).First(context.TODO())
		if err == nil {
			log.Infof("use default storage: %v", r.Name)
			m.StorageID = r.ID
		}
	}
	r, err := c.ent.Media.Create().
		SetTmdbID(m.TmdbID).
		SetStorageID(m.StorageID).
		SetOverview(m.Overview).
		SetNameCn(m.NameCn).
		SetNameEn(m.NameEn).
		SetOriginalName(m.OriginalName).
		SetMediaType(m.MediaType).
		SetAirDate(m.AirDate).
		SetResolution(m.Resolution).
		SetTargetDir(m.TargetDir).
		AddEpisodeIDs(episodes...).
		Save(context.TODO())
	return r, err

}

func (c *Client) GetMediaWatchlist(mediaType media.MediaType) []*ent.Media {
	list, err := c.ent.Media.Query().Where(media.MediaTypeEQ(mediaType)).All(context.TODO())
	if err != nil {
		log.Infof("query wtach list error: %v", err)
		return nil
	}
	return list
}

func (c *Client) GetEpisode(seriesId, seasonNum, episodeNum int) (*ent.Episode, error) {
	return c.ent.Episode.Query().Where(episode.MediaID(seriesId), episode.SeasonNumber(seasonNum),
		episode.EpisodeNumber(episodeNum)).First(context.TODO())
}

func (c *Client) UpdateEpiode(episodeId int, name, overview string) error {
	return c.ent.Episode.Update().Where(episode.ID(episodeId)).SetTitle(name).SetOverview(overview).Exec(context.TODO())
}

type MediaDetails struct {
	*ent.Media
	Episodes []*ent.Episode `json:"episodes"`
}

func (c *Client) GetMediaDetails(id int) *MediaDetails {
	se, err := c.ent.Media.Query().Where(media.ID(id)).First(context.TODO())
	if err != nil {
		log.Errorf("get series %d: %v", id, err)
		return nil
	}
	var md = &MediaDetails{
		Media: se,
	}
	if se.MediaType == media.MediaTypeTv {
		ep, err := se.QueryEpisodes().All(context.Background())
		if err != nil {
			log.Errorf("get episodes %d: %v", id, err)
			return nil
		}
		md.Episodes = ep
	}
	return md
}

func (c *Client) DeleteMedia(id int) error {
	_, err := c.ent.Episode.Delete().Where(episode.MediaID(id)).Exec(context.TODO())
	if err != nil {
		return err
	}
	_, err = c.ent.Media.Delete().Where(media.ID(id)).Exec(context.TODO())
	return err
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
	Name           string            `json:"name"`
	Implementation string            `json:"implementation"`
	Settings       map[string]string `json:"settings"`
	Default        bool              `json:"default"`
}

type LocalDirSetting struct {
	TvPath    string `json:"tv_path"`
	MoviePath string `json:"movie_path"`
}

type WebdavSetting struct {
	URL       string `json:"url"`
	TvPath    string `json:"tv_path"`
	MoviePath string `json:"movie_path"`
	User      string `json:"user"`
	Password  string `json:"password"`
}

func (c *Client) AddStorage(st *StorageInfo) error {
	if !strings.HasSuffix(st.Settings["tv_path"], "/") {
		st.Settings["tv_path"] += "/"
	}
	if !strings.HasSuffix(st.Settings["movie_path"], "/") {
		st.Settings["movie_path"] += "/"
	}


	data, err := json.Marshal(st.Settings)
	if err != nil {
		return errors.Wrap(err, "json marshal")
	}

	count := c.ent.Storage.Query().Where(storage.Name(st.Name)).CountX(context.TODO())
	if count > 0 {
		//storage already exist, edit exist one
		return c.ent.Storage.Update().Where(storage.Name(st.Name)).
			SetImplementation(storage.Implementation(st.Implementation)).
			SetSettings(string(data)).Exec(context.TODO())
	}
	countAll := c.ent.Storage.Query().Where(storage.Deleted(false)).CountX(context.TODO())
	if countAll == 0 {
		log.Infof("first storage, make it default: %s", st.Name)
		st.Default = true
	}
	_, err = c.ent.Storage.Create().SetName(st.Name).
		SetImplementation(storage.Implementation(st.Implementation)).
		SetSettings(string(data)).SetDefault(st.Default).Save(context.TODO())
	if err != nil {
		return err
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

type Storage struct {
	ent.Storage
}

func (s *Storage) ToLocalSetting() LocalDirSetting {
	if s.Implementation != storage.ImplementationLocal {
		panic("not local storage")
	}
	var localSetting LocalDirSetting
	json.Unmarshal([]byte(s.Settings), &localSetting)
	return localSetting
}

func (s *Storage) ToWebDavSetting() WebdavSetting {
	if s.Implementation != storage.ImplementationWebdav {
		panic("not webdav storage")
	}
	var webdavSetting WebdavSetting
	json.Unmarshal([]byte(s.Settings), &webdavSetting)
	return webdavSetting
}

func (s *Storage) GetPath() (tvPath string, moviePath string) {
	var m map[string]string
	json.Unmarshal([]byte(s.Settings), &m)
	return m["tv_path"], m["movie_path"]
}

func (c *Client) GetStorage(id int) *Storage {
	r, err := c.ent.Storage.Query().Where(storage.ID(id)).First(context.TODO())
	if err != nil {
		//use default storage
		r := c.ent.Storage.Query().Where(storage.Default(true)).FirstX(context.TODO())
		return &Storage{*r}
	}
	return &Storage{*r}
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

func (c *Client) SaveHistoryRecord(h ent.History) (*ent.History, error) {
	return c.ent.History.Create().SetMediaID(h.MediaID).SetEpisodeID(h.EpisodeID).SetDate(time.Now()).
		SetStatus(h.Status).SetTargetDir(h.TargetDir).SetSourceTitle(h.SourceTitle).SetSaved(h.Saved).Save(context.TODO())
}

func (c *Client) SetHistoryStatus(id int, status history.Status) error {
	return c.ent.History.Update().Where(history.ID(id)).SetStatus(status).Exec(context.TODO())
}

func (c *Client) GetHistories() ent.Histories {
	h, err := c.ent.History.Query().All(context.TODO())
	if err != nil {
		return nil
	}
	return h
}

func (c *Client) GetRunningHistories() ent.Histories {
	h, err := c.ent.History.Query().Where(history.Or(history.StatusEQ(history.StatusRunning),
		history.StatusEQ(history.StatusUploading))).All(context.TODO())
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

func (c *Client) UpdateEpisodeFile(mediaID int, seasonNum, episodeNum int, file string) error {
	ep, err := c.ent.Episode.Query().Where(episode.MediaID(mediaID)).Where(episode.EpisodeNumber(episodeNum)).
		Where(episode.SeasonNumber(seasonNum)).First(context.TODO())
	if err != nil {
		return errors.Wrap(err, "finding episode")
	}
	return ep.Update().SetFileInStorage(file).Exec(context.TODO())
}

func (c *Client) SetEpisodeStatus(id int, status episode.Status) error {
	return c.ent.Episode.Update().Where(episode.ID(id)).SetStatus(status).Exec(context.TODO())
}
