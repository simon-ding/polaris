package core

import (
	"polaris/ent"
	"polaris/log"
	"polaris/pkg/prowlarr"
	"strings"

	"github.com/pkg/errors"
)

const prowlarrPrefix = "Prowlarr_"

func (c *Client) SyncProwlarrIndexers(apiKey, url string) error {
	client := prowlarr.New(apiKey, url)
	if ins, err := client.GetIndexers(); err != nil {
		return errors.Wrap(err, "connect to prowlarr error")
	} else {
		var prowlarrNames = make(map[string]bool, len(ins))
		for _, in := range ins {
			prowlarrNames[in.Name] = true
		}
		all := c.db.GetAllIndexers()
		for _, index := range all {
			
			if index.Synced {
				if !prowlarrNames[strings.TrimPrefix(index.Name, prowlarrPrefix)] {
					c.db.DeleteIndexer(index.ID) //remove deleted indexers
				}
			}
		}

		for _, indexer := range ins {
			if err := c.db.SaveIndexer(&ent.Indexers{
				Disabled:  indexer.Disabled,
				Name:      prowlarrPrefix + indexer.Name,
				Priority:  indexer.Priority,
				SeedRatio: indexer.SeedRatio,
				//Settings:       indexer.Settings,
				Implementation: "torznab",
				APIKey:         indexer.APIKey,
				URL:            indexer.URL,
				TvSearch:       indexer.TvSearch,
				MovieSearch:    indexer.MovieSearch,
				Synced:         true,
			}); err != nil {
				return errors.Wrap(err, "save prowlarr indexers")
			}
			log.Debugf("synced prowlarr indexer to db: %v", indexer.Name)
		}
	}
	return nil
}

func (c *Client) syncProwlarr() error {
	p, err := c.db.GetProwlarrSetting()
	if err != nil {
		return errors.Wrap(err, "db")
	}
	if p.Disabled {
		return nil
	}
	if err := c.SyncProwlarrIndexers(p.ApiKey, p.URL); err != nil {
		return errors.Wrap(err, "sync prowlarr indexers")
	}

	return nil
}


func (c *Client) DeleteAllProwlarrIndexers() error {
	all := c.db.GetAllIndexers()
	for _, index := range all {
		if index.Synced {
			c.db.DeleteIndexer(index.ID)
			log.Debugf("success delete prowlarr indexer: %s", index.Name)
		}
	}
	return nil
}