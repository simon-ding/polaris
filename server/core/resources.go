package core

import (
	"fmt"
	"polaris/ent"
	"polaris/ent/episode"
	"polaris/ent/history"
	"polaris/log"
	"polaris/pkg/notifier/message"
	"polaris/pkg/torznab"
	"polaris/pkg/utils"

	"github.com/pkg/errors"
)

func (c *Client) DownloadSeasonPackage(r1 torznab.Result, seriesId, seasonNum int) (*string, error) {
	trc, dlClient, err := c.getDownloadClient()
	if err != nil {
		return nil, errors.Wrap(err, "connect transmission")
	}
	downloadDir := c.db.GetDownloadDir()
	size := utils.AvailableSpace(downloadDir)
	if size < uint64(r1.Size) {
		log.Errorf("space available %v, space needed %v", size, r1.Size)
		return nil, errors.New("no enough space")
	}

	torrent, err := trc.Download(r1.Link, c.db.GetDownloadDir())
	if err != nil {
		return nil, errors.Wrap(err, "downloading")
	}
	torrent.Start()

	series := c.db.GetMediaDetails(seriesId)
	if series == nil {
		return nil, fmt.Errorf("no tv series of id %v", seriesId)
	}
	dir := fmt.Sprintf("%s/Season %02d/", series.TargetDir, seasonNum)

	history, err := c.db.SaveHistoryRecord(ent.History{
		MediaID:          seriesId,
		EpisodeID:        0,
		SourceTitle:      r1.Name,
		TargetDir:        dir,
		Status:           history.StatusRunning,
		Size:             r1.Size,
		Saved:            torrent.Save(),
		DownloadClientID: dlClient.ID,
		IndexerID:        r1.IndexerId,
	})
	if err != nil {
		return nil, errors.Wrap(err, "save record")
	}
	c.db.SetSeasonAllEpisodeStatus(seriesId, seasonNum, episode.StatusDownloading)

	c.tasks[history.ID] = &Task{Torrent: torrent}

	c.sendMsg(fmt.Sprintf(message.BeginDownload, r1.Name))
	return &r1.Name, nil

}


func (c *Client) DownloadEpisodeTorrent(r1 torznab.Result, seriesId, seasonNum, episodeNum int) (*string, error) {
	trc, dlc, err := c.getDownloadClient()
	if err != nil {
		return nil, errors.Wrap(err, "connect transmission")
	}
	series := c.db.GetMediaDetails(seriesId)
	if series == nil {
		return nil, fmt.Errorf("no tv series of id %v", seriesId)
	}
	var ep *ent.Episode
	for _, e := range series.Episodes {
		if e.SeasonNumber == seasonNum && e.EpisodeNumber == episodeNum {
			ep = e
		}
	}
	if ep == nil {
		return nil, errors.Errorf("no episode of season %d episode %d", seasonNum, episodeNum)
	}
	torrent, err := trc.Download(r1.Link, c.db.GetDownloadDir())
	if err != nil {
		return nil, errors.Wrap(err, "downloading")
	}
	torrent.Start()

	dir := fmt.Sprintf("%s/Season %02d/", series.TargetDir, seasonNum)

	history, err := c.db.SaveHistoryRecord(ent.History{
		MediaID:          ep.MediaID,
		EpisodeID:        ep.ID,
		SourceTitle:      r1.Name,
		TargetDir:        dir,
		Status:           history.StatusRunning,
		Size:             r1.Size,
		Saved:            torrent.Save(),
		DownloadClientID: dlc.ID,
		IndexerID:        r1.IndexerId,
	})
	if err != nil {
		return nil, errors.Wrap(err, "save record")
	}
	c.db.SetEpisodeStatus(ep.ID, episode.StatusDownloading)

	c.tasks[history.ID] = &Task{Torrent: torrent}
	c.sendMsg(fmt.Sprintf(message.BeginDownload, r1.Name))

	log.Infof("success add %s to download task", r1.Name)
	return &r1.Name, nil

}
func (c *Client) SearchAndDownload(seriesId, seasonNum, episodeNum int) (*string, error) {

	res, err := SearchTvSeries(c.db, seriesId, seasonNum, []int{episodeNum}, true, true)
	if err != nil {
		return nil, err
	}
	r1 := res[0]
	log.Infof("found resource to download: %+v", r1)
	return c.DownloadEpisodeTorrent(r1, seriesId, seasonNum, episodeNum)
}

func (c *Client) DownloadMovie(m *ent.Media,link, name string, size int, indexerID int) (*string, error) {
	trc, dlc, err := c.getDownloadClient()
	if err != nil {
		return nil, errors.Wrap(err, "connect transmission")
	}

	torrent, err := trc.Download(link, c.db.GetDownloadDir())
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
			MediaID:          m.ID,
			EpisodeID:        ep.ID,
			SourceTitle:      name,
			TargetDir:        m.TargetDir,
			Status:           history.StatusRunning,
			Size:             size,
			Saved:            torrent.Save(),
			DownloadClientID: dlc.ID,
			IndexerID:        indexerID,
		})
		if err != nil {
			log.Errorf("save history error: %v", err)
		}

		c.tasks[history.ID] = &Task{Torrent: torrent}

		c.db.SetEpisodeStatus(ep.ID, episode.StatusDownloading)
	}()

	c.sendMsg(fmt.Sprintf(message.BeginDownload, name))
	log.Infof("success add %s to download task", name)
	return &name, nil

}