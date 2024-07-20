package core

import (
	"fmt"
	"polaris/db"
	"polaris/ent"
	"polaris/ent/media"
	"polaris/log"
	"polaris/pkg/torznab"
	"polaris/pkg/utils"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func SearchSeasonPackage(db1 *db.Client, seriesId, seasonNum int, checkResolution bool) ([]torznab.Result, error) {
	series := db1.GetMediaDetails(seriesId)
	if series == nil {
		return nil, fmt.Errorf("no tv series of id %v", seriesId)
	}
	q := fmt.Sprintf("%s S%02d", series.NameEn, seasonNum)

	res := searchWithTorznab(db1, q)
	if len(res) == 0 {
		return nil, fmt.Errorf("no resource found")
	}
	var filtered []torznab.Result
	for _, r := range res {
		if !isNameAcceptable(r.Name, series.Media, seasonNum, -1) {
			continue
		}
		if checkResolution && !IsWantedResolution(r.Name, series.Resolution) {
			continue
		}
		
		filtered = append(filtered, r)
		
	}

	if len(filtered) == 0 {
		return nil, errors.New("no resource found")
	}
	return filtered, nil
}

func SearchEpisode(db1 *db.Client, seriesId, seasonNum, episodeNum int, checkResolution bool) ([]torznab.Result, error) {
	series := db1.GetMediaDetails(seriesId)
	if series == nil {
		return nil, fmt.Errorf("no tv series of id %v", seriesId)
	}

	q := fmt.Sprintf("%s S%02dE%02d", series.NameEn, seasonNum, episodeNum)
	res := searchWithTorznab(db1, q)
	if len(res) == 0 {
		return nil, fmt.Errorf("no resource found")
	}

	var filtered []torznab.Result
	for _, r := range res {
		if !isNameAcceptable(r.Name, series.Media, seasonNum, episodeNum) {
			continue
		}
		if checkResolution && !IsWantedResolution(r.Name, series.Resolution) {
			continue
		}

		filtered = append(filtered, r)
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
		if !isNameAcceptable(r.Name, movieDetail.Media, -1, -1) {
			continue
		}
		if checkResolution && !IsWantedResolution(r.Name, movieDetail.Resolution) {
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
	for _, tor := range allTorznab {
		resp, err := torznab.Search(tor.URL, tor.ApiKey, q)
		if err != nil {
			log.Errorf("search %s error: %v", tor.Name, err)
			continue
		}
		res = append(res, resp...)
	}
	sort.Slice(res, func(i, j int) bool {
		var s1 = res[i]
		var s2 = res[j]
		return s1.Seeders > s2.Seeders
	})

	return res
}

func isNameAcceptable(torrentName string, m *ent.Media, seasonNum, episodeNum int) bool {
	if !utils.IsNameAcceptable(torrentName, m.NameCn) && !utils.IsNameAcceptable(torrentName, m.NameEn)  && !utils.IsNameAcceptable(torrentName, m.OriginalName){
		return false //name not match
	}

	ss := strings.Split(m.AirDate, "-")[0]
	year, _ := strconv.Atoi(ss)
	if m.MediaType == media.MediaTypeMovie {
		if !strings.Contains(torrentName, strconv.Itoa(year)) && !strings.Contains(torrentName, strconv.Itoa(year+1)) && !strings.Contains(torrentName, strconv.Itoa(year-1)) {
			return false //not the same movie, if year is not correct
		}
	}

	if m.MediaType == media.MediaTypeTv {
		if episodeNum != -1 {
			se := fmt.Sprintf("S%02dE%02d", seasonNum, episodeNum)
			if !utils.ContainsIgnoreCase(torrentName, se) {
				return false 
			}	
		} else {
			//season package
			if !utils.IsSeasonPackageName(torrentName) {
				return false
			}

			seNum, err := utils.FindSeasonPackageInfo(torrentName)
			if err != nil {
				return false
			}
			if seNum != seasonNum {
				return false
			}
	
		}
	}
	return true
}

func IsWantedResolution(name string, res media.Resolution) bool {
	switch res {
	case media.Resolution720p:
		return utils.ContainsIgnoreCase(name, "720p")
	case media.Resolution1080p:
		return utils.ContainsIgnoreCase(name, "1080p")
	case media.Resolution4k:
		return utils.ContainsIgnoreCase(name, "4k") || utils.ContainsIgnoreCase(name, "2160p")
	}
	return false
}