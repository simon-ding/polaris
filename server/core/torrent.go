package core

import (
	"fmt"
	"polaris/db"
	"polaris/ent"
	"polaris/ent/media"
	"polaris/log"
	"polaris/pkg/gemini"
	"polaris/pkg/metadata"
	"polaris/pkg/prowlarr"
	"polaris/pkg/torznab"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type SearchParam struct {
	MediaId         int
	SeasonNum       int   //for tv
	Episodes        []int //for tv
	CheckResolution bool
	CheckFileSize   bool
	FilterQiangban  bool //for movie, 是否过滤枪版电影
}

func names2Query(media *ent.Media) []string {
	var names = []string{media.NameEn}

	if media.NameCn != "" {
		hasName := false
		for _, n := range names {
			if media.NameCn == n {
				hasName = true
			}
		}
		if !hasName {
			names = append(names, media.NameCn)
		}

	}
	if media.OriginalName != "" {
		hasName := false
		for _, n := range names {
			if media.OriginalName == n {
				hasName = true
			}
		}
		if !hasName {
			names = append(names, media.OriginalName)
		}

	}

	for _, t := range media.AlternativeTitles {
		if (t.Iso3166_1 == "CN" || t.Iso3166_1 == "US") && t.Type == "" {
			hasName := false
			for _, n := range names {
				if t.Title == n {
					hasName = true
				}
			}
			if !hasName {
				names = append(names, t.Title)
			}
		}
	}
	log.Debugf("name to query %+v", names)
	return names
}

func filterBasedOnGemini(cfg db.AIConfig, res []torznab.Result, names ...string) []torznab.Result {

	var torrentNames []string
	for _, r := range res {
		torrentNames = append(torrentNames, r.Name)
	}
	g, err := gemini.NewClient(cfg.GeminiApiKey, cfg.GeminiModelName)
	if err != nil {
		log.Warnf("create gemini client: %v", err)
		return res
	}
	resf, err := g.FilterTvOrMovies(torrentNames, names...)
	if err != nil {
		log.Warnf("filter with gemini: %v", err)
		return res
	}
	var newRes []torznab.Result
	for _, r := range res {
		if slices.Contains(resf, r.Name) {
			newRes = append(newRes, r)
		}
	}
	return newRes
}

func filterBasedOnRules(res []torznab.Result, names ...string) []torznab.Result {
	var filtered []torznab.Result
	for _, r := range res {
		meta := metadata.ParseTv(r.Name)
		if meta.IsAcceptable(names...) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func filterResourceNames(db1 *db.Client, res []torznab.Result, names ...string) []torznab.Result {
	n1 := len(res)
	cfg, err := db1.GetAIConfig()
	if err != nil {
		log.Warnf("get ai config: %v", err)
	}
	if cfg.Enabled {
		res = filterBasedOnGemini(cfg, res, names...)
	} else {
		res = filterBasedOnRules(res, names...)
	}
	log.Infof("resource before name filtering length is %d, after filtering length is %d", n1, len(res))
	return res

}

func SearchTvSeries(db1 *db.Client, param *SearchParam) ([]torznab.Result, error) {
	series := db1.GetMediaDetails(param.MediaId)
	if series == nil {
		return nil, fmt.Errorf("no tv series of id %v", param.MediaId)
	}
	limiter, err := db1.GetSizeLimiter("tv")
	if err != nil {
		log.Warnf("get tv size limiter: %v", err)
		limiter = &db.MediaSizeLimiter{}
	}
	log.Debugf("check tv series %s, season %d, episode %v", series.NameEn, param.SeasonNum, param.Episodes)

	names := names2Query(series.Media)

	res := searchWithTorznab(db1, prowlarr.TV, names...)

	res = filterResourceNames(db1, res, names...)

	var filtered []torznab.Result
lo:
	for _, r := range res {
		//log.Infof("torrent resource: %+v", r)
		meta := metadata.ParseTv(r.Name)
		meta.ParseExtraDescription(r.Description)

		if isImdbidNotMatch(series.ImdbID, r.ImdbId) { //has imdb id and not match
			continue
		}

		if !imdbIDMatchExact(series.ImdbID, r.ImdbId) { //imdb id not exact match, check file name
			if !torrentNameOk(series, meta) {
				continue
			}
		}

		if !isNoSeasonSeries(series) && meta.Season != param.SeasonNum { //do not check season on series that only rely on episode number
			continue

		}
		if isNoSeasonSeries(series) && len(param.Episodes) == 0 {
			//should not want season
			continue
		}

		if len(param.Episodes) > 0 { //not season pack, but episode number not equal
			if meta.StartEpisode <= 0 {
				continue lo
			}
			for i := meta.StartEpisode; i <= meta.EndEpisode; i++ {
				if !slices.Contains(param.Episodes, i) {
					continue lo
				}
			}
		} else if len(param.Episodes) == 0 && !meta.IsSeasonPack { //want season pack, but not season pack
			continue
		}

		if param.CheckResolution &&
			series.Resolution != media.ResolutionAny &&
			meta.Resolution != series.Resolution.String() {
			continue
		}

		if !torrentSizeOk(series, limiter, r.Size, meta.EndEpisode+1-meta.StartEpisode, param) {
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

// imdbid not exist consider match
func isImdbidNotMatch(id1, id2 string) bool {
	if id1 == "" || id2 == "" {
		return false
	}
	id1 = strings.TrimPrefix(id1, "tt")
	id2 = strings.TrimPrefix(id2, "tt")
	return id1 != id2
}

// imdbid not exist consider not match
func imdbIDMatchExact(id1, id2 string) bool {
	if id1 == "" || id2 == "" {
		return false
	}
	id1 = strings.TrimPrefix(id1, "tt")
	id2 = strings.TrimPrefix(id2, "tt")
	return id1 == id2
}

func torrentSizeOk(detail *db.MediaDetails, globalLimiter *db.MediaSizeLimiter, torrentSize int64,
	torrentEpisodeNum int, param *SearchParam) bool {

	multiplier := 1 //大小倍数，正常为1，如果是季包则为季内集数
	if detail.MediaType == media.MediaTypeTv {
		if len(param.Episodes) == 0 { //want tv season pack
			multiplier = seasonEpisodeCount(detail, param.SeasonNum)
		} else {
			multiplier = torrentEpisodeNum
		}
	}

	if param.CheckFileSize { //check file size when trigger automatic download

		if detail.Limiter.SizeMin > 0 { //min size
			sizeMin := detail.Limiter.SizeMin * int64(multiplier)
			if torrentSize < sizeMin { //比最小要求的大小还要小, min size not qualify
				return false
			}
		} else if globalLimiter != nil {
			resLimiter := globalLimiter.GetLimiter(detail.Resolution)
			sizeMin := resLimiter.MinSize * int64(multiplier)
			if torrentSize < sizeMin { //比最小要求的大小还要小, min size not qualify
				return false
			}
		}

		if detail.Limiter.SizeMax > 0 { //max size
			sizeMax := detail.Limiter.SizeMax * int64(multiplier)
			if torrentSize > sizeMax { //larger than max size wanted, max size not qualify
				return false
			}
		} else if globalLimiter != nil {
			resLimiter := globalLimiter.GetLimiter(detail.Resolution)
			sizeMax := resLimiter.MaxSIze * int64(multiplier)
			if torrentSize > sizeMax { //larger than max size wanted, max size not qualify
				return false
			}
		}
	}
	return true
}

func seasonEpisodeCount(detail *db.MediaDetails, seasonNum int) int {
	count := 0
	for _, ep := range detail.Episodes {
		if ep.SeasonNumber == seasonNum {
			count++
		}
	}
	return count
}

func isNoSeasonSeries(detail *db.MediaDetails) bool {
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

func SearchMovie(db1 *db.Client, param *SearchParam) ([]torznab.Result, error) {
	movieDetail := db1.GetMediaDetails(param.MediaId)
	if movieDetail == nil {
		return nil, errors.New("no media found of id")
	}

	limiter, err := db1.GetSizeLimiter("movie")
	if err != nil {
		log.Warnf("get tv size limiter: %v", err)
		limiter = &db.MediaSizeLimiter{}
	}
	names := names2Query(movieDetail.Media)

	res := searchWithTorznab(db1, prowlarr.Movie, names...)
	res = filterResourceNames(db1, res, names...)

	if movieDetail.Extras.IsJav() {
		res1 := searchWithTorznab(db1, prowlarr.Movie, movieDetail.Extras.JavId)
		res = append(res, res1...)
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("no resource found")
	}
	var filtered []torznab.Result
	for _, r := range res {
		meta := metadata.ParseMovie(r.Name)

		if isImdbidNotMatch(movieDetail.ImdbID, r.ImdbId) { //imdb id not match
			continue
		}

		if !imdbIDMatchExact(movieDetail.ImdbID, r.ImdbId) {
			if !torrentNameOk(movieDetail, meta) {
				continue
			}
			if !movieDetail.Extras.IsJav() {
				ss := strings.Split(movieDetail.AirDate, "-")[0]
				year, _ := strconv.Atoi(ss)
				if meta.Year != year && meta.Year != year-1 && meta.Year != year+1 { //year not match
					continue
				}
			}
		}

		if param.CheckResolution &&
			movieDetail.Resolution != media.ResolutionAny &&
			meta.Resolution != movieDetail.Resolution.String() {
			continue
		}

		if param.FilterQiangban && meta.IsQingban { //过滤枪版电影
			continue
		}

		if !torrentSizeOk(movieDetail, limiter, r.Size, 1, param) {
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

func searchWithTorznab(db *db.Client, t prowlarr.ProwlarrSupportType, queries ...string) []torznab.Result {

	var res []torznab.Result
	allTorznab := db.GetAllTorznabInfo()

	p, err := db.GetProwlarrSetting()
	if err == nil && !p.Disabled { //prowlarr exists
		c := prowlarr.New(p.ApiKey, p.URL)
		all, err := c.GetIndexers(t)
		if err != nil {
			log.Warnf("get prowlarr all indexer error: %v", err)
		} else {
			allTorznab = append(allTorznab, all...)
		}
	}
	resChan := make(chan []torznab.Result)
	var wg sync.WaitGroup

	for _, tor := range allTorznab {
		if tor.Disabled {
			continue
		}
		for _, q := range queries {
			wg.Add(1)

			go func() {
				log.Debugf("search torznab %v with %v", tor.Name, queries)
				defer wg.Done()

				resp, err := torznab.Search(tor, q)
				if err != nil {
					log.Warnf("search %s with query %s error: %v", tor.Name, q, err)
					return
				}
				resChan <- resp
			}()
		}
	}
	go func() {
		wg.Wait()
		close(resChan) // 在所有的worker完成后关闭Channel
	}()

	for result := range resChan {
		res = append(res, result...)
	}

	res = dedup(res)

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
		if s1.IndexerId == s2.IndexerId && s1.IsPrivate && s1.DownloadVolumeFactor == s2.DownloadVolumeFactor {
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
		key := fmt.Sprintf("%s%s%d%d", r.Name, r.Source, r.Seeders, r.Peers)
		if seen[key] {
			continue
		}
		seen[key] = true
		res = append(res, r)
	}
	return res
}

type NameTester interface {
	IsAcceptable(names ...string) bool
}

func torrentNameOk(detail *db.MediaDetails, tester NameTester) bool {
	if detail.Extras.IsJav() && tester.IsAcceptable(detail.Extras.JavId) {
		return true
	}
	names := names2Query(detail.Media)

	return tester.IsAcceptable(names...)
}
