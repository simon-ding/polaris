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
	"polaris/ent/importlist"
	"polaris/ent/indexers"
	"polaris/ent/media"
	"polaris/ent/schema"
	"polaris/ent/settings"
	"polaris/ent/storage"
	"polaris/log"
	"polaris/pkg/utils"
	"strings"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
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
	c.init()

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
		c.SetSetting(SettingDownloadDir, "/downloads")
	}
	logLevel := c.GetSetting(SettingLogLevel)
	if logLevel == "" {
		log.Infof("set default log level")
		c.SetSetting(SettingLogLevel, "info")
	}
	if tr := c.GetAllDonloadClients(); len(tr) == 0 {
		log.Warnf("no download client, set default download client")
		c.SaveDownloader(&ent.DownloadClients{
			Name:           "transmission",
			Implementation: downloadclients.ImplementationTransmission,
			URL:            "http://transmission:9091",
		})
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
		TvPath:         "/data/tv/",
		MoviePath:      "/data/movies/",
		Default:        true,
	})
}

func (c *Client) GetSetting(key string) string {
	v, err := c.ent.Settings.Query().Where(settings.Key(key)).Only(context.TODO())
	if err != nil {
		log.Debugf("get setting by key: %s error: %v", key, err)
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
		return LanguageCN
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
		SetImdbID(m.ImdbID).
		SetStorageID(m.StorageID).
		SetOverview(m.Overview).
		SetNameCn(m.NameCn).
		SetNameEn(m.NameEn).
		SetOriginalName(m.OriginalName).
		SetMediaType(m.MediaType).
		SetAirDate(m.AirDate).
		SetResolution(m.Resolution).
		SetTargetDir(m.TargetDir).
		SetDownloadHistoryEpisodes(m.DownloadHistoryEpisodes).
		SetLimiter(m.Limiter).
		SetExtras(m.Extras).
		AddEpisodeIDs(episodes...).
		Save(context.TODO())
	return r, err

}

func (c *Client) GetMediaWatchlist(mediaType media.MediaType) []*ent.Media {
	list, err := c.ent.Media.Query().Where(media.MediaTypeEQ(mediaType)).Order(ent.Desc(media.FieldID)).All(context.TODO())
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
func (c *Client) GetEpisodeByID(epID int) (*ent.Episode, error) {
	return c.ent.Episode.Query().Where(episode.ID(epID)).First(context.TODO())
}

func (c *Client) UpdateEpiode(episodeId int, name, overview string) error {
	return c.ent.Episode.Update().Where(episode.ID(episodeId)).SetTitle(name).SetOverview(overview).Exec(context.TODO())
}

func (c *Client) UpdateEpiode2(episodeId int, name, overview, airdate string) error {
	return c.ent.Episode.Update().Where(episode.ID(episodeId)).SetTitle(name).SetOverview(overview).SetAirDate(airdate).Exec(context.TODO())
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

	ep, err := se.QueryEpisodes().All(context.Background())
	if err != nil {
		log.Errorf("get episodes %d: %v", id, err)
		return nil
	}
	md.Episodes = ep

	return md
}

func (c *Client) GetMedia(id int) (*ent.Media, error) {
	return c.ent.Media.Query().Where(media.ID(id)).First(context.TODO())
}

func (c *Client) DeleteMedia(id int) error {
	_, err := c.ent.Episode.Delete().Where(episode.MediaID(id)).Exec(context.TODO())
	if err != nil {
		return err
	}
	_, err = c.ent.Media.Delete().Where(media.ID(id)).Exec(context.TODO())
	if err != nil {
		return err
	}
	return c.CleanAllDanglingEpisodes()
}

func (c *Client) SaveEposideDetail(d *ent.Episode) (int, error) {
	ep, err := c.ent.Episode.Create().
		SetAirDate(d.AirDate).
		SetSeasonNumber(d.SeasonNumber).
		SetEpisodeNumber(d.EpisodeNumber).
		SetOverview(d.Overview).
		SetMonitored(d.Monitored).
		SetTitle(d.Title).Save(context.TODO())
	if err != nil {
		return 0, errors.Wrap(err, "save episode")
	}
	return ep.ID, nil
}

func (c *Client) SaveEposideDetail2(d *ent.Episode) (int, error) {
	ep, err := c.ent.Episode.Create().
		SetAirDate(d.AirDate).
		SetSeasonNumber(d.SeasonNumber).
		SetEpisodeNumber(d.EpisodeNumber).
		SetMediaID(d.MediaID).
		SetStatus(d.Status).
		SetOverview(d.Overview).
		SetTitle(d.Title).Save(context.TODO())

	return ep.ID, err
}

type TorznabSetting struct {
	URL    string `json:"url"`
	ApiKey string `json:"api_key"`
}

func (c *Client) SaveIndexer(in *ent.Indexers) error {

	if in.ID != 0 {
		//update setting
		return c.ent.Indexers.Update().Where(indexers.ID(in.ID)).SetName(in.Name).SetImplementation(in.Implementation).
			SetPriority(in.Priority).SetSettings(in.Settings).SetSeedRatio(in.SeedRatio).SetDisabled(in.Disabled).Exec(context.Background())
	}
	//create new one
	count := c.ent.Indexers.Query().Where(indexers.Name(in.Name)).CountX(context.TODO())
	if count > 0 {
		return fmt.Errorf("name already esxits: %v", in.Name)
	}

	_, err := c.ent.Indexers.Create().
		SetName(in.Name).SetImplementation(in.Implementation).SetPriority(in.Priority).SetSettings(in.Settings).SetSeedRatio(in.SeedRatio).
		SetDisabled(in.Disabled).Save(context.TODO())
	if err != nil {
		return errors.Wrap(err, "save db")
	}

	return nil
}

func (c *Client) DeleteTorznab(id int) {
	c.ent.Indexers.Delete().Where(indexers.ID(id)).Exec(context.TODO())
}

func (c *Client) GetIndexer(id int) (*TorznabInfo, error) {
	res, err := c.ent.Indexers.Query().Where(indexers.ID(id)).First(context.TODO())
	if err != nil {
		return nil, err
	}
	var ss TorznabSetting
	err = json.Unmarshal([]byte(res.Settings), &ss)
	if err != nil {

		return nil, fmt.Errorf("unmarshal torznab %s error: %v", res.Name, err)
	}
	return &TorznabInfo{Indexers: res, TorznabSetting: ss}, nil
}

type TorznabInfo struct {
	*ent.Indexers
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
			Indexers:       r,
			TorznabSetting: ss,
		})
	}
	return l
}

func (c *Client) SaveDownloader(downloader *ent.DownloadClients) error {
	count := c.ent.DownloadClients.Query().Where(downloadclients.Name(downloader.Name)).CountX(context.TODO())
	if count != 0 {
		err := c.ent.DownloadClients.Update().Where(downloadclients.Name(downloader.Name)).SetImplementation(downloader.Implementation).
			SetURL(downloader.URL).SetUser(downloader.User).SetPassword(downloader.Password).SetPriority1(downloader.Priority1).Exec(context.TODO())
		return err
	}

	_, err := c.ent.DownloadClients.Create().SetEnable(true).SetImplementation(downloader.Implementation).
		SetName(downloader.Name).SetURL(downloader.URL).SetUser(downloader.User).SetPriority1(downloader.Priority1).SetPassword(downloader.Password).Save(context.TODO())
	return err
}

func (c *Client) GetAllDonloadClients() []*ent.DownloadClients {
	cc, err := c.ent.DownloadClients.Query().Order(ent.Asc(downloadclients.FieldPriority1)).All(context.TODO())
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
	Name           string            `json:"name" binding:"required"`
	Implementation string            `json:"implementation" binding:"required"`
	Settings       map[string]string `json:"settings" binding:"required"`
	TvPath         string            `json:"tv_path" binding:"required"`
	MoviePath      string            `json:"movie_path" binding:"required"`
	Default        bool              `json:"default"`
}

func (s *StorageInfo) ToWebDavSetting() WebdavSetting {
	if s.Implementation != storage.ImplementationWebdav.String() {
		panic("not webdav storage")
	}
	return WebdavSetting{
		URL:            s.Settings["url"],
		User:           s.Settings["user"],
		Password:       s.Settings["password"],
		ChangeFileHash: s.Settings["change_file_hash"],
	}
}

type WebdavSetting struct {
	URL            string `json:"url"`
	User           string `json:"user"`
	Password       string `json:"password"`
	ChangeFileHash string `json:"change_file_hash"`
}

func (c *Client) AddStorage(st *StorageInfo) error {
	if !strings.HasSuffix(st.TvPath, "/") {
		st.TvPath += "/"
	}
	if !strings.HasSuffix(st.MoviePath, "/") {
		st.MoviePath += "/"
	}
	if st.Settings == nil {
		st.Settings = map[string]string{}
	}

	data, err := json.Marshal(st.Settings)
	if err != nil {
		return errors.Wrap(err, "json marshal")
	}

	count := c.ent.Storage.Query().Where(storage.Name(st.Name)).CountX(context.TODO())
	if count > 0 {
		//storage already exist, edit exist one
		return c.ent.Storage.Update().Where(storage.Name(st.Name)).
			SetImplementation(storage.Implementation(st.Implementation)).SetTvPath(st.TvPath).SetMoviePath(st.MoviePath).
			SetSettings(string(data)).Exec(context.TODO())
	}
	countAll := c.ent.Storage.Query().Where(storage.Deleted(false)).CountX(context.TODO())
	if countAll == 0 {
		log.Infof("first storage, make it default: %s", st.Name)
		st.Default = true
	}
	_, err = c.ent.Storage.Create().SetName(st.Name).
		SetImplementation(storage.Implementation(st.Implementation)).SetTvPath(st.TvPath).SetMoviePath(st.MoviePath).
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

func (s *Storage) ToWebDavSetting() WebdavSetting {
	if s.Implementation != storage.ImplementationWebdav {
		panic("not webdav storage")
	}
	var webdavSetting WebdavSetting
	json.Unmarshal([]byte(s.Settings), &webdavSetting)
	return webdavSetting
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
	if h.Link != "" {
		r, err := utils.Link2Magnet(h.Link)
		if err != nil {
			log.Warnf("convert link to magnet error, link %v, error: %v", h.Link, err)
		} else {
			h.Link = r
		}
	}
	return c.ent.History.Create().SetMediaID(h.MediaID).SetEpisodeID(h.EpisodeID).SetDate(time.Now()).
		SetStatus(h.Status).SetTargetDir(h.TargetDir).SetSourceTitle(h.SourceTitle).SetIndexerID(h.IndexerID).
		SetDownloadClientID(h.DownloadClientID).SetSize(h.Size).SetSaved(h.Saved).SetLink(h.Link).Save(context.TODO())
}

func (c *Client) SetHistoryStatus(id int, status history.Status) error {
	return c.ent.History.Update().Where(history.ID(id)).SetStatus(status).Exec(context.TODO())
}

func (c *Client) GetHistories() ent.Histories {
	h, err := c.ent.History.Query().Order(history.ByID(sql.OrderDesc())).All(context.TODO())
	if err != nil {
		return nil
	}
	return h
}

func (c *Client) GetRunningHistories() ent.Histories {
	h, err := c.ent.History.Query().Where(history.Or(history.StatusEQ(history.StatusRunning),
		history.StatusEQ(history.StatusUploading), history.StatusEQ(history.StatusSeeding))).All(context.TODO())
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

func (c *Client) UpdateEpisodeStatus(mediaID int, seasonNum, episodeNum int) error {
	ep, err := c.ent.Episode.Query().Where(episode.MediaID(mediaID)).Where(episode.EpisodeNumber(episodeNum)).
		Where(episode.SeasonNumber(seasonNum)).First(context.TODO())
	if err != nil {
		return errors.Wrap(err, "finding episode")
	}
	return ep.Update().SetStatus(episode.StatusDownloaded).Exec(context.TODO())
}

func (c *Client) SetEpisodeStatus(id int, status episode.Status) error {
	return c.ent.Episode.Update().Where(episode.ID(id)).SetStatus(status).Exec(context.TODO())
}

func (c *Client) IsEpisodeDownloadingOrDownloaded(id int) bool {
	his := c.ent.History.Query().Where(history.EpisodeID(id)).AllX(context.Background())
	for _, h := range his {
		if h.Status != history.StatusFail {
			return true
		}
	}
	return false
}

func (c *Client) SetSeasonAllEpisodeStatus(mediaID, seasonNum int, status episode.Status) error {
	return c.ent.Episode.Update().Where(episode.MediaID(mediaID), episode.SeasonNumber(seasonNum)).SetStatus(status).Exec(context.TODO())
}

func (c *Client) TmdbIdInWatchlist(tmdb_id int) bool {
	return c.ent.Media.Query().Where(media.TmdbID(tmdb_id)).CountX(context.TODO()) > 0
}

func (c *Client) GetDownloadHistory(mediaID int) ([]*ent.History, error) {
	return c.ent.History.Query().Where(history.MediaID(mediaID)).All(context.TODO())
}

func (c *Client) GetMovieDummyEpisode(movieId int) (*ent.Episode, error) {
	_, err := c.ent.Media.Query().Where(media.ID(movieId), media.MediaTypeEQ(media.MediaTypeMovie)).First(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "get movie")
	}
	ep, err := c.ent.Episode.Query().Where(episode.MediaID(movieId)).First(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "query episode")
	}
	return ep, nil
}

func (c *Client) GetDownloadClient(id int) (*ent.DownloadClients, error) {
	return c.ent.DownloadClients.Query().Where(downloadclients.ID(id)).First(context.Background())
}

func (c *Client) SetEpisodeMonitoring(id int, b bool) error {
	return c.ent.Episode.Update().Where(episode.ID(id)).SetMonitored(b).Exec(context.Background())
}

type EditMediaData struct {
	ID         int                 `json:"id"`
	Resolution media.Resolution    `json:"resolution"`
	TargetDir  string              `json:"target_dir"`
	Limiter    schema.MediaLimiter `json:"limiter"`
}

func (c *Client) EditMediaMetadata(in EditMediaData) error {
	return c.ent.Media.Update().Where(media.ID(in.ID)).SetResolution(in.Resolution).SetTargetDir(in.TargetDir).SetLimiter(in.Limiter).
		Exec(context.Background())
}

func (c *Client) UpdateEpisodeTargetFile(id int, filename string) error {
	return c.ent.Episode.Update().Where(episode.ID(id)).SetTargetFile(filename).Exec(context.Background())
}

func (c *Client) GetSeasonEpisodes(mediaId, seasonNum int) ([]*ent.Episode, error) {
	return c.ent.Episode.Query().Where(episode.MediaID(mediaId), episode.SeasonNumber(seasonNum)).All(context.Background())
}

func (c *Client) GetAllImportLists() ([]*ent.ImportList, error) {
	return c.ent.ImportList.Query().All(context.Background())
}

func (c *Client) AddImportlist(il *ent.ImportList) error {
	count, err := c.ent.ImportList.Query().Where(importlist.Name(il.Name)).Count(context.Background())
	if err != nil {
		return err
	}
	if count > 0 {
		//edit exist record
		return c.ent.ImportList.Update().Where(importlist.Name(il.Name)).
			SetURL(il.URL).SetQulity(il.Qulity).SetType(il.Type).SetStorageID(il.StorageID).Exec(context.Background())
	}
	return c.ent.ImportList.Create().SetName(il.Name).SetURL(il.URL).SetQulity(il.Qulity).SetStorageID(il.StorageID).
		SetType(il.Type).Exec(context.Background())
}

func (c *Client) DeleteImportlist(id int) error {
	return c.ent.ImportList.DeleteOneID(id).Exec(context.TODO())
}

func (c *Client) GetSizeLimiter() (*SizeLimiter, error) {
	v := c.GetSetting(SetttingSizeLimiter)
	var limiter SizeLimiter
	err := json.Unmarshal([]byte(v), &limiter)
	return &limiter, err
}

func (c *Client) SetSizeLimiter(limiter *SizeLimiter) error {
	data, err := json.Marshal(limiter)
	if err != nil {
		return err
	}
	return c.SetSetting(SetttingSizeLimiter, string(data))
}

func (c *Client) GetTvNamingFormat() string {
	s := c.GetSetting(SettingTvNamingFormat)
	if s == "" {
		return DefaultNamingFormat
	}
	return s
}

func (c *Client) GetMovingNamingFormat() string {
	s := c.GetSetting(SettingMovieNamingFormat)
	if s == "" {
		return DefaultNamingFormat
	}
	return s
}

func (c *Client) CleanAllDanglingEpisodes() error {
	_, err := c.ent.Episode.Delete().Where(episode.Not(episode.HasMedia())).Exec(context.Background())
	return err
}

func (c *Client) AddBlacklistItem(item *ent.Blacklist) error {
	return c.ent.Blacklist.Create().SetType(item.Type).SetValue(item.Value).SetNotes(item.Notes).Exec(context.Background())
}
