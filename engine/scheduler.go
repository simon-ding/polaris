package engine

import (
	"fmt"
	"os"
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
	"time"

	"github.com/pkg/errors"
)

type DownloadOptions struct {
	HashFilterFn func(hash string) bool
	SeasonNum    int
	EpisodeNums  []int
	MediaId      int
}

func (c *Engine) addSysCron() {
	c.registerCronJob("check_running_tasks", "@every 1m", c.checkTasks)
	c.registerCronJob("check_available_medias_to_download", "0 0 * * * *", func() error {
		v := os.Getenv("POLARIS_NO_AUTO_DOWNLOAD")
		if v == "true" {
			return nil
		}
		if err := c.syncProwlarr(); err != nil {
			log.Warnf("sync prowlarr error: %v", err)
		}
		c.downloadAllTvSeries()
		c.downloadAllMovies()
		return nil
	})
	c.registerCronJob("check_series_new_release", "0 0 */12 * * *", c.checkAllSeriesNewSeason)
	c.registerCronJob("update_import_lists", "0 30 * * * *", c.periodicallyUpdateImportlist)

	c.schedulers.Range(func(key string, value scheduler) bool {
		log.Debugf("add cron job: %v", key)
		c.mustAddCron(value.cron, func() {
			if err := value.f(); err != nil {
				log.Errorf("exexuting cron job %s error: %v", key, err)
			}
		})
		return true
	})
	c.cron.Start()
	log.Infof("--------- add cron jobs done --------")
}

func (c *Engine) mustAddCron(spec string, cmd func()) {
	if err := c.cron.AddFunc(spec, cmd); err != nil {
		log.Errorf("add func error: %v", err)
		panic(err)
	}
}

func (c *Engine) TriggerCronJob(name string) error {
	job, ok := c.schedulers.Load(name)
	if !ok {
		return fmt.Errorf("job name not exists: %s", name)
	}
	return job.f()
}

func (c *Engine) checkTasks() error {
	log.Debug("begin check tasks...")
	c.tasks.Range(func(id int, t *Task) bool {
		r := c.db.GetHistory(id)
		if !t.Exists() {
			log.Infof("task no longer exists: %v", id)

			c.tasks.Delete(id)
			return true
		}
		name, err := t.Name()
		if err != nil {
			log.Warnf("get task name error: %v", err)
			return true
		}

		progress, err := t.Progress()
		if err != nil {
			log.Warnf("get task progress error: %v", err)
			return true
		}
		log.Infof("task (%s) percentage done: %d%%", name, progress)
		if progress == 100 {

			if r.Status == history.StatusSeeding {
				//task already success, check seed ratio
				torrent, _ := c.tasks.Load(id)
				ratio, ok := c.isSeedRatioLimitReached(r.IndexerID, torrent)
				if ok {
					log.Infof("torrent file seed ratio reached, remove: %v, current seed ratio: %v", name, ratio)
					torrent.Remove()
					c.tasks.Delete(id)
					c.setHistoryStatus(id, history.StatusSuccess)
				} else {
					log.Infof("torrent file still sedding: %v, current seed ratio: %v", name, ratio)
				}
				return true
			} else if r.Status == history.StatusRunning {
				log.Infof("task is done: %v", name)
				c.sendMsg(fmt.Sprintf(message.DownloadComplete, name))
				go c.postTaskProcessing(id)
			}
		}

		return true
	})
	return nil
}

/*
episode 状态有3种：missing、downloading、downloaded

history状态有5种：running, success, fail, uploading, seeding

没有下载的剧集状态都是missing，已下载完成的都是downloaded，正在下载的是downloading

对应的history状态，下载任务创建成功，正常跑着是running，出了问题失败了，就是fail，下载完成的任务会先进入uploading状态进一步处理，
uploading状态下会传输到对应的存储里面，uploading成功如果需要做种会进入seeding状态，如果不做种进入success状态，失败了会进入fail状态

seeding状态中，会定时检查做种状态，达到指定分享率，会置为success

任务创建成功，episode状态会由missing置为downloading，如果任务失败重新置为missing，如果任务成功进入success或seeding，episode状态应置为downloaded

*/

func (c *Engine) setHistoryStatus(id int, status history.Status) {
	r := c.db.GetHistory(id)

	episodeIds := c.GetEpisodeIds(r)

	switch status {
	case history.StatusRunning:
		c.db.SetHistoryStatus(id, history.StatusRunning)
		c.setEpsideoStatus(episodeIds, episode.StatusDownloading)
	case history.StatusSuccess:
		c.db.SetHistoryStatus(id, history.StatusSuccess)
		c.setEpsideoStatus(episodeIds, episode.StatusDownloaded)

	case history.StatusUploading:
		c.db.SetHistoryStatus(id, history.StatusUploading)

	case history.StatusSeeding:
		c.db.SetHistoryStatus(id, history.StatusSeeding)
		c.setEpsideoStatus(episodeIds, episode.StatusDownloaded)

	case history.StatusFail:
		c.db.SetHistoryStatus(id, history.StatusFail)
		c.setEpsideoStatus(episodeIds, episode.StatusMissing)
	default:
		panic(fmt.Sprintf("unkown status %v", status))
	}
}

func (c *Engine) setEpsideoStatus(episodeIds []int, status episode.Status) error {
	for _, id := range episodeIds {
		ep, err := c.db.GetEpisodeByID(id)
		if err != nil {
			return err
		}
		if ep.Status == episode.StatusDownloaded {
			//已经下载完成的任务，不再重新设置状态
			continue
		}

		if err := c.db.SetEpisodeStatus(id, status); err != nil {
			return err
		}
	}

	return nil
}

func (c *Engine) postTaskProcessing(id int) {
	if err := c.findEpisodeFilesPreMoving(id); err != nil {
		log.Errorf("finding all episode file error: %v", err)
	} else {
		if err := c.writePlexmatch(id); err != nil {
			log.Errorf("write plexmatch file error: %v", err)
		}
		if err := c.writeNfoFile(id); err != nil {
			log.Errorf("write nfo file error: %v", err)
		}
	}
	if err := c.moveCompletedTask(id); err != nil {
		log.Infof("post tasks for id %v fail: %v", id, err)
	}
}

func getSeasonNum(h *ent.History) int {
	if h.SeasonNum != 0 {
		return h.SeasonNum
	}
	seasonNum, err := utils.SeasonId(h.TargetDir)
	if err != nil {
		log.Errorf("no season id: %v", h.TargetDir)
		seasonNum = -1
	}
	return seasonNum
}

func (c *Engine) GetEpisodeIds(r *ent.History) []int {

	series, err := c.db.GetMediaDetails(r.MediaID)
	if err != nil {
		log.Errorf("get media details error: %v", err)
		return []int{}
	}
	if series.MediaType == media.MediaTypeMovie { //movie
		ep, _ := c.db.GetMovieDummyEpisode(series.ID)
		return []int{ep.ID}
	} else { //tv
		var episodeIds []int
		seasonNum := getSeasonNum(r)

		if len(r.EpisodeNums) > 0 {
			for _, epNum := range r.EpisodeNums {
				for _, ep := range series.Episodes {
					if ep.SeasonNumber == seasonNum && ep.EpisodeNumber == epNum {
						episodeIds = append(episodeIds, ep.ID)
					}
				}
			}
		} else {
			for _, ep := range series.Episodes {
				if ep.SeasonNumber == seasonNum {
					episodeIds = append(episodeIds, ep.ID)
				}
			}

		}
		return episodeIds
	}
}

func (c *Engine) moveCompletedTask(id int) (err1 error) {
	torrent, _ := c.tasks.Load(id)
	r := c.db.GetHistory(id)
	// if r.Status == history.StatusUploading {
	// 	log.Infof("task %d is already uploading, skip", id)
	// 	return nil
	// }

	c.setHistoryStatus(r.ID, history.StatusUploading)

	downloadclient, err := c.db.GetDownloadClient(r.DownloadClientID)
	if err != nil {
		log.Errorf("get task download client error: %v, use default one", err)
		downloadclient = &ent.DownloadClients{RemoveCompletedDownloads: true, RemoveFailedDownloads: true}
	}
	torrentName, err := torrent.Name()
	if err != nil {
		return err
	}

	defer func() {

		if err1 != nil {
			c.setHistoryStatus(r.ID, history.StatusFail)
			c.sendMsg(fmt.Sprintf(message.ProcessingFailed, err1))
			if downloadclient.RemoveFailedDownloads {
				log.Debugf("task failed, remove failed torrent and files related")
				c.tasks.Delete(r.ID)
				torrent.Remove()
			}
		}
	}()

	series, err := c.db.GetMediaDetails(r.MediaID)
	if err != nil {
		return err
	}

	st := c.db.GetStorage(series.StorageID)
	log.Infof("move task files to target dir: %v", r.TargetDir)
	stImpl, err := c.GetStorage(st.ID, series.MediaType)
	if err != nil {
		return err
	}

	//如果种子是路径，则会把路径展开，只移动文件，类似 move dir/* dir2/, 如果种子是文件，则会直接移动文件，类似 move file dir/
	if err := stImpl.Copy(filepath.Join(c.db.GetDownloadDir(), torrentName), r.TargetDir, torrent.WalkFunc()); err != nil {
		return errors.Wrap(err, "move file")
	}
	torrent.UploadProgresser = stImpl.UploadProgress

	c.sendMsg(fmt.Sprintf(message.ProcessingComplete, torrentName))

	//判断是否需要删除本地文件, TODO prowlarr has no indexer id
	r1, ok := c.isSeedRatioLimitReached(r.IndexerID, torrent)
	if downloadclient.RemoveCompletedDownloads && ok {
		log.Debugf("download complete,remove torrent and files related, torrent: %v, seed ratio: %v", torrentName, r1)
		c.setHistoryStatus(r.ID, history.StatusSuccess)
		c.tasks.Delete(r.ID)
		torrent.Remove()
	} else {
		log.Infof("task complete but still needs seeding: %v", torrentName)
		c.setHistoryStatus(r.ID, history.StatusSeeding)
	}

	log.Infof("move downloaded files to target dir success, file: %v, target dir: %v", torrentName, r.TargetDir)
	return nil
}

func (c *Engine) CheckDownloadedSeriesFiles(m *ent.Media) error {
	if m.MediaType != media.MediaTypeTv {
		return nil
	}
	log.Infof("check files in directory: %s", m.TargetDir)

	var storageImpl, err = c.GetStorage(m.StorageID, media.MediaTypeTv)
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
	UploadProgresser func() float64
}

func (c *Engine) DownloadSeriesAllEpisodes(id int) []string {
	tvDetail, err := c.db.GetMediaDetails(id)
	if err != nil {
		log.Errorf("get media details error: %v", err)
		return nil
	}
	m := make(map[int][]*ent.Episode)
	for _, ep := range tvDetail.Episodes {
		m[ep.SeasonNumber] = append(m[ep.SeasonNumber], ep)
	}
	var allNames []string
	for seasonNum, epsides := range m {
		if seasonNum == 0 {
			continue
		}
		wantedSeasonPack := true
		seasonEpisodesWanted := make(map[int][]int, 0)
		for _, ep := range epsides {
			if !ep.Monitored {
				wantedSeasonPack = false
				continue
			}
			if ep.Status != episode.StatusMissing {
				wantedSeasonPack = false
				continue
			}
			if ep.AirDate != "" {
				t, err := time.Parse("2006-01-02", ep.AirDate)
				if err != nil {
					continue
				}
				/*
				 -------- now ------ t -----
				 t - 1day < now 要检测的剧集
				 提前一天开始检测
				*/
				if time.Now().Before(t.Add(-24 * time.Hour)) { //not aired
					wantedSeasonPack = false
					continue
				}
			}
			seasonEpisodesWanted[ep.SeasonNumber] = append(seasonEpisodesWanted[ep.SeasonNumber], ep.EpisodeNumber)
		}
		if wantedSeasonPack {
			names, err := c.SearchAndDownload(id, seasonNum)
			if err == nil {
				allNames = append(allNames, names...)
				log.Infof("begin download torrent resource: %v", names)
			} else {
				log.Warnf("finding season pack error: %v", err)
				wantedSeasonPack = false
			}
		}
		if !wantedSeasonPack {

			for se, eps := range seasonEpisodesWanted {
				names, err := c.SearchAndDownload(id, se, eps...)
				if err != nil {
					log.Warnf("finding resoruces of season %d episode %v error: %v", se, eps, err)
					continue
				} else {
					allNames = append(allNames, names...)
					log.Infof("begin download torrent resource: %v", names)
				}
			}

		}

	}
	return allNames
}

func (c *Engine) downloadAllTvSeries() {
	log.Infof("begin check all tv series resources")
	allSeries := c.db.GetMediaWatchlist(media.MediaTypeTv)
	for _, series := range allSeries {
		c.DownloadSeriesAllEpisodes(series.ID)
	}
}

func (c *Engine) downloadAllMovies() {
	log.Infof("begin check all movie resources")
	allSeries := c.db.GetMediaWatchlist(media.MediaTypeMovie)

	for _, series := range allSeries {
		if _, err := c.DownloadMovieByID(series.ID); err != nil {
			log.Errorf("download movie error: %v", err)
		}
	}
}

func (c *Engine) DownloadMovieByID(id int) (string, error) {
	detail, err := c.db.GetMediaDetails(id)
	if err != nil {
		return "", errors.Wrap(err, "get media details")
	}
	if len(detail.Episodes) == 0 {
		return "", fmt.Errorf("no related dummy episode: %v", detail.NameEn)
	}
	ep := detail.Episodes[0]
	if ep.Status != episode.StatusMissing {
		return "", nil
	}

	if name, err := c.downloadMovieSingleEpisode(detail.Media, ep); err != nil {
		return "", errors.Wrap(err, "download movie")
	} else {
		return name, nil
	}
}

func (c *Engine) downloadMovieSingleEpisode(m *ent.Media, ep *ent.Episode) (string, error) {

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
		return "", errors.Wrap(err, "search movie")
	}
	r1 := res[0]
	log.Infof("begin download torrent resource: %v", r1.Name)

	s, err := c.downloadTorrent(m, r1, DownloadOptions{MediaId: m.ID, SeasonNum: 0, HashFilterFn: c.hashInBlacklist})
	if err != nil {
		return "", err
	}
	return *s, nil
}

func (c *Engine) checkAllSeriesNewSeason() error {
	log.Infof("begin checking series all new season")
	allSeries := c.db.GetMediaWatchlist(media.MediaTypeTv)
	for _, series := range allSeries {
		err := c.checkSeiesNewSeason(series)
		if err != nil {
			log.Errorf("check series new season error: series name %v, error: %v", series.NameEn, err)
		}
	}
	return nil
}

func (c *Engine) checkSeiesNewSeason(media *ent.Media) error {
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

func (c *Engine) isSeedRatioLimitReached(indexId int, t pkg.Torrent) (float64, bool) {
	indexer, err := c.db.GetIndexer(indexId)
	if err != nil {
		return 0, true
	}
	currentRatio, err := t.SeedRatio()
	if err != nil {
		log.Warnf("get current seed ratio error: %v", err)
		return currentRatio, indexer.SeedRatio == 0
	}
	return currentRatio, currentRatio >= float64(indexer.SeedRatio)
}
