package db

import (
	"polaris/ent"
	"polaris/ent/episode"
	"polaris/ent/history"
	"polaris/ent/media"
)

type Database interface {
	GetSetting(key string) string
	SetSetting(key, value string) error
	GetLanguage() string
	GetDownloadDir() string
	GetTvNamingFormat() string
	GetMovingNamingFormat() string
	GetProwlarrSetting() (*ProwlarrSetting, error)
	SaveProwlarrSetting(se *ProwlarrSetting) error
	GetAcceptedVideoFormats() ([]string, error)
	SetAcceptedVideoFormats(key string, v []string) error
	GetAcceptedSubtitleFormats() ([]string, error)
	SetAcceptedSubtitleFormats(key string, v []string) error
	GetTmdbApiKey() string

	AddMediaWatchlist(m *ent.Media, episodes []int) (*ent.Media, error)
	GetMediaWatchlist(mediaType media.MediaType) []*ent.Media
	GetMediaDetails(id int) (*MediaDetails, error)
	GetMedia(id int) (*ent.Media, error)
	DeleteMedia(id int) error
	TmdbIdInWatchlist(tmdb_id int) bool
	EditMediaMetadata(in EditMediaData) error
	GetSizeLimiter(mediaType string) (*MediaSizeLimiter, error)
	SetSizeLimiter(mediaType string, limiter *MediaSizeLimiter) error

	GetEpisode(seriesId, seasonNum, episodeNum int) (*ent.Episode, error)
	GetEpisodeByID(epID int) (*ent.Episode, error)
	UpdateEpiode(episodeId int, name, overview string) error
	UpdateEpiode2(episodeId int, name, overview, airdate string) error
	SaveEposideDetail(d *ent.Episode) (int, error)
	SaveEposideDetail2(d *ent.Episode) (int, error)
	UpdateEpisodeStatus(mediaID int, seasonNum, episodeNum int) error
	SetEpisodeStatus(id int, status episode.Status) error
	IsEpisodeDownloadingOrDownloaded(id int) bool
	SetSeasonAllEpisodeStatus(mediaID, seasonNum int, status episode.Status) error
	SetEpisodeMonitoring(id int, b bool) error
	UpdateEpisodeTargetFile(id int, filename string) error
	GetSeasonEpisodes(mediaId, seasonNum int) ([]*ent.Episode, error)
	CleanAllDanglingEpisodes() error

	SaveIndexer(in *ent.Indexers) error
	DeleteIndexer(id int)
	GetIndexer(id int) (*ent.Indexers, error)
	GetAllIndexers() []*ent.Indexers

	SaveDownloader(downloader *ent.DownloadClients) error
	GetAllDonloadClients() []*ent.DownloadClients
	DeleteDownloadCLient(id int)
	GetDownloadClient(id int) (*ent.DownloadClients, error)

	AddStorage(st *StorageInfo) error
	GetAllStorage() []*ent.Storage
	GetStorage(id int) *Storage
	DeleteStorage(id int) error
	SetDefaultStorage(id int) error
	SetDefaultStorageByName(name string) error

	SaveHistoryRecord(h ent.History) (*ent.History, error)
	SetHistoryStatus(id int, status history.Status) error
	GetRunningHistories() ent.Histories
	GetHistory(id int) *ent.History
	GetHistories() ent.Histories
	DeleteHistory(id int) error
	GetDownloadHistory(mediaID int) ([]*ent.History, error)
	GetMovieDummyEpisode(movieId int) (*ent.Episode, error)

	GetAllImportLists() ([]*ent.ImportList, error)
	AddImportlist(il *ent.ImportList) error
	DeleteImportlist(id int) error

	GetAllNotificationClients2() ([]*ent.NotificationClient, error)
	GetAllNotificationClients() ([]*NotificationClient, error)
	AddNotificationClient(name, service string, setting string, enabled bool) error
	DeleteNotificationClient(id int) error
	GetNotificationClient(id int) (*NotificationClient, error)
}
