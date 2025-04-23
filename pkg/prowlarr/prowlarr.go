package prowlarr

import (
	"encoding/json"
	"fmt"
	"polaris/db"
	"polaris/ent"
	"strings"
	"time"

	"golift.io/starr"
	"golift.io/starr/prowlarr"
)

type ProwlarrSupportType string

const (
	TV    ProwlarrSupportType = "tv"
	Movie ProwlarrSupportType = "movie"
)

type Client struct {
	p      *prowlarr.Prowlarr
	apiKey string
	url    string
}

func New(apiKey, url string) *Client {
	c := starr.New(apiKey, url, 10*time.Second)
	p := prowlarr.New(c)
	return &Client{p: p, apiKey: apiKey, url: url}
}

func (c *Client) GetIndexers() ([]*ent.Indexers, error) {
	ins, err := c.p.GetIndexers()
	if err != nil {
		return nil, err
	}
	var indexers []*ent.Indexers
	for _, in := range ins {

		tvSearch := true
		movieSearch := true
		if len(in.Capabilities.TvSearchParams) == 0 { //no tv resource in this indexer
			tvSearch = false
		}
		if len(in.Capabilities.MovieSearchParams) == 0 { //no movie resource in this indexer
			movieSearch = false
		}
		seedRatio := 0.0
		for _, f := range in.Fields {
			if f.Name == "torrentBaseSettings.seedRatio" && f.Value != nil {
				if r, ok := f.Value.(float64); ok {
					seedRatio = r
					break
				}
			}
		}
		setting := db.TorznabSetting{
			URL:    fmt.Sprintf("%s/%d/api", strings.TrimSuffix(c.url, "/"), in.ID),
			ApiKey: c.apiKey,
		}
		data, _ := json.Marshal(&setting)

		entIndexer := ent.Indexers{
			Disabled:       !in.Enable,
			Name:           in.Name,
			Implementation: "torznab",
			Priority:       int(in.Priority),
			SeedRatio:      float32(seedRatio),
			Settings:       string(data),
			TvSearch:       tvSearch,
			MovieSearch:    movieSearch,
			APIKey:         c.apiKey,
			URL:            fmt.Sprintf("%s/%d/api", strings.TrimSuffix(c.url, "/"), in.ID),
		}
		indexers = append(indexers, &entIndexer)
	}
	return indexers, nil
}
