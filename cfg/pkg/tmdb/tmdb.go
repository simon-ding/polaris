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

func (c *Client) GetTvDetails(id int, language string) {
	language = wrapLanguage(language)
	d, err := c.tmdbClient.GetTVDetails(id, map[string]string{
		"language": language,
	})
	log.Infof("error %v", err)
	log.Infof("detail %+v", d)
}


func wrapLanguage(lang string) string {
	if lang == "" {
		lang = "zh-CN"
	}
	return lang
}