package core

import (
	"bytes"
	"fmt"
	"polaris/ent"
	"polaris/ent/episode"
	"polaris/ent/history"
	"polaris/log"
	"polaris/pkg/metadata"
	"polaris/pkg/notifier/message"
	"polaris/pkg/torznab"
	"polaris/pkg/utils"

	"github.com/pkg/errors"
)

func (c *Client) DownloadEpisodeTorrent(r1 torznab.Result, seriesId, seasonNum int, episodeNums ...int) (*string, error) {
	trc, dlc, err := c.GetDownloadClient()
	if err != nil {
		return nil, errors.Wrap(err, "connect transmission")
	}
	series := c.db.GetMediaDetails(seriesId)
	if series == nil {
		return nil, fmt.Errorf("no tv series of id %v", seriesId)
	}

	//check space available
	downloadDir := c.db.GetDownloadDir()
	size := utils.AvailableSpace(downloadDir)
	if size < uint64(r1.Size) {
		log.Errorf("space available %v, space needed %v", size, r1.Size)
		return nil, errors.New("no enough space")
	}

	magnet, err := utils.Link2Magnet(r1.Link)
	if err != nil {
		return nil, errors.Errorf("converting link to magnet error, link: %v, error: %v", r1.Link, err)
	}

	torrent, err := trc.Download(magnet, downloadDir)
	if err != nil {
		return nil, errors.Wrap(err, "downloading")
	}
	torrent.Start()

	dir := fmt.Sprintf("%s/Season %02d/", series.TargetDir, seasonNum)

	if len(episodeNums) > 0 {
		for _, epNum := range episodeNums {
			var ep *ent.Episode
			for _, e := range series.Episodes {
				if e.SeasonNumber == seasonNum && e.EpisodeNumber == epNum {
					ep = e
				}
			}
			if ep == nil {
				return nil, errors.Errorf("no episode of season %d episode %d", seasonNum, epNum)
			}

			if ep.Status == episode.StatusMissing {
				c.db.SetEpisodeStatus(ep.ID, episode.StatusDownloading)
			}

		}
	} else { //season package download
		c.db.SetSeasonAllEpisodeStatus(seriesId, seasonNum, episode.StatusDownloading)

	}
	history, err := c.db.SaveHistoryRecord(ent.History{
		MediaID:     seriesId,
		EpisodeNums: episodeNums,
		SeasonNum:   seasonNum,
		SourceTitle: r1.Name,
		TargetDir:   dir,
		Status:      history.StatusRunning,
		Size:        r1.Size,
		//Saved:            torrent.Save(),
		Link:             magnet,
		DownloadClientID: dlc.ID,
		IndexerID:        r1.IndexerId,
	})
	if err != nil {
		return nil, errors.Wrap(err, "save record")
	}

	c.tasks[history.ID] = &Task{Torrent: torrent}
	name := r1.Name

	if len(episodeNums) > 0 {
		buff := &bytes.Buffer{}
		for i, ep := range episodeNums {
			if i != 0 {
				buff.WriteString(",")

			}
			buff.WriteString(fmt.Sprint(ep))
		}
		name = fmt.Sprintf("第%s集 (%s)", buff.String(), name)
	} else {
		name = fmt.Sprintf("全集 (%s)", name)
	}

	c.sendMsg(fmt.Sprintf(message.BeginDownload, name))

	log.Infof("success add %s to download task", r1.Name)
	return &r1.Name, nil
}

/*
tmdb 校验获取的资源名，如果用资源名在tmdb搜索出来的结果能匹配上想要的资源，则认为资源有效，否则无效
解决名称过于简单的影视会匹配过多资源的问题， 例如：梦魇绝镇 FROM
*/
func (c *Client) checkBtReourceWithTmdb(r *torznab.Result, seriesId int) bool {
	m := metadata.ParseTv(r.Name)
	se, err := c.MustTMDB().SearchMedia(m.NameEn, "", 1)
	if err != nil {
		log.Warnf("tmdb search error, consider this torrent ok: ", err)
		return true
	} else {
		if len(se.Results) == 0 {
			log.Debugf("tmdb search no result, consider this torrent ok: %s", r.Name) //because tv name parse is not accurate
			return true
		}
		series := c.db.GetMediaDetails(seriesId)

		se0 := se.Results[0]
		if se0.ID != int64(series.TmdbID) {
			log.Warnf("bt reosurce name not match tmdb id: %s", r.Name)
			return false
		} else { //resource tmdb id match
			return true
		}
	}
}

func (c *Client) SearchAndDownload(seriesId, seasonNum int, episodeNums ...int) ([]string, error) {

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
				name, err := c.DownloadEpisodeTorrent(r, seriesId, seasonNum)
				if err != nil {
					return nil, err
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
			name, err := c.DownloadEpisodeTorrent(r, seriesId, seasonNum, torrentEpisodes...)
			if err != nil {
				return nil, err
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

func (c *Client) DownloadMovie(m *ent.Media, link, name string, size int, indexerID int) (*string, error) {
	trc, dlc, err := c.GetDownloadClient()
	if err != nil {
		return nil, errors.Wrap(err, "connect transmission")
	}
	magnet, err := utils.Link2Magnet(link)
	if err != nil {
		return nil, errors.Errorf("converting link to magnet error, link: %v, error: %v", link, err)
	}

	torrent, err := trc.Download(magnet, c.db.GetDownloadDir())
	if err != nil {
		return nil, errors.Wrap(err, "downloading")
	}
	torrent.Start()

	if name == "" {
		name = m.OriginalName
	}
	go func() {
		ep, _ := c.db.GetMovieDummyEpisode(m.ID)
		history, err := c.db.SaveHistoryRecord(ent.History{
			MediaID:     m.ID,
			EpisodeID:   ep.ID,
			SourceTitle: name,
			TargetDir:   m.TargetDir,
			Status:      history.StatusRunning,
			Size:        size,
			//Saved:            torrent.Save(),
			Link:             magnet,
			DownloadClientID: dlc.ID,
			IndexerID:        indexerID,
		})
		if err != nil {
			log.Errorf("save history error: %v", err)
		}

		c.tasks[history.ID] = &Task{Torrent: torrent}

		if ep.Status == episode.StatusMissing {
			c.db.SetEpisodeStatus(ep.ID, episode.StatusDownloading)
		}

	}()

	c.sendMsg(fmt.Sprintf(message.BeginDownload, name))
	log.Infof("success add %s to download task", name)
	return &name, nil

}
