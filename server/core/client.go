package core

import (
	"polaris/db"
	"polaris/ent"
	"polaris/log"
	"polaris/pkg/tmdb"
	"polaris/pkg/transmission"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

func NewClient(db *db.Client, language string) *Client {
	return &Client{
		db:       db,
		cron:     cron.New(),
		tasks:    make(map[int]*Task, 0),
		language: language,
	}
}

type Client struct {
	db       *db.Client
	cron     *cron.Cron
	tasks    map[int]*Task
	language string
}

func (c *Client) Init() {
	c.reloadTasks()
	c.addSysCron()
}

func (c *Client) reloadTasks() {
	allTasks := c.db.GetHistories()
	for _, t := range allTasks {
		torrent, err := transmission.ReloadTorrent(t.Saved)
		if err != nil {
			log.Errorf("relaod task %s failed: %v", t.SourceTitle, err)
			continue
		}
		if !torrent.Exists() { //只要种子还存在于客户端中，就重新加载，有可能是还在做种中
			continue
		}
		log.Infof("reloading task: %d %s", t.ID, t.SourceTitle)
		c.tasks[t.ID] = &Task{Torrent: torrent}
	}
}

func (c *Client) getDownloadClient() (*transmission.Client, *ent.DownloadClients, error) {
	tr := c.db.GetTransmission()
	trc, err := transmission.NewClient(transmission.Config{
		URL:      tr.URL,
		User:     tr.User,
		Password: tr.Password,
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "connect transmission")
	}
	return trc, tr, nil
}

func (c *Client) TMDB() (*tmdb.Client, error) {
	api := c.db.GetSetting(db.SettingTmdbApiKey)
	if api == "" {
		return nil, errors.New("TMDB apiKey not set")
	}
	return tmdb.NewClient(api)
}

func (c *Client) MustTMDB() *tmdb.Client {
	t, err := c.TMDB()
	if err != nil {
		log.Panicf("get tmdb: %v", err)
	}
	return t
}


func (c *Client) RemoveTaskAndTorrent(id int)error {
	torrent := c.tasks[id]
	if torrent != nil {
		if err := torrent.Remove(); err != nil {
			return errors.Wrap(err, "remove torrent")
		}
		delete(c.tasks, id)
	}
	return nil
}

func (c *Client) GetTasks() map[int]*Task {
	return c.tasks
}