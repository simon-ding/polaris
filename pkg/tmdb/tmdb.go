package tmdb

import (
	"polaris/log"

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
	log.Infof("tv id %d, language %s", id, language)
	language = wrapLanguage(language)
	d, err := c.tmdbClient.GetTVDetails(id, withLangOption(language))
	return d, err
}

func (c *Client) SearchTvShow(query string, lang string) (*tmdb.SearchTVShows, error) {
	r, err := c.tmdbClient.GetSearchTVShow(query, withLangOption(lang))
	if err != nil {
		return nil, errors.Wrap(err, "tmdb search tv")
	}
	return r, nil
}

func (c *Client) GetEposideDetail(id, seasonNumber, eposideNumber int, language string) (*tmdb.TVEpisodeDetails, error) {
	d, err := c.tmdbClient.GetTVEpisodeDetails(id, seasonNumber, eposideNumber, withLangOption(language))
	return d, err
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
