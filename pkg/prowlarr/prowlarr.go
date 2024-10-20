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

func (c *Client) GetIndexers() ([]*db.TorznabInfo, error) {
	ins, err := c.p.GetIndexers()
	if err != nil {
		return nil, err
	}
	var indexers []*db.TorznabInfo
	for _, in := range ins {
		if !in.Enable {
			continue
		}
		seedRatio := 0.0
		for _, f := range in.Fields {
			if f.Name == "torrentBaseSettings.seedRatio" && f.Value != nil {
				if r, ok := f.Value.(float64); ok {
					seedRatio = r
				}
			}
		}
		setting := db.TorznabSetting{
			URL:    fmt.Sprintf("%s/%d/api", strings.TrimSuffix(c.url, "/"), in.ID),
			ApiKey: c.apiKey,
		}
		data, _ := json.Marshal(&setting)

		entIndexer := ent.Indexers{
			Name:           in.Name,
			Implementation: "torznab",
			Priority:       128 - int(in.Priority),
			SeedRatio:      float32(seedRatio),
			Settings:       string(data),
		}

		indexers = append(indexers, &db.TorznabInfo{
			Indexers: &entIndexer,
			TorznabSetting: setting,
		})
	}
	return indexers, nil
}
