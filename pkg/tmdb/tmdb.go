package tmdb

import (
	"polaris/log"
	"strconv"
	"strings"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/pkg/errors"
)

type Client struct {
	apiKey     string
	tmdbClient *tmdb.Client
}

func NewClient(apiKey string) (*Client, error) {
	tmdbClient, err := tmdb.Init(apiKey)
	if err != nil {
		return nil, errors.Wrap(err, "new tmdb client")
	}

	return &Client{
		apiKey:     apiKey,
		tmdbClient: tmdbClient,
	}, nil
}

func (c *Client) GetTvDetails(id int, language string) (*tmdb.TVDetails, error) {
	d, err := c.tmdbClient.GetTVDetails(id, withLangOption(language))
	if err != nil {
		return nil, errors.Wrap(err, "get tv detail")
	}

	log.Infof("tv id %d, language %s", id, language)
	if !episodeNameUseful(d.LastEpisodeToAir.Name) {
		log.Debug("should fetch english version")
		var detailEN *tmdb.TVDetails
		if language == "zh-CN" {
			detailEN, err = c.tmdbClient.GetTVDetails(id, withLangOption("en-US"))
			if err != nil {
				return d, nil
			}
			if episodeNameUseful(detailEN.LastEpisodeToAir.Name) {
				d.LastEpisodeToAir.Name = detailEN.LastEpisodeToAir.Name
				d.LastEpisodeToAir.Overview = detailEN.LastEpisodeToAir.Overview
				d.NextEpisodeToAir.Name = detailEN.NextEpisodeToAir.Name
				d.NextEpisodeToAir.Overview = detailEN.NextEpisodeToAir.Overview
			}
		}
	}


	return d, err
}

func (c *Client) GetMovieDetails(id int, language string) (*tmdb.MovieDetails, error) {
	return c.tmdbClient.GetMovieDetails(id, withLangOption(language))
}

func (c *Client) SearchTvShow(query string, lang string) (*tmdb.SearchTVShows, error) {
	r, err := c.tmdbClient.GetSearchTVShow(query, withLangOption(lang))
	if err != nil {
		return nil, errors.Wrap(err, "tmdb search tv")
	}
	return r, nil
}

type SearchResult struct {
	Page         int                 `json:"page"`
	Results      []*SearchResultItem `json:"results"`
	TotalResults int64               `json:"total_results"`
	TotalPages   int64               `json:"total_pages"`
}

type SearchResultItem struct {
	PosterPath       string   `json:"poster_path"`
	ID               int64    `json:"id"`
	Overview         string   `json:"overview"`
	MediaType        string   `json:"media_type"`
	FirstAirDate     string   `json:"first_air_date"`
	OriginCountry    []string `json:"origin_country"`
	GenreIDs         []int64  `json:"genre_ids"`
	OriginalLanguage string   `json:"original_language"`
	Name             string   `json:"name"`
	OriginalName     string   `json:"original_name"`
	Adult            bool     `json:"adult"`
	InWatchlist      bool     `json:"in_watchlist"`
}

func (c *Client) SearchMedia(query string, lang string, page int) (*SearchResult, error) {
	if page == 0 {
		page = 1
	}
	options := withLangOption(lang)
	options["page"] = strconv.Itoa(page)
	res, err := c.tmdbClient.GetSearchMulti(query, options)
	if err != nil {
		return nil, errors.Wrap(err, "query imdb")
	}

	searchResult := &SearchResult{
		Page:         res.Page,
		TotalResults: res.TotalResults,
		TotalPages:   res.TotalPages,
	}

	for _, r := range res.Results {
		if r.MediaType != "tv" && r.MediaType != "movie" {
			continue
		}
		item := &SearchResultItem{
			PosterPath:       r.PosterPath,
			ID:               r.ID,
			Overview:         r.Overview,
			MediaType:        r.MediaType,
			OriginCountry:    r.OriginCountry,
			OriginalLanguage: r.OriginalLanguage,
			GenreIDs:         r.GenreIDs,
			Adult:            r.Adult,
		}
		if r.MediaType == "tv" {
			item.Name = r.Name
			item.OriginalName = r.OriginalName
			item.FirstAirDate = r.FirstAirDate
		} else if r.MediaType == "movie" {
			item.Name = r.Title
			item.OriginalName = r.OriginalTitle
			item.FirstAirDate = r.ReleaseDate
		}
		searchResult.Results = append(searchResult.Results, item)

	}
	return searchResult, nil
}

func (c *Client) GetEposideDetail(id, seasonNumber, eposideNumber int, language string) (*tmdb.TVEpisodeDetails, error) {

	d, err := c.tmdbClient.GetTVEpisodeDetails(id, seasonNumber, eposideNumber, withLangOption(language))
	if err != nil {
		return nil, err
	}
	if !episodeNameUseful(d.Name) {
		var detailEN *tmdb.TVEpisodeDetails
		if language == "zh-CN" {
			detailEN, err = c.tmdbClient.GetTVEpisodeDetails(id, seasonNumber, eposideNumber, withLangOption("en-US"))
			if err != nil {
				return d, nil
			}
		}
		if episodeNameUseful(detailEN.Name) {
			return d, err
		}
		d.Name = detailEN.Name
		d.Overview = detailEN.Overview	
	}

	return d, err
}

func (c *Client) GetSeasonDetails(id, seasonNumber int, language string) (*tmdb.TVSeasonDetails, error) {
	detailCN, err := c.tmdbClient.GetTVSeasonDetails(id, seasonNumber, withLangOption(language))
	if err != nil {
		return nil, err
	}
	var detailEN *tmdb.TVSeasonDetails
	if language == "zh-CN" {
		detailEN, err = c.tmdbClient.GetTVSeasonDetails(id, seasonNumber, withLangOption("en-US"))
		if err != nil {
			return detailCN, nil
		}
	}

	for i, ep := range detailCN.Episodes {
		if !episodeNameUseful(ep.Name) && episodeNameUseful(detailEN.Episodes[i].Name){
			detailCN.Episodes[i].Name = detailEN.Episodes[i].Name
			detailCN.Episodes[i].Overview = detailEN.Episodes[i].Overview
		}
	}
	return detailCN, nil
}

func (c *Client) GetTVAlternativeTitles(id int, language string) (*tmdb.TVAlternativeTitles, error) {
	return c.tmdbClient.GetTVAlternativeTitles(id, withLangOption(language))
}

func wrapLanguage(lang string) string {
	if lang == "" {
		lang = "zh-CN"
	}
	return lang
}

func withLangOption(language string) map[string]string {
	language = wrapLanguage(language)
	return map[string]string{
		"language": language,
	}
}


func episodeNameUseful(name string) bool {
	return !strings.HasSuffix(name, "集") && !strings.HasPrefix(strings.ToLower(name), "episode")
}