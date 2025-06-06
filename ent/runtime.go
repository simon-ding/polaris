// Code generated by ent, DO NOT EDIT.

package ent

import (
	"polaris/ent/blacklist"
	"polaris/ent/downloadclients"
	"polaris/ent/episode"
	"polaris/ent/history"
	"polaris/ent/indexers"
	"polaris/ent/media"
	"polaris/ent/notificationclient"
	"polaris/ent/schema"
	"polaris/ent/storage"
	"time"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	blacklistFields := schema.Blacklist{}.Fields()
	_ = blacklistFields
	// blacklistDescCreateTime is the schema descriptor for create_time field.
	blacklistDescCreateTime := blacklistFields[4].Descriptor()
	// blacklist.DefaultCreateTime holds the default value on creation for the create_time field.
	blacklist.DefaultCreateTime = blacklistDescCreateTime.Default.(func() time.Time)
	downloadclientsFields := schema.DownloadClients{}.Fields()
	_ = downloadclientsFields
	// downloadclientsDescUser is the schema descriptor for user field.
	downloadclientsDescUser := downloadclientsFields[4].Descriptor()
	// downloadclients.DefaultUser holds the default value on creation for the user field.
	downloadclients.DefaultUser = downloadclientsDescUser.Default.(string)
	// downloadclientsDescPassword is the schema descriptor for password field.
	downloadclientsDescPassword := downloadclientsFields[5].Descriptor()
	// downloadclients.DefaultPassword holds the default value on creation for the password field.
	downloadclients.DefaultPassword = downloadclientsDescPassword.Default.(string)
	// downloadclientsDescSettings is the schema descriptor for settings field.
	downloadclientsDescSettings := downloadclientsFields[6].Descriptor()
	// downloadclients.DefaultSettings holds the default value on creation for the settings field.
	downloadclients.DefaultSettings = downloadclientsDescSettings.Default.(string)
	// downloadclientsDescPriority1 is the schema descriptor for priority1 field.
	downloadclientsDescPriority1 := downloadclientsFields[7].Descriptor()
	// downloadclients.DefaultPriority1 holds the default value on creation for the priority1 field.
	downloadclients.DefaultPriority1 = downloadclientsDescPriority1.Default.(int)
	// downloadclients.Priority1Validator is a validator for the "priority1" field. It is called by the builders before save.
	downloadclients.Priority1Validator = downloadclientsDescPriority1.Validators[0].(func(int) error)
	// downloadclientsDescUseNatTraversal is the schema descriptor for use_nat_traversal field.
	downloadclientsDescUseNatTraversal := downloadclientsFields[8].Descriptor()
	// downloadclients.DefaultUseNatTraversal holds the default value on creation for the use_nat_traversal field.
	downloadclients.DefaultUseNatTraversal = downloadclientsDescUseNatTraversal.Default.(bool)
	// downloadclientsDescRemoveCompletedDownloads is the schema descriptor for remove_completed_downloads field.
	downloadclientsDescRemoveCompletedDownloads := downloadclientsFields[9].Descriptor()
	// downloadclients.DefaultRemoveCompletedDownloads holds the default value on creation for the remove_completed_downloads field.
	downloadclients.DefaultRemoveCompletedDownloads = downloadclientsDescRemoveCompletedDownloads.Default.(bool)
	// downloadclientsDescRemoveFailedDownloads is the schema descriptor for remove_failed_downloads field.
	downloadclientsDescRemoveFailedDownloads := downloadclientsFields[10].Descriptor()
	// downloadclients.DefaultRemoveFailedDownloads holds the default value on creation for the remove_failed_downloads field.
	downloadclients.DefaultRemoveFailedDownloads = downloadclientsDescRemoveFailedDownloads.Default.(bool)
	// downloadclientsDescTags is the schema descriptor for tags field.
	downloadclientsDescTags := downloadclientsFields[11].Descriptor()
	// downloadclients.DefaultTags holds the default value on creation for the tags field.
	downloadclients.DefaultTags = downloadclientsDescTags.Default.(string)
	// downloadclientsDescCreateTime is the schema descriptor for create_time field.
	downloadclientsDescCreateTime := downloadclientsFields[12].Descriptor()
	// downloadclients.DefaultCreateTime holds the default value on creation for the create_time field.
	downloadclients.DefaultCreateTime = downloadclientsDescCreateTime.Default.(func() time.Time)
	episodeFields := schema.Episode{}.Fields()
	_ = episodeFields
	// episodeDescMonitored is the schema descriptor for monitored field.
	episodeDescMonitored := episodeFields[7].Descriptor()
	// episode.DefaultMonitored holds the default value on creation for the monitored field.
	episode.DefaultMonitored = episodeDescMonitored.Default.(bool)
	// episodeDescCreateTime is the schema descriptor for create_time field.
	episodeDescCreateTime := episodeFields[9].Descriptor()
	// episode.DefaultCreateTime holds the default value on creation for the create_time field.
	episode.DefaultCreateTime = episodeDescCreateTime.Default.(func() time.Time)
	historyFields := schema.History{}.Fields()
	_ = historyFields
	// historyDescSize is the schema descriptor for size field.
	historyDescSize := historyFields[6].Descriptor()
	// history.DefaultSize holds the default value on creation for the size field.
	history.DefaultSize = historyDescSize.Default.(int)
	// historyDescCreateTime is the schema descriptor for create_time field.
	historyDescCreateTime := historyFields[12].Descriptor()
	// history.DefaultCreateTime holds the default value on creation for the create_time field.
	history.DefaultCreateTime = historyDescCreateTime.Default.(func() time.Time)
	indexersFields := schema.Indexers{}.Fields()
	_ = indexersFields
	// indexersDescSettings is the schema descriptor for settings field.
	indexersDescSettings := indexersFields[2].Descriptor()
	// indexers.DefaultSettings holds the default value on creation for the settings field.
	indexers.DefaultSettings = indexersDescSettings.Default.(string)
	// indexersDescEnableRss is the schema descriptor for enable_rss field.
	indexersDescEnableRss := indexersFields[3].Descriptor()
	// indexers.DefaultEnableRss holds the default value on creation for the enable_rss field.
	indexers.DefaultEnableRss = indexersDescEnableRss.Default.(bool)
	// indexersDescPriority is the schema descriptor for priority field.
	indexersDescPriority := indexersFields[4].Descriptor()
	// indexers.DefaultPriority holds the default value on creation for the priority field.
	indexers.DefaultPriority = indexersDescPriority.Default.(int)
	// indexers.PriorityValidator is a validator for the "priority" field. It is called by the builders before save.
	indexers.PriorityValidator = indexersDescPriority.Validators[0].(func(int) error)
	// indexersDescSeedRatio is the schema descriptor for seed_ratio field.
	indexersDescSeedRatio := indexersFields[5].Descriptor()
	// indexers.DefaultSeedRatio holds the default value on creation for the seed_ratio field.
	indexers.DefaultSeedRatio = indexersDescSeedRatio.Default.(float32)
	// indexersDescDisabled is the schema descriptor for disabled field.
	indexersDescDisabled := indexersFields[6].Descriptor()
	// indexers.DefaultDisabled holds the default value on creation for the disabled field.
	indexers.DefaultDisabled = indexersDescDisabled.Default.(bool)
	// indexersDescTvSearch is the schema descriptor for tv_search field.
	indexersDescTvSearch := indexersFields[7].Descriptor()
	// indexers.DefaultTvSearch holds the default value on creation for the tv_search field.
	indexers.DefaultTvSearch = indexersDescTvSearch.Default.(bool)
	// indexersDescMovieSearch is the schema descriptor for movie_search field.
	indexersDescMovieSearch := indexersFields[8].Descriptor()
	// indexers.DefaultMovieSearch holds the default value on creation for the movie_search field.
	indexers.DefaultMovieSearch = indexersDescMovieSearch.Default.(bool)
	// indexersDescSynced is the schema descriptor for synced field.
	indexersDescSynced := indexersFields[11].Descriptor()
	// indexers.DefaultSynced holds the default value on creation for the synced field.
	indexers.DefaultSynced = indexersDescSynced.Default.(bool)
	// indexersDescCreateTime is the schema descriptor for create_time field.
	indexersDescCreateTime := indexersFields[12].Descriptor()
	// indexers.DefaultCreateTime holds the default value on creation for the create_time field.
	indexers.DefaultCreateTime = indexersDescCreateTime.Default.(func() time.Time)
	mediaFields := schema.Media{}.Fields()
	_ = mediaFields
	// mediaDescCreatedAt is the schema descriptor for created_at field.
	mediaDescCreatedAt := mediaFields[7].Descriptor()
	// media.DefaultCreatedAt holds the default value on creation for the created_at field.
	media.DefaultCreatedAt = mediaDescCreatedAt.Default.(time.Time)
	// mediaDescAirDate is the schema descriptor for air_date field.
	mediaDescAirDate := mediaFields[8].Descriptor()
	// media.DefaultAirDate holds the default value on creation for the air_date field.
	media.DefaultAirDate = mediaDescAirDate.Default.(string)
	// mediaDescDownloadHistoryEpisodes is the schema descriptor for download_history_episodes field.
	mediaDescDownloadHistoryEpisodes := mediaFields[12].Descriptor()
	// media.DefaultDownloadHistoryEpisodes holds the default value on creation for the download_history_episodes field.
	media.DefaultDownloadHistoryEpisodes = mediaDescDownloadHistoryEpisodes.Default.(bool)
	// mediaDescCreateTime is the schema descriptor for create_time field.
	mediaDescCreateTime := mediaFields[16].Descriptor()
	// media.DefaultCreateTime holds the default value on creation for the create_time field.
	media.DefaultCreateTime = mediaDescCreateTime.Default.(func() time.Time)
	notificationclientFields := schema.NotificationClient{}.Fields()
	_ = notificationclientFields
	// notificationclientDescEnabled is the schema descriptor for enabled field.
	notificationclientDescEnabled := notificationclientFields[3].Descriptor()
	// notificationclient.DefaultEnabled holds the default value on creation for the enabled field.
	notificationclient.DefaultEnabled = notificationclientDescEnabled.Default.(bool)
	storageFields := schema.Storage{}.Fields()
	_ = storageFields
	// storageDescDeleted is the schema descriptor for deleted field.
	storageDescDeleted := storageFields[5].Descriptor()
	// storage.DefaultDeleted holds the default value on creation for the deleted field.
	storage.DefaultDeleted = storageDescDeleted.Default.(bool)
	// storageDescDefault is the schema descriptor for default field.
	storageDescDefault := storageFields[6].Descriptor()
	// storage.DefaultDefault holds the default value on creation for the default field.
	storage.DefaultDefault = storageDescDefault.Default.(bool)
	// storageDescCreateTime is the schema descriptor for create_time field.
	storageDescCreateTime := storageFields[7].Descriptor()
	// storage.DefaultCreateTime holds the default value on creation for the create_time field.
	storage.DefaultCreateTime = storageDescCreateTime.Default.(func() time.Time)
}
