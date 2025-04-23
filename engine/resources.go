package engine

import (
	"bytes"
	"fmt"
	"polaris/ent"
	"polaris/ent/episode"
	"polaris/ent/history"
	"polaris/ent/media"
	"polaris/log"
	"polaris/pkg/metadata"
	"polaris/pkg/notifier/message"
	"polaris/pkg/torznab"
	"polaris/pkg/utils"

	"github.com/pkg/errors"
)

func (c *Engine) DownloadEpisodeTorrent(r1 torznab.Result, op DownloadOptions) (*string, error) {

	series, err := c.db.GetMedia(op.MediaId)
	if err != nil {
		return nil, fmt.Errorf("no tv series of id %v", op.MediaId)
	}

	return c.downloadTorrent(series, r1, op)
}

/*
tmdb 校验获取的资源名，如果用资源名在tmdb搜索出来的结果能匹配上想要的资源，则认为资源有效，否则无效
解决名称过于简单的影视会匹配过多资源的问题， 例如：梦魇绝镇 FROM
*/
func (c *Engine) checkBtReourceWithTmdb(r *torznab.Result, seriesId int) bool {
	m := metadata.ParseTv(r.Name)
	se, err := c.MustTMDB().SearchMedia(m.NameEn, "", 1)
	if err != nil {
		log.Warnf("tmdb search error, consider this torrent ok: %v", err)
		return true
	} else {
		if len(se.Results) == 0 {
			log.Debugf("tmdb search no result, consider this torrent ok: %s", r.Name) //because tv name parse is not accurate
			return true
		}
		series, err := c.db.GetMediaDetails(seriesId)
		if err != nil {
			log.Warnf("get media details error: %v", err)
			return false
		}

		se0 := se.Results[0]
		if se0.ID != int64(series.TmdbID) {
			log.Warnf("bt reosurce name not match tmdb id: %s", r.Name)
			return false
		} else { //resource tmdb id match
			return true
		}
	}
}

func (c *Engine) SearchAndDownload(seriesId, seasonNum int, episodeNums ...int) ([]string, error) {

	res, err := SearchTvSeries(c.db, &SearchParam{
		MediaId:         seriesId,
		SeasonNum:       seasonNum,
		Episodes:        episodeNums,
		CheckFileSize:   true,
		CheckResolution: true,
	})
	if err != nil {
		return nil, err
	}
	wanted := make(map[int]bool, len(episodeNums))
	for _, ep := range episodeNums {
		wanted[ep] = true
	}
	var torrentNames []string
lo:
	for _, r := range res {
		if !c.checkBtReourceWithTmdb(&r, seriesId) {
			continue
		}
		m := metadata.ParseTv(r.Name)
		m.ParseExtraDescription(r.Description)
		if len(episodeNums) == 0 { //want season pack
			if m.IsSeasonPack {
				name, err := c.DownloadEpisodeTorrent(r, DownloadOptions{
					SeasonNum: seasonNum,
					MediaId:   seriesId,
					HashFilterFn: c.hashInBlacklist,
				})
				if err != nil {
					log.Warnf("download season pack error, continue next item: %v", err)
					continue lo
				}
				torrentNames = append(torrentNames, *name)
				break lo
			}
		} else {
			torrentEpisodes := make([]int, 0)
			for i := m.StartEpisode; i <= m.EndEpisode; i++ {
				if !wanted[i] { //torrent has episode not wanted
					continue lo
				}
				torrentEpisodes = append(torrentEpisodes, i)
			}
			name, err := c.DownloadEpisodeTorrent(r, DownloadOptions{
				SeasonNum: seasonNum,
				MediaId:   seriesId,
				EpisodeNums: torrentEpisodes,
				HashFilterFn: c.hashInBlacklist,
			})
			if err != nil {
				log.Warnf("download episode error, continue next item: %v", err)
				continue lo	
			}
			torrentNames = append(torrentNames, *name)

			for _, num := range torrentEpisodes {
				delete(wanted, num) //delete downloaded episode from wanted
			}
		}
	}
	if len(wanted) > 0 {
		log.Warnf("still wanted but not downloaded episodes: %v", wanted)
	}
	return torrentNames, nil
}

func (c *Engine) DownloadMovie(m *ent.Media, r1 torznab.Result) (*string, error) {
	return c.downloadTorrent(m, r1, DownloadOptions{
		SeasonNum: 0,
		MediaId:   m.ID,
	})
}

func (c *Engine) hashInBlacklist(hash string) bool {
	blacklist, err := c.db.GetTorrentBlacklist()
	if err!= nil {
		log.Warnf("get torrent blacklist error: %v", err)
		return false	
	}
	for _, b := range blacklist {
		if b.TorrentHash == hash {
			return true
		}	
	}
	return false
}

func (c *Engine) downloadTorrent(m *ent.Media, r1 torznab.Result, op DownloadOptions) (*string, error) {
	trc, dlc, err := c.GetDownloadClient()
	if err != nil {
		return nil, errors.Wrap(err, "get download client")
	}

	downloadDir := c.db.GetDownloadDir()
	
	//due to reported bug by user, this will be temporarily disabled
	// size := utils.AvailableSpace(downloadDir)
	// if size < uint64(r1.Size) {
	// 	log.Errorf("space available %v, space needed %v", size, r1.Size)
	// 	return nil, errors.New("not enough space")
	// }

	var name = r1.Name
	var targetDir = m.TargetDir
	if m.MediaType == media.MediaTypeTv { //tv download
		targetDir = fmt.Sprintf("%s/Season %02d/", m.TargetDir, op.SeasonNum)

		if len(op.EpisodeNums) > 0 {
			for _, epNum := range op.EpisodeNums {
				ep, err := c.db.GetEpisode(m.ID, op.SeasonNum, epNum)
				if err != nil {
					return nil, errors.Errorf("no episode of season %d episode %d", op.SeasonNum, epNum)

				}
				if ep.Status == episode.StatusMissing {
					c.db.SetEpisodeStatus(ep.ID, episode.StatusDownloading)
				}
			}
			buff := &bytes.Buffer{}
			for i, ep := range op.EpisodeNums {
				if i != 0 {
					buff.WriteString(",")

				}
				buff.WriteString(fmt.Sprint(ep))
			}
			name = fmt.Sprintf("第%s集 (%s)", buff.String(), name)

		} else { //season package download
			name = fmt.Sprintf("全集 (%s)", name)
			c.db.SetSeasonAllEpisodeStatus(m.ID, op.SeasonNum, episode.StatusDownloading)
		}

	} else {//movie download
		ep, _ := c.db.GetMovieDummyEpisode(m.ID)
		if ep.Status == episode.StatusMissing {
			c.db.SetEpisodeStatus(ep.ID, episode.StatusDownloading)
		}

	}
	
	link, hash, err := utils.GetRealLinkAndHash(r1.Link)
	if err != nil {
		return nil, errors.Wrap(err, "get hash")
	}

	if op.HashFilterFn != nil && op.HashFilterFn(hash) {
		return nil, errors.Errorf("hash is filtered: %s", hash)
	}

	r1.Link = link

	history, err := c.db.SaveHistoryRecord(ent.History{
		MediaID:     m.ID,
		EpisodeNums: op.EpisodeNums,
		SeasonNum:   op.SeasonNum,
		SourceTitle: r1.Name,
		TargetDir:   targetDir,
		Status:      history.StatusRunning,
		Size:        int(r1.Size),
		//Saved:            torrent.Save(),
		Link:             r1.Link,
		Hash:             hash,
		DownloadClientID: dlc.ID,
		IndexerID:        r1.IndexerId,
	})
	if err != nil {
		return nil, errors.Wrap(err, "save record")
	}

	torrent, err := trc.Download(r1.Link, hash, downloadDir)
	if err != nil {
		return nil, errors.Wrap(err, "downloading")
	}
	torrent.Start()

	c.tasks.Store(history.ID, &Task{Torrent: torrent})

	c.sendMsg(fmt.Sprintf(message.BeginDownload, name))

	log.Infof("success add %s to download task", r1.Name)

	return &r1.Name, nil
}
