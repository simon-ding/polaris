package core

import (
	"fmt"
	"path/filepath"
	"polaris/db"
	"polaris/ent"
	"polaris/ent/episode"
	"polaris/ent/history"
	"polaris/ent/media"
	"polaris/log"
	"polaris/pkg"
	"polaris/pkg/notifier/message"
	"polaris/pkg/utils"

	"github.com/pkg/errors"
)

func (c *Client) addSysCron() {
	c.mustAddCron("@every 1m", c.checkTasks)
	c.mustAddCron("0 0 * * * *", func() {
		c.downloadTvSeries()
		c.downloadMovie()
	})
	c.mustAddCron("0 0 */12 * * *", c.checkAllSeriesNewSeason)
	c.cron.Start()
}

func (c *Client) mustAddCron(spec string, cmd func()) {
	if err := c.cron.AddFunc(spec, cmd); err != nil {
		log.Errorf("add func error: %v", err)
		panic(err)
	}
}

func (c *Client) checkTasks() {
	log.Debug("begin check tasks...")
	for id, t := range c.tasks {
		if !t.Exists() {
			log.Infof("task no longer exists: %v", id)
			delete(c.tasks, id)
			continue
		}
		log.Infof("task (%s) percentage done: %d%%", t.Name(), t.Progress())
		if t.Progress() == 100 {
			r := c.db.GetHistory(id)
			if r.Status == history.StatusSuccess {
				//task already success, check seed ratio
				torrent := c.tasks[id]
				ok := c.isSeedRatioLimitReached(r.IndexerID, torrent)
				if ok {
					log.Infof("torrent file seed ratio reached, remove: %v", torrent.Name())
					torrent.Remove()
					delete(c.tasks, id)
				} else {
					log.Infof("torrent file still sedding: %v", torrent.Name())
				}
				continue
			}
			log.Infof("task is done: %v", t.Name())
			c.sendMsg(fmt.Sprintf(message.DownloadComplete, t.Name()))
			go func() {
				if err := c.moveCompletedTask(id); err != nil {
					log.Infof("post tasks for id %v fail: %v", id, err)
				}
			}()
		}
	}
}

func (c *Client) moveCompletedTask(id int) (err1 error) {
	torrent := c.tasks[id]
	r := c.db.GetHistory(id)
	if r.Status == history.StatusUploading {
		log.Infof("task %d is already uploading, skip", id)
		return nil
	}
	c.db.SetHistoryStatus(r.ID, history.StatusUploading)
	seasonNum, err := utils.SeasonId(r.TargetDir)
	if err != nil {
		log.Errorf("no season id: %v", r.TargetDir)
		seasonNum = -1
	}
	downloadclient, err := c.db.GetDownloadClient(r.DownloadClientID)
	if err != nil {
		log.Errorf("get task download client error: %v, use default one", err)
		downloadclient = &ent.DownloadClients{RemoveCompletedDownloads: true, RemoveFailedDownloads: true}
	}
	torrentName := torrent.Name()

	defer func() {

		if err1 != nil {
			c.db.SetHistoryStatus(r.ID, history.StatusFail)
			if r.EpisodeID != 0 {
				c.db.SetEpisodeStatus(r.EpisodeID, episode.StatusMissing)
			} else {
				c.db.SetSeasonAllEpisodeStatus(r.MediaID, seasonNum, episode.StatusMissing)
			}
			c.sendMsg(fmt.Sprintf(message.ProcessingFailed, err))
			if downloadclient.RemoveFailedDownloads {
				log.Debugf("task failed, remove failed torrent and files related")
				delete(c.tasks, r.ID)
				torrent.Remove()
			}
		}
	}()

	series := c.db.GetMediaDetails(r.MediaID)
	if series == nil {
		return nil
	}
	st := c.db.GetStorage(series.StorageID)
	log.Infof("move task files to target dir: %v", r.TargetDir)
	stImpl, err := c.getStorage(st.ID, series.MediaType)
	if err != nil {
		return err
	}

	// .plexmatch file
	if err := c.writePlexmatch(r.MediaID, r.EpisodeID, r.TargetDir, torrentName); err != nil {
		log.Errorf("create .plexmatch file error: %v", err)
	}

	//如果种子是路径，则会把路径展开，只移动文件，类似 move dir/* dir2/, 如果种子是文件，则会直接移动文件，类似 move file dir/
	if err := stImpl.Copy(filepath.Join(c.db.GetDownloadDir(), torrentName), r.TargetDir); err != nil {
		return errors.Wrap(err, "move file")
	}

	c.db.SetHistoryStatus(r.ID, history.StatusSuccess)
	if r.EpisodeID != 0 {
		c.db.SetEpisodeStatus(r.EpisodeID, episode.StatusDownloaded)
	} else {
		c.db.SetSeasonAllEpisodeStatus(r.MediaID, seasonNum, episode.StatusDownloaded)
	}
	c.sendMsg(fmt.Sprintf(message.ProcessingComplete, torrentName))

	//判断是否需要删除本地文件
	ok := c.isSeedRatioLimitReached(r.IndexerID, torrent)
	if downloadclient.RemoveCompletedDownloads && ok {
		log.Debugf("download complete,remove torrent and files related")
		delete(c.tasks, r.ID)
		torrent.Remove()
	}

	log.Infof("move downloaded files to target dir success, file: %v, target dir: %v", torrentName, r.TargetDir)
	return nil
}

func (c *Client) CheckDownloadedSeriesFiles(m *ent.Media) error {
	if m.MediaType != media.MediaTypeTv {
		return nil
	}
	log.Infof("check files in directory: %s", m.TargetDir)

	var storageImpl, err = c.getStorage(m.StorageID, media.MediaTypeTv)
	if err != nil {
		return err
	}

	files, err := storageImpl.ReadDir(m.TargetDir)
	if err != nil {
		return errors.Wrapf(err, "read dir %s", m.TargetDir)
	}

	for _, in := range files {
		if !in.IsDir() { //season dir, ignore file
			continue
		}
		dir := filepath.Join(m.TargetDir, in.Name())
		epFiles, err := storageImpl.ReadDir(dir)
		if err != nil {
			log.Errorf("read dir %s error: %v", dir, err)
			continue
		}
		for _, ep := range epFiles {
			log.Infof("found file: %v", ep.Name())
			seNum, epNum, err := utils.FindSeasonEpisodeNum(ep.Name())
			if err != nil {
				log.Errorf("find season episode num error: %v", err)
				continue
			}
			log.Infof("found match, season num %d, episode num %d", seNum, epNum)
			ep, err := c.db.GetEpisode(m.ID, seNum, epNum)
			if err != nil {
				log.Error("update episode: %v", err)
				continue
			}
			err = c.db.SetEpisodeStatus(ep.ID, episode.StatusDownloaded)
			if err != nil {
				log.Error("update episode: %v", err)
				continue
			}

		}
	}
	return nil

}

type Task struct {
	//Processing bool
	pkg.Torrent
}

func (c *Client) downloadTvSeries() {
	log.Infof("begin check all tv series resources")
	allSeries := c.db.GetMediaWatchlist(media.MediaTypeTv)
	for _, series := range allSeries {
		tvDetail := c.db.GetMediaDetails(series.ID)
		m := make(map[int][]*ent.Episode)
		for _, ep := range tvDetail.Episodes {
			m[ep.SeasonNumber] = append(m[ep.SeasonNumber], ep)
		}
		for seasonNum, epsides := range m {
			wantedSeasonPack := true
			for _, ep := range epsides {
				if !ep.Monitored {
					wantedSeasonPack = false
				}
				if ep.Status != episode.StatusMissing {
					wantedSeasonPack = false
				}
			}
			if wantedSeasonPack {
				name, err := c.SearchAndDownload(series.ID, seasonNum, -1)
				if err != nil {
					log.Infof("cannot find resource to download : %v", err)
				} else {
					log.Infof("begin download torrent resource: %v", name)
				}

			} else {
				for _, ep := range epsides {
					if !ep.Monitored {
						continue
					}
					if ep.Status != episode.StatusMissing {
						continue
					}
					name, err := c.SearchAndDownload(series.ID, ep.SeasonNumber, ep.EpisodeNumber)
					if err != nil {
						log.Infof("cannot find resource to download for %s: %v", ep.Title, err)
					} else {
						log.Infof("begin download torrent resource: %v", name)
					}
				}

			}
		}
	}
}

func (c *Client) downloadMovie() {
	log.Infof("begin check all movie resources")
	allSeries := c.db.GetMediaWatchlist(media.MediaTypeMovie)

	for _, series := range allSeries {
		detail := c.db.GetMediaDetails(series.ID)
		if len(detail.Episodes) == 0 {
			log.Errorf("no related dummy episode: %v", detail.NameEn)
			continue
		}
		ep := detail.Episodes[0]
		if ep.Status != episode.StatusMissing {
			continue
		}

		if err := c.downloadMovieSingleEpisode(ep, series.TargetDir); err != nil {
			log.Errorf("download movie error: %v", err)
		}
	}
}

func (c *Client) downloadMovieSingleEpisode(ep *ent.Episode, targetDir string) error {
	trc, dlc, err := c.getDownloadClient()
	if err != nil {
		return errors.Wrap(err, "connect transmission")
	}
	qiangban := c.db.GetSetting(db.SettingAllowQiangban)
	allowQiangban := false
	if qiangban == "true" {
		allowQiangban = false
	}

	res, err := SearchMovie(c.db, &SearchParam{
		MediaId:         ep.MediaID,
		CheckFileSize:   true,
		CheckResolution: true,
		FilterQiangban:  !allowQiangban,
	})
	if err != nil {

		return errors.Wrap(err, "search movie")
	}
	r1 := res[0]
	log.Infof("begin download torrent resource: %v", r1.Name)
	torrent, err := trc.Download(r1.Link, c.db.GetDownloadDir())
	if err != nil {
		return errors.Wrap(err, "downloading")
	}
	torrent.Start()

	history, err := c.db.SaveHistoryRecord(ent.History{
		MediaID:          ep.MediaID,
		EpisodeID:        ep.ID,
		SourceTitle:      r1.Name,
		TargetDir:        targetDir,
		Status:           history.StatusRunning,
		Size:             r1.Size,
		Saved:            torrent.Save(),
		DownloadClientID: dlc.ID,
		IndexerID:        r1.IndexerId,
	})
	if err != nil {
		log.Errorf("save history error: %v", err)
	}

	c.tasks[history.ID] = &Task{Torrent: torrent}

	c.db.SetEpisodeStatus(ep.ID, episode.StatusDownloading)
	return nil
}

func (c *Client) checkAllSeriesNewSeason() {
	log.Infof("begin checking series all new season")
	allSeries := c.db.GetMediaWatchlist(media.MediaTypeTv)
	for _, series := range allSeries {
		err := c.checkSeiesNewSeason(series)
		if err != nil {
			log.Errorf("check series new season error: series name %v, error: %v", series.NameEn, err)
		}
	}
}

func (c *Client) checkSeiesNewSeason(media *ent.Media) error {
	d, err := c.MustTMDB().GetTvDetails(media.TmdbID, c.language)
	if err != nil {
		return errors.Wrap(err, "tmdb")
	}
	lastsSason := d.NumberOfSeasons
	seasonDetail, err := c.MustTMDB().GetSeasonDetails(media.TmdbID, lastsSason, c.language)
	if err != nil {
		return errors.Wrap(err, "tmdb season")
	}

	for _, ep := range seasonDetail.Episodes {
		epDb, err := c.db.GetEpisode(media.ID, ep.SeasonNumber, ep.EpisodeNumber)
		if err != nil {
			if ent.IsNotFound(err) {
				log.Infof("add new episode: %+v", ep)
				episode := &ent.Episode{
					MediaID:       media.ID,
					SeasonNumber:  ep.SeasonNumber,
					EpisodeNumber: ep.EpisodeNumber,
					Title:         ep.Name,
					Overview:      ep.Overview,
					AirDate:       ep.AirDate,
					Status:        episode.StatusMissing,
					Monitored:     true,
				}
				c.db.SaveEposideDetail2(episode)
			}
		} else { //update episode
			if ep.Name != epDb.Title || ep.Overview != epDb.Overview || ep.AirDate != epDb.AirDate {
				log.Infof("update new episode: %+v", ep)
				c.db.UpdateEpiode2(epDb.ID, ep.Name, ep.Overview, ep.AirDate)
			}
		}
	}
	return nil
}

func (c *Client) isSeedRatioLimitReached(indexId int, t pkg.Torrent) bool {
	indexer, err := c.db.GetIndexer(indexId)
	if err != nil {
		return true
	}
	currentRatio := t.SeedRatio()
	if currentRatio == nil {
		log.Warnf("get current seed ratio error, current ratio is nil")
		return indexer.SeedRatio == 0
	}
	return *currentRatio >= float64(indexer.SeedRatio)
}
