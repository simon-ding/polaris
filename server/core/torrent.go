package core

import (
	"fmt"
	"polaris/db"
	"polaris/log"
	"polaris/pkg/metadata"
	"polaris/pkg/torznab"
	"polaris/pkg/utils"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

func SearchSeasonPackage(db1 *db.Client, seriesId, seasonNum int, checkResolution bool) ([]torznab.Result, error) {
	return SearchEpisode(db1, seriesId, seasonNum, -1, checkResolution)
}

func isNumberedSeries(detail *db.MediaDetails) bool {
	hasSeason2 := false
	season2HasEpisode1 := false
	for _, ep := range detail.Episodes {
		if ep.SeasonNumber == 2 {
			hasSeason2 = true
			if ep.EpisodeNumber == 1 {
				season2HasEpisode1 = true
			}
	
		}
	}
	return hasSeason2 && !season2HasEpisode1//only one 1st episode
}

func SearchEpisode(db1 *db.Client, seriesId, seasonNum, episodeNum int, checkResolution bool) ([]torznab.Result, error) {
	series := db1.GetMediaDetails(seriesId)
	if series == nil {
		return nil, fmt.Errorf("no tv series of id %v", seriesId)
	}

	res := searchWithTorznab(db1, series.NameEn)
	resCn := searchWithTorznab(db1, series.NameCn)
	res = append(res, resCn...)

	var filtered []torznab.Result
	for _, r := range res {
		//log.Infof("torrent resource: %+v", r)
		meta := metadata.ParseTv(r.Name)
		if meta == nil { //cannot parse name
			continue
		}
		if !isNumberedSeries(series) { //do not check season on series that only rely on episode number
			if meta.Season != seasonNum {
				continue
			}	
		}
		if isNumberedSeries(series) && episodeNum == -1 {
			//should not want season
			continue
		}

		if episodeNum != -1 && meta.Episode != episodeNum { //not season pack, episode number equals
			continue
		} else if seasonNum == -1 && !meta.IsSeasonPack { //want season pack, but not season pack
			continue
		}
		if checkResolution && meta.Resolution != series.Resolution.String() {
			continue
		}
		if !utils.IsNameAcceptable(meta.NameEn, series.NameEn) && !utils.IsNameAcceptable(meta.NameCn, series.NameCn) {
			continue
		}
		filtered = append(filtered, r)
	}
	if len(filtered) == 0 {
		return nil, errors.New("no resource found")
	}

	return filtered, nil

}

func SearchMovie(db1 *db.Client, movieId int, checkResolution bool) ([]torznab.Result, error) {
	movieDetail := db1.GetMediaDetails(movieId)
	if movieDetail == nil {
		return nil, errors.New("no media found of id")
	}

	res := searchWithTorznab(db1, movieDetail.NameEn)

	res1 := searchWithTorznab(db1, movieDetail.NameCn)
	res = append(res, res1...)

	if len(res) == 0 {
		return nil, fmt.Errorf("no resource found")
	}
	var filtered []torznab.Result
	for _, r := range res {
		meta := metadata.ParseMovie(r.Name)
		if !utils.IsNameAcceptable(meta.NameEn, movieDetail.NameEn) {
			continue
		}
		if checkResolution && meta.Resolution != movieDetail.Resolution.String() {
			continue
		}
		ss := strings.Split(movieDetail.AirDate, "-")[0]
		year, _ := strconv.Atoi(ss)
		if meta.Year != year && meta.Year != year-1 && meta.Year != year+1 { //year not match
			continue
		}
		if utils.ContainsIgnoreCase(r.Name, "soundtrack") {
			//ignore soundtracks
			continue
		}

		filtered = append(filtered, r)

	}
	if len(filtered) == 0 {
		return nil, errors.New("no resource found")
	}

	return filtered, nil

}

func searchWithTorznab(db *db.Client, q string) []torznab.Result {

	var res []torznab.Result
	allTorznab := db.GetAllTorznabInfo()
	resChan := make(chan []torznab.Result)
	var wg sync.WaitGroup

	for _, tor := range allTorznab {
		wg.Add(1)
		go func ()  {
			log.Debugf("search torznab %v with %v", tor.Name, q)
			defer wg.Done()
			resp, err := torznab.Search(tor.URL, tor.ApiKey, q)
			if err != nil {
				log.Errorf("search %s error: %v", tor.Name, err)
				return
			}
			resChan <- resp
	
		}()
	}
	go func() {
		wg.Wait()
		close(resChan) // 在所有的worker完成后关闭Channel
	}()

	for result := range resChan {
		res = append(res, result...)
	}

	sort.Slice(res, func(i, j int) bool {
		var s1 = res[i]
		var s2 = res[j]
		return s1.Seeders > s2.Seeders
	})

	return res
}
