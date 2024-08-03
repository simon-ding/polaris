package core

import (
	"fmt"
	"polaris/db"
	"polaris/log"
	"polaris/pkg/metadata"
	"polaris/pkg/torznab"
	"polaris/pkg/utils"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

func SearchTvSeries(db1 *db.Client, seriesId, seasonNum int, episodes []int, checkResolution bool, checkFileSize bool) ([]torznab.Result, error) {
	series := db1.GetMediaDetails(seriesId)
	if series == nil {
		return nil, fmt.Errorf("no tv series of id %v", seriesId)
	}
	slices.Contains(episodes, 1)

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
		if isNumberedSeries(series) && len(episodes) == 0 {
			//should not want season
			continue
		}

		if len(episodes) > 0 && slices.Contains(episodes, meta.Episode) { //not season pack, episode number equals
			continue
		}else if len(episodes) == 0 && !meta.IsSeasonPack { //want season pack, but not season pack
			continue
		}
		if checkResolution && meta.Resolution != series.Resolution.String() {
			continue
		}
		if !utils.IsNameAcceptable(meta.NameEn, series.NameEn) && !utils.IsNameAcceptable(meta.NameCn, series.NameCn) {
			continue
		}

		if checkFileSize && series.Limiter != nil  {
			if series.Limiter.SizeMin > 0 &&  r.Size < series.Limiter.SizeMin {
				//min size not satified
				continue
			}
			if series.Limiter.SizeMax > 0 && r.Size > series.Limiter.SizeMax {
				//max size not satified
				continue
			}
		}
		filtered = append(filtered, r)
	}
	if len(filtered) == 0 {
		return nil, errors.New("no resource found")
	}
	filtered = dedup(filtered)
	return filtered, nil

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
	return hasSeason2 && !season2HasEpisode1 //only one 1st episode
}

func SearchMovie(db1 *db.Client, movieId int, checkResolution bool, checkFileSize bool) ([]torznab.Result, error) {
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

		if checkFileSize && movieDetail.Limiter != nil  {
			if movieDetail.Limiter.SizeMin > 0 &&  r.Size < movieDetail.Limiter.SizeMin {
				//min size not satified
				continue
			}
			if movieDetail.Limiter.SizeMax > 0 && r.Size > movieDetail.Limiter.SizeMax {
				//max size not satified
				continue
			}
		}

		ss := strings.Split(movieDetail.AirDate, "-")[0]
		year, _ := strconv.Atoi(ss)
		if meta.Year != year && meta.Year != year-1 && meta.Year != year+1 { //year not match
			continue
		}

		filtered = append(filtered, r)

	}
	if len(filtered) == 0 {
		return nil, errors.New("no resource found")
	}
	filtered = dedup(filtered)

	return filtered, nil

}

func searchWithTorznab(db *db.Client, q string) []torznab.Result {

	var res []torznab.Result
	allTorznab := db.GetAllTorznabInfo()
	resChan := make(chan []torznab.Result)
	var wg sync.WaitGroup

	for _, tor := range allTorznab {
		if tor.Disabled {
			continue
		}
		wg.Add(1)
		go func() {
			log.Debugf("search torznab %v with %v", tor.Name, q)
			defer wg.Done()
			resp, err := torznab.Search(tor, tor.ApiKey, q)
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

	//res = dedup(res)

	sort.SliceStable(res, func(i, j int) bool { //先按做种人数排序
		var s1 = res[i]
		var s2 = res[j]
		return s1.Seeders > s2.Seeders
	})

	sort.SliceStable(res, func(i, j int) bool { //再按优先级排序，优先级高的种子排前面
		var s1 = res[i]
		var s2 = res[j]
		return s1.Priority > s2.Priority
	})

	//pt资源中，同一indexer内部，优先下载free的资源
	sort.SliceStable(res, func(i, j int) bool {
		var s1 = res[i]
		var s2 = res[j]
		if s1.IndexerId == s2.IndexerId && s1.IsPrivate { 
			return s1.DownloadVolumeFactor < s2.DownloadVolumeFactor
		}
		return false
	})

	//同一indexer内部，如果下载消耗一样，则优先下载上传奖励较多的
	sort.SliceStable(res, func(i, j int) bool { 
		var s1 = res[i]
		var s2 = res[j]
		if s1.IndexerId == s2.IndexerId && s1.IsPrivate && s1.DownloadVolumeFactor == s2.DownloadVolumeFactor{ 
			return s1.UploadVolumeFactor > s2.UploadVolumeFactor
		}
		return false
	})

	return res
}


func dedup(list []torznab.Result) []torznab.Result {
	var res = make([]torznab.Result, 0, len(list))
	seen := make(map[string]bool, 0)
	for _, r := range list {
		key := fmt.Sprintf("%s%s%d%d", r.Name, r.Source, r.Seeders,r.Peers)
		if seen[key] {
			continue
		}
		seen[key] = true
		res = append(res, r)
	}
	return res
}