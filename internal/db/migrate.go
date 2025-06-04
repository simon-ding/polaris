package db

import (
	"context"
	"encoding/json"
	"polaris/log"

	"github.com/pkg/errors"
)

func (c *client) migrate() error {
	// Run the auto migration tool.
	if err := c.ent.Schema.Create(context.Background()); err != nil {
		return errors.Wrap(err, "failed creating schema resources")
	}

	if err := c.migrateIndexerSetting(); err != nil {
		return errors.Wrap(err, "migrate indexer setting")
	}
	return nil
}

func (c *client) migrateIndexerSetting() error {
	indexers := c.GetAllIndexers()
	for _, in := range indexers {

		if in.Settings == "" {
			continue
		}
		if in.APIKey != "" && in.URL != "" {
			continue
		}
		var setting TorznabSetting
		err := json.Unmarshal([]byte(in.Settings), &setting)
		if err != nil {
			return err
		}
		in.APIKey = setting.ApiKey
		in.URL = setting.URL
		if err := c.SaveIndexer(in); err != nil {
			return errors.Wrap(err, "save indexer")
		}
		log.Infof("success migrate indexer setting field: %s", in.Name)
	}
	return nil
}
