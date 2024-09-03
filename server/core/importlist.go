package core

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"polaris/db"
	"polaris/ent"
	"polaris/ent/importlist"
	"polaris/ent/media"
	"polaris/ent/schema"
	"polaris/log"
	"polaris/pkg/importlist/plexwatchlist"
	"regexp"
	"strings"
	"time"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/pkg/errors"
)

func (c *Client) periodicallyUpdateImportlist() error {
	log.Infof("begin check import list")
	lists, err := c.db.GetAllImportLists()
	if err != nil {
		return errors.Wrap(err, "get from db")
	}
	for _, l := range lists {
		log.Infof("check import list content for %v", l.Name)
		if l.Type == importlist.TypePlex {
			res, err := plexwatchlist.ParsePlexWatchlist(l.URL)
			if err != nil {
				log.Errorf("parse plex watchlist: %v", err)
				continue
			}
			for _, item := range res.Items {
				var tmdbRes *tmdb.FindByID
				if item.ImdbID != "" {
					tmdbRes1, err := c.MustTMDB().GetByImdbId(item.ImdbID, c.language)
					if err != nil {
						log.Errorf("get by imdb id error: %v", err)
						continue
					}
					tmdbRes = tmdbRes1
				} else if item.TvdbID != "" {
					tmdbRes1, err := c.MustTMDB().GetByTvdbId(item.TvdbID, c.language)
					if err != nil {
						log.Errorf("get by imdb id error: %v", err)
						continue
					}
					tmdbRes = tmdbRes1
				}
				if tmdbRes == nil {
					log.Errorf("can not find media for : %+v", item)
					continue
				}
				if len(tmdbRes.MovieResults) > 0 {
					d := tmdbRes.MovieResults[0]
					name, err := c.SuggestedMovieFolderName(int(d.ID))
					if err != nil {
						log.Errorf("suggesting name error: %v", err)
						continue
					}
					_, err = c.AddMovie2Watchlist(AddWatchlistIn{
						TmdbID:     int(d.ID),
						StorageID:  l.StorageID,
						Resolution: l.Qulity,
						Folder:     name,
					})
					if err != nil {
						log.Errorf("[update_import_lists] add movie to watchlist error: %v", err)
					} else {
						log.Infof("[update_import_lists] add movie to watchlist success")
					}
				} else if len(tmdbRes.TvResults) > 0 {
					d := tmdbRes.TvResults[0]
					name, err := c.SuggestedSeriesFolderName(int(d.ID))
					if err != nil {
						log.Errorf("suggesting name error: %v", err)
						continue
					}

					_, err = c.AddTv2Watchlist(AddWatchlistIn{
						TmdbID:     int(d.ID),
						StorageID:  l.StorageID,
						Resolution: l.Qulity,
						Folder:     name,
					})
					if err != nil {
						log.Errorf("[update_import_lists] add tv to watchlist error: %v", err)
					} else {
						log.Infof("[update_import_lists] add tv to watchlist success")
					}

				}

			}
		}
	}
	return nil
}

type AddWatchlistIn struct {
	TmdbID                  int    `json:"tmdb_id" binding:"required"`
	StorageID               int    `json:"storage_id" `
	Resolution              string `json:"resolution" binding:"required"`
	Folder                  string `json:"folder" binding:"required"`
	DownloadHistoryEpisodes bool   `json:"download_history_episodes"` //for tv
	SizeMin                 int    `json:"size_min"`
	SizeMax                 int    `json:"size_max"`
}

func (c *Client) AddTv2Watchlist(in AddWatchlistIn) (interface{}, error) {
	log.Debugf("add tv watchlist input %+v", in)
	if in.Folder == "" {
		return nil, errors.New("folder should be provided")
	}
	detailCn, err := c.MustTMDB().GetTvDetails(in.TmdbID, db.LanguageCN)
	if err != nil {
		return nil, errors.Wrap(err, "get tv detail")
	}
	var nameCn = detailCn.Name

	detailEn, _ := c.MustTMDB().GetTvDetails(in.TmdbID, db.LanguageEN)
	var nameEn = detailEn.Name
	var detail *tmdb.TVDetails
	if c.language == "" || c.language == db.LanguageCN {
		detail = detailCn
	} else {
		detail = detailEn
	}
	log.Infof("find detail for tv id %d: %v", in.TmdbID, detail)

	var epIds []int
	for _, season := range detail.Seasons {
		seasonId := season.SeasonNumber
		se, err := c.MustTMDB().GetSeasonDetails(int(detail.ID), seasonId, c.language)
		if err != nil {
			log.Errorf("get season detail (%s) error: %v", detail.Name, err)
			continue
		}
		for _, ep := range se.Episodes {
			shouldMonitor := false
			//如果设置下载往期剧集，则监控所有剧集。如果没有则监控未上映的剧集，考虑时差等问题留24h余量
			if in.DownloadHistoryEpisodes {
				shouldMonitor = true
			} else {
				t, err := time.Parse("2006-01-02", ep.AirDate)
				if err != nil {
					log.Error("air date not known, will monitor: %v", ep.AirDate)
					shouldMonitor = true

				} else {
					if time.Since(t) < 24*time.Hour { //monitor episode air 24h before now
						shouldMonitor = true
					}
				}
			}

			epid, err := c.db.SaveEposideDetail(&ent.Episode{
				SeasonNumber:  seasonId,
				EpisodeNumber: ep.EpisodeNumber,
				Title:         ep.Name,
				Overview:      ep.Overview,
				AirDate:       ep.AirDate,
				Monitored:     shouldMonitor,
			})
			if err != nil {
				log.Errorf("save episode info error: %v", err)
				continue
			}
			epIds = append(epIds, epid)
		}
	}
	m := &ent.Media{
		TmdbID:                  int(detail.ID),
		ImdbID:                  detail.IMDbID,
		MediaType:               media.MediaTypeTv,
		NameCn:                  nameCn,
		NameEn:                  nameEn,
		OriginalName:            detail.OriginalName,
		Overview:                detail.Overview,
		AirDate:                 detail.FirstAirDate,
		Resolution:              media.Resolution(in.Resolution),
		StorageID:               in.StorageID,
		TargetDir:               in.Folder,
		DownloadHistoryEpisodes: in.DownloadHistoryEpisodes,
		Limiter:                 schema.MediaLimiter{SizeMin: in.SizeMin, SizeMax: in.SizeMax},
		Extras: schema.MediaExtras{
			OriginalLanguage: detail.OriginalLanguage,
			Genres:           detail.Genres,
		},
	}

	r, err := c.db.AddMediaWatchlist(m, epIds)
	if err != nil {
		return nil, errors.Wrap(err, "add to list")
	}
	go func() {
		if err := c.downloadPoster(detail.PosterPath, r.ID); err != nil {
			log.Errorf("download poster error: %v", err)
		}
		if err := c.downloadBackdrop(detail.BackdropPath, r.ID); err != nil {
			log.Errorf("download poster error: %v", err)
		}
		if err := c.CheckDownloadedSeriesFiles(r); err != nil {
			log.Errorf("check downloaded files error: %v", err)
		}

	}()

	log.Infof("add tv %s to watchlist success", detail.Name)
	return nil, nil
}

func (c *Client) AddMovie2Watchlist(in AddWatchlistIn) (interface{}, error) {
	log.Infof("add movie watchlist input: %+v", in)
	detailCn, err := c.MustTMDB().GetMovieDetails(in.TmdbID, db.LanguageCN)
	if err != nil {
		return nil, errors.Wrap(err, "get movie detail")
	}
	var nameCn = detailCn.Title

	detailEn, _ := c.MustTMDB().GetMovieDetails(in.TmdbID, db.LanguageEN)
	var nameEn = detailEn.Title
	var detail *tmdb.MovieDetails
	if c.language == "" || c.language == db.LanguageCN {
		detail = detailCn
	} else {
		detail = detailEn
	}
	log.Infof("find detail for movie id %d: %v", in.TmdbID, detail)

	epid, err := c.db.SaveEposideDetail(&ent.Episode{
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "dummy episode for movies",
		Overview:      "dummy episode for movies",
		AirDate:       detail.ReleaseDate,
		Monitored:     true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "add dummy episode")
	}
	log.Infof("added dummy episode for movie: %v", nameEn)

	movie := ent.Media{
		TmdbID:       int(detail.ID),
		ImdbID:       detail.IMDbID,
		MediaType:    media.MediaTypeMovie,
		NameCn:       nameCn,
		NameEn:       nameEn,
		OriginalName: detail.OriginalTitle,
		Overview:     detail.Overview,
		AirDate:      detail.ReleaseDate,
		Resolution:   media.Resolution(in.Resolution),
		StorageID:    in.StorageID,
		TargetDir:    in.Folder,
		Limiter:      schema.MediaLimiter{SizeMin: in.SizeMin, SizeMax: in.SizeMax},
	}

	extras := schema.MediaExtras{
		IsAdultMovie:     detail.Adult,
		OriginalLanguage: detail.OriginalLanguage,
		Genres:           detail.Genres,
	}
	if IsJav(detail) {
		javid := c.GetJavid(in.TmdbID)
		extras.JavId = javid
	}

	movie.Extras = extras
	r, err := c.db.AddMediaWatchlist(&movie, []int{epid})
	if err != nil {
		return nil, errors.Wrap(err, "add to list")
	}
	go func() {
		if err := c.downloadPoster(detail.PosterPath, r.ID); err != nil {
			log.Errorf("download poster error: %v", err)
		}
		if err := c.downloadBackdrop(detail.BackdropPath, r.ID); err != nil {
			log.Errorf("download backdrop error: %v", err)
		}
	}()

	log.Infof("add movie %s to watchlist success", detail.Title)
	return nil, nil

}

func IsJav(detail *tmdb.MovieDetails) bool {
	if detail.Adult && len(detail.ProductionCountries) > 0 && strings.ToUpper(detail.ProductionCountries[0].Iso3166_1) == "JP" {
		return true
	}
	return false
}

func (c *Client) GetJavid(id int) string {
	alters, err := c.MustTMDB().GetMovieAlternativeTitles(id, c.language)
	if err != nil {
		return ""
	}
	for _, t := range alters.Titles {
		if t.Iso3166_1 == "JP" && t.Type == "" {
			return t.Title
		}
	}
	return ""
}

func (c *Client) downloadBackdrop(path string, mediaID int) error {
	url := "https://image.tmdb.org/t/p/original" + path
	return c.downloadImage(url, mediaID, "backdrop.jpg")
}

func (c *Client) downloadPoster(path string, mediaID int) error {
	var url = "https://image.tmdb.org/t/p/original" + path

	return c.downloadImage(url, mediaID, "poster.jpg")
}

func (c *Client) downloadImage(url string, mediaID int, name string) error {

	log.Infof("try to download image: %v", url)
	var resp, err = http.Get(url)
	if err != nil {
		return errors.Wrap(err, "http get")
	}
	targetDir := fmt.Sprintf("%v/%d", db.ImgPath, mediaID)
	os.MkdirAll(targetDir, 0777)
	//ext := filepath.Ext(path)
	targetFile := filepath.Join(targetDir, name)
	f, err := os.Create(targetFile)
	if err != nil {
		return errors.Wrap(err, "new file")
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return errors.Wrap(err, "copy http response")
	}
	log.Infof("image successfully downlaoded: %v", targetFile)
	return nil

}

func (c *Client) SuggestedMovieFolderName(tmdbId int) (string, error) {

	d1, err := c.MustTMDB().GetMovieDetails(tmdbId, c.language)
	if err != nil {
		return "", errors.Wrap(err, "get movie details")
	}
	name := d1.Title

	if IsJav(d1) {
		javid := c.GetJavid(tmdbId)
		if javid != "" {
			return javid, nil
		}
	}

	if c.language == db.LanguageCN {
		en, err := c.MustTMDB().GetMovieDetails(tmdbId, db.LanguageEN)
		if err != nil {
			log.Errorf("get en movie detail error: %v", err)
		} else {
			name = fmt.Sprintf("%s %s", d1.Title, en.Title)
		}
	}
	//remove extra characters
	re := regexp.MustCompile(`[^\p{L}\w\s]`)
	name = re.ReplaceAllString(strings.ToLower(name), " ")
	name = strings.Join(strings.Fields(name), " ")
	year := strings.Split(d1.ReleaseDate, "-")[0]
	if year != "" {
		name = fmt.Sprintf("%s (%s)", name, year)
	}

	log.Infof("tv series of tmdb id %v suggestting name is %v", tmdbId, name)
	return name, nil
}

func (c *Client) SuggestedSeriesFolderName(tmdbId int) (string, error) {

	d, err := c.MustTMDB().GetTvDetails(tmdbId, c.language)
	if err != nil {
		return "", errors.Wrap(err, "get tv details")
	}

	name := d.Name

	if c.language == db.LanguageCN {
		en, err := c.MustTMDB().GetTvDetails(tmdbId, db.LanguageEN)
		if err != nil {
			log.Errorf("get en tv detail error: %v", err)
		} else {
			name = fmt.Sprintf("%s %s", d.Name, en.Name)
		}
	}
	//remove extra characters
	re := regexp.MustCompile(`[^\p{L}\w\s]`)
	name = re.ReplaceAllString(strings.ToLower(name), " ")
	name = strings.Join(strings.Fields(name), " ")
	year := strings.Split(d.FirstAirDate, "-")[0]
	if year != "" {
		name = fmt.Sprintf("%s (%s)", name, year)
	}

	log.Infof("tv series of tmdb id %v suggestting name is %v", tmdbId, name)
	return name, nil
}
